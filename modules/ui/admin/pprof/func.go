package pprof

import (
	"fmt"
	"github.com/infinitbyte/gopa/core/util"
	"regexp"
	"sort"
)

var reg = regexp.MustCompile("\\s+?/[\\w|\\W]*?\\.go:\\d+")
var split = regexp.MustCompile(`\d+: \d+ \[\d+: \d+\]`)
var snapshotValue = map[string]int{}

func getSnapshot() []string {
	url := "http://localhost:6060/debug/pprof/heap?debug=1"
	r, err := util.HttpGet(url)
	if err != nil {
		panic(err)
	}

	str := string(r.Body)
	strs := split.Split(str, -1)

	result := []string{}

	for _, s := range strs {
		all := reg.FindAllString(s, -1)
		if len(all) > 1 {
			s1 := util.MergeSpace(all[0]) + " - " + util.MergeSpace(all[len(all)-1])
			result = append(result, s1)
		}
	}

	return result
}

func takeSnapshot() {
	all := getSnapshot()
	snapshotValue = map[string]int{}

	for _, k := range all {
		k = util.MergeSpace(k)
		v, ok := snapshotValue[k]

		if !ok {
			snapshotValue[k] = 1
		} else {
			snapshotValue[k] = v + 1
		}
	}
}

func compareNow() map[string]int {
	nowValue := map[string]int{}

	all := getSnapshot()
	for _, k := range all {
		k = util.MergeSpace(k)
		v, ok := nowValue[k]

		if !ok {
			nowValue[k] = 1
		} else {
			nowValue[k] = v + 1
		}
	}

	newValue := map[string]int{}
	for k := range nowValue {
		newValue[k] = nowValue[k]
	}

	for k := range snapshotValue {
		snapV := snapshotValue[k]

		v, ok := newValue[k]

		if !ok {
			//not exist in new snapshot, must disappeared
			newValue[k] = 0 - snapV
		} else {
			//get diff value
			newValue[k] = v - snapV
		}
	}

	finalValue := map[string]int{}
	for k, v := range newValue {
		if v != 0 {
			finalValue[k] = v
		}
	}

	//sort by value
	type kv struct {
		Key   string
		Value int
	}

	var ss []kv
	for k, v := range finalValue {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	fmt.Println("- - ")
	len := len(ss)
	if len > 10 {
		len = 10
	}
	for i := 0; i < len; i++ {
		kv := ss[i]
		fmt.Printf("%d: %s, %d\n", i, kv.Key, kv.Value)
	}

	//fmt.Println(util.ToJson(finalValue, true))

	return newValue

}
