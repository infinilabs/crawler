/*
Copyright Medcl (m AT medcl.net)

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

package filter

import (
	"bytes"
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/framework/core/pipeline"
	"github.com/infinitbyte/framework/core/util"
	"github.com/infinitbyte/gopa/model"
	"github.com/ledongthuc/pdf"
)

type ParsePDFJoint struct {
	pipeline.Parameters
}

func (joint ParsePDFJoint) Name() string {
	return "parse_pdf"
}

const contentTypeKey pipeline.ParaKey = "accepted_content_type"

func (joint ParsePDFJoint) Process(context *pipeline.Context) error {

	snapshot := context.MustGet(model.CONTEXT_SNAPSHOT).(*model.Snapshot)

	acceptedContentType, ok := joint.GetStringArray(contentTypeKey)

	valid := false
	if ok {
		valid = false
		for _, v := range acceptedContentType {
			if snapshot.ContentType == v {
				valid = true
				break
			}
		}
	} else {
		if snapshot.ContentType == "application/pdf" {
			valid = true
		}
	}

	if !valid {
		log.Debugf("snapshot is supported or not pdf, %s, %s , %s", snapshot.ID, snapshot.Url, snapshot.ContentType)
		return nil
	}

	r := bytes.NewReader(snapshot.Payload)

	size := len(snapshot.Payload)
	f, err := pdf.NewReader(r, int64(size))
	if err != nil {
		log.Error(err)
		return err
	}

	var buf bytes.Buffer
	b, err := f.GetPlainText()
	if err != nil {
		log.Error(err)
		return err
	}

	_, err = buf.ReadFrom(b)
	if err != nil {
		log.Error(err)
		return err
	}

	snapshot.Text = util.XSSHandle(buf.String())
	if snapshot.Title == "" && len(snapshot.RawFileName) > 0 {
		snapshot.Title = util.NoWordBreak(util.XSSHandle(snapshot.RawFileName))
	}

	return nil
}

func ParsePDF2Text(file string) string {
	content, err := readPdf(file)
	if err != nil {
		panic(err)
	}
	return content
}

func readPdf(path string) (string, error) {
	f, r, err := pdf.Open(path)
	defer f.Close()
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	buf.ReadFrom(b)
	return buf.String(), nil
}
