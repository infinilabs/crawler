package util

import (
	"net"
	log "github.com/cihub/seelog"
	"strconv"
	"errors"
)

func TestPort(port int)bool  {
	host := ":" + strconv.Itoa(port)
	ln, err := net.Listen("tcp", host)

	if err != nil {
		log.Debugf("Can't listen on port %q: %s\n", port, err)
		return false
	}
	ln.Close()
	return true
}

func TestListenPort(ip string,port int)bool  {
	host := ip+":" + strconv.Itoa(port)
	ln, err := net.Listen("tcp", host)

	if err != nil {
		log.Debugf("Can't listen on port %s:%q: %s\n",ip, port, err)
		return false
	}
	ln.Close()
	return true
}

func GetAvailablePort(ip string, port int) int  {

	maxRetry:=500

	for i :=0; i<maxRetry ; i++ {
		ok:=TestListenPort(ip,port)
		if(ok){
			return port
		}
		port++
	}

	panic(errors.New("no ports available"))
}