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

package joint

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/gopa/core/errors"
	"github.com/infinitbyte/gopa/core/model"
	"github.com/infinitbyte/gopa/core/util"
	"strings"
)

type ExtractJoint struct {
	model.Parameters
}

func (joint ExtractJoint) Name() string {
	return "extract"
}

var htmlBlock model.ParaKey = "html_block"

func (joint ExtractJoint) Process(context *model.Context) error {

	snapshot := context.MustGet(model.CONTEXT_SNAPSHOT).(*model.Snapshot)

	if !util.PrefixStr(snapshot.ContentType, "text/html") {
		log.Debugf("snapshot is not html, %s, %s , %s", snapshot.ID, snapshot.Url, snapshot.ContentType)
		return nil
	}

	fileContent := snapshot.Payload

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(fileContent))
	if err != nil {
		panic(err)
	}

	kv, ok := joint.GetMap(htmlBlock)
	if ok {

		if snapshot.EnrichedFeatures == nil {
			snapshot.EnrichedFeatures = &map[string]interface{}{}
		}

		for k, v := range kv {
			o := parseTag(doc, v.(string))
			if strings.TrimSpace(o) != "" {
				(*snapshot.EnrichedFeatures)[k] = o
			}
		}
	} else {
		panic(errors.New("no extract rule was defined"))
	}

	return nil
}

func parseTag(doc *goquery.Document, tag string) string {
	found := doc.Find(tag)
	if found.Size() > 1 {
		panic(errors.Errorf("tag have multi instances, %v", tag))
	}
	ret, err := found.Html()
	if err != nil {
		panic(err)
	}
	return ret
}
