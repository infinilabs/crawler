/*
Copyright 2016 Medcl (m AT medcl.net)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestIDGenerator(t *testing.T) {
	SetIDPersistencePath("/tmp/")
	for j := 0; j < 5; j++ {
		go func() {
			for i := 0; i < 5; i++ {
				id := GetIncrementID("a")
				fmt.Println(id)
			}
		}()
		go func() {
			for i := 0; i < 5; i++ {
				id := GetIncrementID("b")
				fmt.Println(id)
			}
		}()
	}

	time.Sleep(1 * time.Second)

}
func TestIDGenerator1(t *testing.T) {
	var set = map[string]interface{}{}
	var s = sync.RWMutex{}
	for j := 0; j < 5; j++ {
		go func() {
			for i := 0; i < 5; i++ {
				id := GetUUID()
				fmt.Println(id)
				s.Lock()
				if _, ok := idseed[id]; ok {
					panic(id)
				} else {
					set[id] = true
				}
				s.Unlock()
			}
		}()
	}

	time.Sleep(1 * time.Second)

}
