package util

import (
	"os"
	"net"
	log "github.com/cihub/seelog"
)

func TestPort(port string)bool  {
	ln, err := net.Listen("tcp", ":" + port)

	if err != nil {
		log.Debug(os.Stderr, "Can't listen on port %q: %s\n", port, err)
		return false
	}
	ln.Close()
	return true
}