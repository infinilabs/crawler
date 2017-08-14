// +build !windows

package util

import (
	"syscall"
)

// CheckProcessExists check if the pid is running
func CheckProcessExists(pid int) bool {
	err := syscall.Kill(pid, syscall.Signal(0))
	return err == nil || err == syscall.EPERM
}
