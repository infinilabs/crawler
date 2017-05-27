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

package blevesearch

import (
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/cihub/seelog"
)

func main() {
	// open a new index
	mapping := bleve.NewIndexMapping()
	index, err := bleve.New("example.bleve", mapping)
	if err != nil {
		seelog.Error(err)
	}

	data := struct {
		Name string
	}{
		Name: "text",
	}

	// index some data
	err = index.Index("id", data)
	if err != nil {
		seelog.Error(err)
	}
	// search for some text
	query := bleve.NewMatchQuery("text")
	search := bleve.NewSearchRequest(query)
	searchResults, err := index.Search(search)
	if err != nil {
		seelog.Error(err)
	}

	fmt.Println(searchResults)
}
