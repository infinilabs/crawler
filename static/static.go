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
    "github.com/infinitbyte/framework/core/util"
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

	if fs.CheckLocalFirst && util.FileExists(name){
		p := path.Join(fs.BaseFolder, ".", )
		f2, err := os.Open(p)
		if err == nil {
			return f2, err
		}
	}

	f, err := fs.prepare(name)
	if err != nil {
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
}
