//Package config , actually copied from github.com/elastic/beats
package config

import (
	"errors"
	"flag"
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/elastic/go-ucfg"
	"github.com/elastic/go-ucfg/cfgutil"
	cfgflag "github.com/elastic/go-ucfg/flag"
	"github.com/elastic/go-ucfg/yaml"
	"github.com/infinitbyte/gopa/core/util/file"
	"os"
	"path/filepath"
	"runtime"
)

// IsStrictPerms returns true if strict permission checking on config files is
// enabled.
func IsStrictPerms() bool {
	return false
}

// Config object to store hierarchical configurations into.
// See https://godoc.org/github.com/elastic/go-ucfg#Config
type Config ucfg.Config

// Namespace storing at most one configuration section by name and sub-section.
type Namespace struct {
	C map[string]*Config `config:",inline"`
}

type flagOverwrite struct {
	config *ucfg.Config
	path   string
	value  string
}

var configOpts = []ucfg.Option{
	ucfg.PathSep("."),
	ucfg.ResolveEnv,
	ucfg.VarExp,
}

// NewConfig create a pretty new config
func NewConfig() *Config {
	return fromConfig(ucfg.New())
}

// NewConfigFrom get config instance
func NewConfigFrom(from interface{}) (*Config, error) {
	c, err := ucfg.NewFrom(from, configOpts...)
	return fromConfig(c), err
}

// MergeConfigs just merge configs together
func MergeConfigs(cfgs ...*Config) (*Config, error) {
	config := NewConfig()
	for _, c := range cfgs {
		if err := config.Merge(c); err != nil {
			return nil, err
		}
	}
	return config, nil
}

// NewConfigWithYAML load config from yaml
func NewConfigWithYAML(in []byte, source string) (*Config, error) {
	opts := append(
		[]ucfg.Option{
			ucfg.MetaData(ucfg.Meta{Source: source}),
		},
		configOpts...,
	)
	c, err := yaml.NewConfig(in, opts...)
	return fromConfig(c), err
}

// NewFlagConfig will use flags
func NewFlagConfig(
	set *flag.FlagSet,
	def *Config,
	name string,
	usage string,
) *Config {
	opts := append(
		[]ucfg.Option{
			ucfg.MetaData(ucfg.Meta{Source: "command line flag"}),
		},
		configOpts...,
	)

	var to *ucfg.Config
	if def != nil {
		to = def.access()
	}

	config := cfgflag.ConfigVar(set, to, name, usage, opts...)
	return fromConfig(config)
}

// NewFlagOverwrite will use flags which specified
func NewFlagOverwrite(
	set *flag.FlagSet,
	config *Config,
	name, path, def, usage string,
) *string {
	if config == nil {
		panic("Missing configuration")
	}
	if path == "" {
		panic("empty path")
	}

	if def != "" {
		err := config.SetString(path, -1, def)
		if err != nil {
			panic(err)
		}
	}

	f := &flagOverwrite{
		config: config.access(),
		path:   path,
		value:  def,
	}

	if set == nil {
		flag.Var(f, name, usage)
	} else {
		set.Var(f, name, usage)
	}

	return &f.value
}

// LoadFile will load config from specify file
func LoadFile(path string) (*Config, error) {
	if IsStrictPerms() {
		if err := ownerHasExclusiveWritePerms(path); err != nil {
			return nil, err
		}
	}

	c, err := yaml.NewConfigWithFile(path, configOpts...)
	if err != nil {
		return nil, err
	}

	cfg := fromConfig(c)

	log.Debugf("load config file '%v'", path)
	return cfg, err
}

// LoadFiles will load configs from specify files
func LoadFiles(paths ...string) (*Config, error) {
	merger := cfgutil.NewCollector(nil, configOpts...)
	for _, path := range paths {
		cfg, err := LoadFile(path)
		if err := merger.Add(cfg.access(), err); err != nil {
			return nil, err
		}
	}
	return fromConfig(merger.Config()), nil
}

// Merge a map, a slice, a struct or another Config object into c.
func (c *Config) Merge(from interface{}) error {
	return c.access().Merge(from, configOpts...)
}

// Unpack unpacks c into a struct, a map, or a slice allocating maps, slices,
// and pointers as necessary.
func (c *Config) Unpack(to interface{}) error {
	return c.access().Unpack(to, configOpts...)
}

// Path gets the absolute path of c separated by sep. If c is a root-Config an
// empty string will be returned.
func (c *Config) Path() string {
	return c.access().Path(".")
}

// PathOf gets the absolute path of a potential setting field in c with name
// separated by sep.
func (c *Config) PathOf(field string) string {
	return c.access().PathOf(field, ".")
}

// HasField checks if c has a top-level named key name.
func (c *Config) HasField(name string) bool {
	return c.access().HasField(name)
}

// CountField returns number of entries in a table or 1 if entry is a primitive value.
// Primitives settings can be handled like a list with 1 entry.
func (c *Config) CountField(name string) (int, error) {
	return c.access().CountField(name)
}

// Bool reads a boolean setting returning an error if the setting has no
// boolean value.
func (c *Config) Bool(name string, idx int) (bool, error) {
	return c.access().Bool(name, idx, configOpts...)
}

