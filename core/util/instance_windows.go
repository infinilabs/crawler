package util

import (
	log "github.com/cihub/seelog"
	"os"
	"path"
	"syscall"
)

// CheckProcessExists check if the pid is running
func CheckProcessExists(pid int) bool {
	log.Warn("process running check is not supported on Windows, please manually check with your working dir")
	return true
}
