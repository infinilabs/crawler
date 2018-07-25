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
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/framework/core/fs"
)

func (fs StaticFS) prepare(name string) (*fs.VFile, error) {
	name=path.Clean(name)
	f, present := data[name]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	fs.once.Do(func() {
		f.FileName = path.Base(name)

		if f.FileSize == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.Compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			log.Error(err)
			return
		}
		f.Data, err = ioutil.ReadAll(gr)

	})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return f, nil
}

func (fs StaticFS) Open(name string) (http.File, error) {

	name=path.Clean(name)

	if fs.CheckLocalFirst {
		p := path.Join(fs.BaseFolder, ".", )
		f2, err := os.Open(p)
		if err == nil {
			return f2, err
		}
	}

	f, err := fs.prepare(name)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return f.File()
}

type StaticFS struct {
	once sync.Once
	BaseFolder      string
	CheckLocalFirst bool
}

var data = map[string]*fs.VFile{

	"/assets/config/common-config.js": {
		FileName:   "assets/config/common-config.js",
		FileSize:    58,
		ModifyTime: 1513268105,
		Compressed: `
H4sIAAAAAAAA/ypLLFJIzs/Nzc9Lzs9Ly0xXsFWo5lJQUEosyCwtylGyUlDKKCkpsNLXz8lPTszJyC8u
sbIwMDBU4qoFBAAA///AVh1/OgAAAA==
`,
	},

	"/assets/css/layout.css": {
		FileName:   "assets/css/layout.css",
		FileSize:    3386,
		ModifyTime: 1511010150,
		Compressed: `
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

	"/assets": {
		IsFolder: true,
		FileName: "/assets",
	},

	"/assets/config": {
		IsFolder: true,
		FileName: "/assets/config",
	},

	"/assets/css": {
		IsFolder: true,
		FileName: "/assets/css",
	},
}