// Strings reads a string setting returning an error if the setting has
// no string or primitive value convertible to string.
func (c *Config) String(name string, idx int) (string, error) {
	return c.access().String(name, idx, configOpts...)
}

// Int reads an int64 value returning an error if the setting is
// not integer value, the primitive value is not convertible to int or a conversion
// would create an integer overflow.
func (c *Config) Int(name string, idx int) (int64, error) {
	return c.access().Int(name, idx, configOpts...)
}

// Float reads a float64 value returning an error if the setting is
// not a float value or the primitive value is not convertible to float.
func (c *Config) Float(name string, idx int) (float64, error) {
	return c.access().Float(name, idx, configOpts...)
}

// Child returns a child configuration or an error if the setting requested is a
// primitive value only.
func (c *Config) Child(name string, idx int) (*Config, error) {
	sub, err := c.access().Child(name, idx, configOpts...)
	return fromConfig(sub), err
}

// SetBool sets a boolean primitive value. An error is returned if the new name
// is invalid.
func (c *Config) SetBool(name string, idx int, value bool) error {
	return c.access().SetBool(name, idx, value, configOpts...)
}

// SetInt sets an integer primitive value. An error is returned if the new
// name is invalid.
func (c *Config) SetInt(name string, idx int, value int64) error {
	return c.access().SetInt(name, idx, value, configOpts...)
}

// SetFloat sets an floating point primitive value. An error is returned if
// the name is invalid.
func (c *Config) SetFloat(name string, idx int, value float64) error {
	return c.access().SetFloat(name, idx, value, configOpts...)
}

// SetString sets string value. An error is returned if the name is invalid.
func (c *Config) SetString(name string, idx int, value string) error {
	return c.access().SetString(name, idx, value, configOpts...)
}

// SetChild adds a sub-configuration. An error is returned if the name is invalid.
func (c *Config) SetChild(name string, idx int, value *Config) error {
	return c.access().SetChild(name, idx, value.access(), configOpts...)
}

// IsDict checks if c has named keys.
func (c *Config) IsDict() bool {
	return c.access().IsDict()
}

// IsArray checks if c has index only accessible settings.
func (c *Config) IsArray() bool {
	return c.access().IsArray()
}

// Enabled was a predefined config, enabled by default if no config was found
func (c *Config) Enabled() bool {
	testEnabled := struct {
		Enabled bool `config:"enabled"`
	}{true}

	if c == nil {
		return true
	}
	if err := c.Unpack(&testEnabled); err != nil {
		// if unpacking fails, expect 'enabled' being set to default value
		return true
	}
	return testEnabled.Enabled
}

func fromConfig(in *ucfg.Config) *Config {
	return (*Config)(in)
}

func (c *Config) access() *ucfg.Config {
	return (*ucfg.Config)(c)
}

// GetFields returns a list of all top-level named keys in c.
func (c *Config) GetFields() []string {
	return c.access().GetFields()
}

func (f *flagOverwrite) String() string {
	return f.value
}

func (f *flagOverwrite) Set(v string) error {
	opts := append(
		[]ucfg.Option{
			ucfg.MetaData(ucfg.Meta{Source: "command line flag"}),
		},
		configOpts...,
	)

	err := f.config.SetString(f.path, -1, v, opts...)
	if err != nil {
		return err
	}
	f.value = v
	return nil
}

func (f *flagOverwrite) Get() interface{} {
	return f.value
}

// Validate checks at most one sub-namespace being set.
func (ns *Namespace) Validate() error {
	if len(ns.C) > 1 {
		return errors.New("more then one namespace configured")
	}
	return nil
}

// Name returns the configuration sections it's name if a section has been set.
func (ns *Namespace) Name() string {
	for name := range ns.C {
		return name
	}
	return ""
}

// Config return the sub-configuration section if a section has been set.
func (ns *Namespace) Config() *Config {
	for _, cfg := range ns.C {
		return cfg
	}
	return nil
}

// IsSet returns true if a sub-configuration section has been set.
func (ns *Namespace) IsSet() bool {
	return len(ns.C) != 0
}

// ownerHasExclusiveWritePerms asserts that the current user or root is the
// owner of the config file and that the config file is (at most) writable by
// the owner or root (e.g. group and other cannot have write access).
func ownerHasExclusiveWritePerms(name string) error {
	if runtime.GOOS == "windows" {
		return nil
	}

	info, err := file.Stat(name)
	if err != nil {
		return err
	}

	euid := os.Geteuid()
	fileUID, _ := info.UID()
	perm := info.Mode().Perm()

	if fileUID != 0 && euid != fileUID {
		return fmt.Errorf(`config file ("%v") must be owned by the beat user `+
			`(uid=%v) or root`, name, euid)
	}

	// Test if group or other have write permissions.
	if perm&0022 > 0 {
		nameAbs, err := filepath.Abs(name)
		if err != nil {
			nameAbs = name
		}
		return fmt.Errorf(`config file ("%v") can only be writable by the `+
			`owner but the permissions are "%v" (to fix the permissions use: `+
			`'chmod go-w %v')`,
			name, perm, nameAbs)
	}

	return nil
}
