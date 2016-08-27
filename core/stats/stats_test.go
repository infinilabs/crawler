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

package stats

import "testing"

func TestStats(t *testing.T) {

	domain:="www.google.com"
	key := "key1"
	key2 := "key2"
	Increment(domain,key)

	result := Stat(domain,key)
	expected := 1
	if result != expected {
		t.Errorf("Value '%d' doesn't match expected '%d'!", result, expected)
	}

	Decrement(domain,key)
	result = Stat(domain,key)
	expected = 0
	if result != expected {
		t.Errorf("Value '%d' doesn't match expected '%d'!", result, expected)
	}

	Increment(domain,key2)
	data := StatsAll()

	expected = 0
	if data[domain][key] != expected {
		t.Errorf("Value '%d' doesn't match expected '%d'!", result, expected)
	}

	expected = 1
	if data[domain][key2] != expected {
		t.Errorf("Value '%d' doesn't match expected '%d'!", data[key2], expected)
	}

	IncrementBy(domain,key, 1)
	DecrementBy(domain,key2, 3)

	data = StatsAll()

	expected = 1
	if data[domain][key] != expected {
		t.Errorf("Value '%d' doesn't match expected '%d'!", data[key], expected)
	}

	expected = -2
	if data[domain][key2] != expected {
		t.Errorf("Value '%d' doesn't match expected '%d'!", data[key2], expected)
	}

}
