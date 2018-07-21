package static

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// FS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func FS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// Dir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func Dir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// FSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func FSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		f.Close()
		return b, err
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// FSMustByte is the same as FSByte, but panics if name is not present.
func FSMustByte(useLocal bool, name string) []byte {
	b, err := FSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// FSString is the string version of FSByte.
func FSString(useLocal bool, name string) (string, error) {
	b, err := FSByte(useLocal, name)
	return string(b), err
}

// FSMustString is the string version of FSMustByte.
func FSMustString(useLocal bool, name string) string {
	return string(FSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/static/assets/config/common-config.js": {
		local:   "../static/assets/config/common-config.js",
		size:    58,
		modtime: 1513268105,
		compressed: `
H4sIAAAAAAAA/ypLLFJIzs/Nzc9Lzs9Ly0xXsFWo5lJQUEosyCwtylGyUlDKKCkpsNLXz8lPTszJyC8u
sbIwMDBU4qoFBAAA///AVh1/OgAAAA==
`,
	},

	"/static/assets/css/layout.css": {
		local:   "../static/assets/css/layout.css",
		size:    3386,
		modtime: 1511010150,
		compressed: `
H4sIAAAAAAAA/6xWwY7bNhC96ysGzmXXlWzZsXc3Mhq0TRqkQBMESHvoLZRISawpUiAprx3BQD6jP9Ci
ubT3Fu1h90/yJQVJSSvZ3uayq4NX5Axn+ObNG03HHozhmShKwQnXEbi/5yJRZuM5UYmkpaaCu63nJKWc
KFB6x4iCVEjQOYEfv1tTDVgkVUG4Rsbegy8f7A/GU8/zpmN4IbhWD35yLPDOh3zmQz73IX/sQ77wIV/6
kF/43qRaBxxtYiTNDzwFRuEpIKghFVwHKSoo20UweiPFlhYIXosNGvkweknYhmiaIHhNKjLyoVvwR28o
z14gnsHbZyN/9IomUiiRavgJvSR05I9uf7359+aP299v/rn58/bjpw9/ffrw2+3Hm78/ffhl5MPXkiLm
g0JcBYpImq5g7/D5Biny4PDkumBQe6b8YkNkysR1sItAJVIwtrLrwTWJ11QHFhFVCKFzyrMIENcUMYoU
wY1hId4HQm2PLDOJdipBjKy8vWcxjyutBYc6FhITGUiEaaUiWJTbVWOhUcwIaOxD7y1vUm3cYqG1KCKY
lVtQglEMj8jSPDbOdAzfo52o9MOz1XOt9YpizIhn1ya6CAr7DjXEKFlnUlQcR/AoTW0NewbmRongGlFO
ZA0F5UFOaJbrCC7C0GJwnzWYDUUxiZF0WJQIY8qzgJFURxC6UhwfaZEfHmpBnTgs+7gOAP12aZ4D9yAn
CBMJdT/O48PMW6tJJikOCiQzymv3E2hRRrPloYNp1qbE/ZMXS3sDu2zdI5iH5dZcduCeCqHvycqw4Qdz
Xw2IY4hplhFpCvdVQTBFcGZ8rinWeQSXF1fl9hxqz0bsAd6k1mOgdEFO4HVg6Kpzj93+OJKpT4k4YVB3
FW6CmatbUXC3cvKctCoPhcA0pUQ+vJQ60r9GG8t4k22rmYOuNKW976alUNQNHEkY0nRD2hK2R0UMKR0k
OWX4c71+QM1OwCeMBomouK4PYqJYCVZp4nJpa9dxyyXevSaCCen617SwS6Zr7KDdDkN0RS5W/V6MQlj2
Wm44WVq/y8vL1f7EvkuFbHWAGM14BAnhmkgXYMBqu8IoJ8HxspNg+p5EMLtocpmO4aUR+cPqmbhRbnb8
4/VUJJUarhtuosQUz2bcanJf814szTPEcbFYtGm8reJAxD+TREfwjqNNoxTv+pk1o7mTmsFZ8/m8d8/r
5vpcyAKxPu62lZoBb+vtjkmZQDpyDGhy6lPbeqsq7jjJde4oeca/mJ9D3chQ09VdQ07HY898x8iiOaZa
B6mQBbT/WB2XgimYMBQTZgefgPqgw2eNkvbdVRUXVEPdaNSVMWnF0Ni7LA6cEsQTIyF9p3YCI5yRBlVM
VcnQLgLKLZ1iJpL1kNHQUNqsTceDWofhE4yvVuPpEfGGbGyrdBmGJ6g7W5zou9X/NsOGSPPFxdo9NwV6
PloirgwYhhq82TmecTKL0Vnom2cyP18dq07b633Dx+dHiiAkteXofE/gFdACZSTqPqoMBkgGmSEC4fpM
i9J/FIbxIl2a3yucLM9baE+dc+wPLud7Tzn45Jp36mcQUznC4tpUOzD4hEfY7L0TE/PJRVhuzx2XOqrn
QtL35rOF3a1Z0t8N0YaWs7AjykG5mwZtdwY0igXD7fj8bOCu7+6CD1p4NuvYvff23n8BAAD//x/s7pY6
DQAA
`,
	},

	"/": {
		isDir: true,
		local: "/",
	},

	"/static": {
		isDir: true,
		local: "/static",
	},

	"/static/assets": {
		isDir: true,
		local: "/static/assets",
	},

	"/static/assets/config": {
		isDir: true,
		local: "/static/assets/config",
	},

	"/static/assets/css": {
		isDir: true,
		local: "/static/assets/css",
	},
}
