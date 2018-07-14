package model

import (
	"github.com/infinitbyte/framework/core/persist"
	"github.com/infinitbyte/framework/core/util"
	"time"
)

// Host is host struct
type Host struct {
	Host        string        `json:"host,omitempty" index:"id"`
	Favicon     string        `json:"favicon,omitempty"`
	Enabled     bool          `json:"enabled"`
	HostConfigs *[]HostConfig `json:"host_configs,omitempty"`
	Created     time.Time     `json:"created,omitempty"`
	Updated     time.Time     `json:"updated,omitempty"`
}

// CreateHost create a domain host
func CreateHost(host string) Host {
	h := Host{}
	h.Host = host
	time := time.Now().UTC()
	h.Created = time
	h.Updated = time
	err := persist.Save(&h)
	if err != nil {
		panic(err)
	}
	return h
}

// GetHostList return host list
func GetHostList(from, size int, host string) (int, []Host, error) {
	var hosts []Host

	query := persist.Query{From: from, Size: size}
	if len(host) > 0 {
		query.Conds = persist.And(persist.Eq("host", host))
	}

	err, r := persist.Search(Host{}, &hosts, &query)

	if hosts == nil && r.Result != nil {
		t, ok := r.Result.([]interface{})
		if ok {
			for _, i := range t {
				js := util.ToJson(i, false)
				t := Host{}
				util.FromJson(js, &t)
				hs := GetHostConfig("", t.Host)
				if len(hosts) > 0 {
					t.HostConfigs = &hs
				}
				hosts = append(hosts, t)
			}
		}
	}

	return r.Total, hosts, err
}

// GetHost return a single host
func GetHost(host string) (Host, error) {
	var d = Host{Host: host}
	err := persist.Get(&d)
	hosts := GetHostConfig("", host)
	if len(hosts) > 0 {
		d.HostConfigs = &hosts
	}
	return d, err
}
