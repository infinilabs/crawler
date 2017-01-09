package util

import (
	log "github.com/cihub/seelog"
	"os"
	"path"
	"time"
)

var locked bool
var file string

func CheckInstanceLock(p string) {
	file = path.Join(p, ".lock")
	if FileExists(file) {
		log.Errorf("lock file:%s exists, if you only have one instance, please remove it", file)
		log.Flush()
		os.Exit(1)
	}
	FilePutContent(file, time.Now().String())
	log.Trace("lock placed,", file)
	locked = true
}

func ClearInstanceLock() {
	if locked {
		os.Remove(path.Join(file))
	}
}
