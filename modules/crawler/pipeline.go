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

package crawler

import (
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/gopa/core/model"
	"github.com/infinitbyte/gopa/modules/crawler/pipe"
)

var inited bool = false

// InitJoints register crawler joint types in order to create pipeline with these joints on the fly
func InitJoints() {

	if inited {
		return
	}
	inited = true

	log.Debug("start register joints")
	model.RegisterPipeJoint(pipe.EmptyJoint{})
	model.RegisterPipeJoint(pipe.FilterCheckJoint{})
	model.RegisterPipeJoint(pipe.FetchJoint{})
	model.RegisterPipeJoint(pipe.UrlNormalizationJoint{})
	model.RegisterPipeJoint(pipe.SaveTaskJoint{})
	model.RegisterPipeJoint(pipe.HtmlToTextJoint{})
	model.RegisterPipeJoint(pipe.IgnoreTimeoutJoint{})
	model.RegisterPipeJoint(pipe.LoadMetadataJoint{})
	model.RegisterPipeJoint(pipe.ParsePageJoint{})
	model.RegisterPipeJoint(pipe.SaveSnapshotToDBJoint{})
	model.RegisterPipeJoint(pipe.SaveSnapshotToFileSystemJoint{})
	model.RegisterPipeJoint(pipe.InitTaskJoint{})
	model.RegisterPipeJoint(pipe.UrlFilterJoint{})
	model.RegisterPipeJoint(pipe.HashJoint{})
	model.RegisterPipeJoint(pipe.IndexJoint{})
	model.RegisterPipeJoint(pipe.TaskDeduplicationJoint{})
	model.RegisterPipeJoint(pipe.ContentDeduplicationJoint{})
	model.RegisterPipeJoint(pipe.UpdateCheckTimeJoint{})
	model.RegisterPipeJoint(pipe.LanguageDetectJoint{})
	log.Debug("end register joints")

}
