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
	"sync"
	"testing"
	"time"
)

func BenchmarkGetIncrementID(b *testing.B) {

	for i := 0; i < b.N; i++ {
		GetIncrementID("a")
	}
}

func TestIDGenerator(t *testing.T) {
	var set = map[string]interface{}{}
	var s = sync.RWMutex{}
	for j := 0; j < 50; j++ {
		go func() {
			for i := 0; i < 5000000; i++ {
				id := GetUUID()
				s.Lock()
				if _, ok := set[id]; ok {
					panic(id)
				} else {
					set[id] = true
				}
				s.Unlock()
			}
		}()
	}
	time.Sleep(3 * time.Second)

}

func BenchmarkGetUUID(t *testing.B) {

	for i := 0; i < t.N; i++ {
		GetUUID()
	}
}
