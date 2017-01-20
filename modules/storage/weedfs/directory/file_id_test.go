package directory

import (
	"fmt"
	"testing"
)

func TestSerialDeserialization(t *testing.T) {
	f1 := &FileId{VolumeId: 345, Key: 8698, Hashcode: 23849095}
	fmt.Println("vid", f1.VolumeId, "key", f1.Key, "hash", f1.Hashcode)

	//f2 := ParseFileId(t.String())
	//
	//fmt.Println("vvid", f2.VolumeId, "vkey", f2.Key, "vhash", f2.Hashcode)
}
