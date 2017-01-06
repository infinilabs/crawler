package util

import (
	"errors"
	log "github.com/cihub/seelog"
	"net"
	"strconv"
	"strings"
)

func TestPort(port int) bool {
	host := ":" + strconv.Itoa(port)
	ln, err := net.Listen("tcp", host)

	if err != nil {
		log.Debugf("can't listen on port %s, %s", host, err)
		return false
	}
	ln.Close()
	return true
}

func TestListenPort(ip string, port int) bool {

	log.Tracef("testing port %s:%d", ip, port)
	host := ip + ":" + strconv.Itoa(port)
	ln, err := net.Listen("tcp", host)

	if err != nil {
		log.Debugf("can't listen on port %s, %s", host, err)
		return false
	}
	ln.Close()
	return true
}

/**
get valid port to listen, if the specify port is not available, auto choose the next one
*/
func GetAvailablePort(ip string, port int) int {

	maxRetry := 500

	for i := 0; i < maxRetry; i++ {
		ok := TestListenPort(ip, port)
		if ok {
			return port
		}
		port++
	}

	panic(errors.New("no ports available"))
}

/**
get valid address to listen, if the specify port is not available, auto choose the next one
*/
func AutoGetAddress(addr string) string {
	if strings.Index(addr, ":") < 0 {
		panic(errors.New("invalid address, eg ip:port, " + addr))
	}

	array := strings.Split(addr, ":")
	p, _ := strconv.Atoi(array[1])
	port := GetAvailablePort(array[0], p)
	array[1] = strconv.Itoa(port)
	return strings.Join(array, ":")
}

/**
get valid address, input: :8001 -> output: 127.0.0.1:8001
*/
func GetValidAddress(addr string) string {
	if strings.Index(addr, ":") >= 0 {
		array := strings.Split(addr, ":")
		if len(array[0]) == 0 {
			array[0] = "127.0.0.1"
			addr = strings.Join(array, ":")
		}
	}
	return addr
}
