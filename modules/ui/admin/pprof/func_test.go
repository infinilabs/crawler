package pprof

import (
	"testing"
	"time"
)

func TestSnapshot(t *testing.T) {

	takeSnapshot()
	for {
		time.Sleep(20 * time.Second)
		compareNow()
	}

}
