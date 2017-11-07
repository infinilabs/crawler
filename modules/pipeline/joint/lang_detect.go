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
	"github.com/infinitbyte/gopa/core/errors"
	"github.com/infinitbyte/gopa/core/model"
	"github.com/jbowles/nlpt-detect"
)

// LanguageDetectJoint used to detect the language of the webpage
type LanguageDetectJoint struct {
}

// Name return lang_detect
func (joint LanguageDetectJoint) Name() string {
	return "lang_detect"
}

// Process language detect
func (joint LanguageDetectJoint) Process(c *model.Context) error {
	snapshot := c.MustGet(CONTEXT_CRAWLER_SNAPSHOT).(*model.Snapshot)

	if snapshot == nil {
		return errors.Errorf("snapshot is nil, %s , %s", snapshot.ID, snapshot.Url)
	}

	if snapshot.Text != "" {
		code := nlpt_detect.GetLanguageCode(snapshot.Text)
		snapshot.Language = code
	}

	return nil
}
