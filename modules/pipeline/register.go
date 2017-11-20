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

package pipeline

import (
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/gopa/core/model"
	"github.com/infinitbyte/gopa/modules/pipeline/joint"
)

var inited bool = false

// InitJoints register crawler joint types in order to create pipeline with these joints on the fly
func InitJoints() {

	if inited {
		return
	}
	inited = true

	log.Debug("start register joints")
	model.RegisterPipeJoint(joint.EmptyJoint{})
	model.RegisterPipeJoint(joint.FilterCheckJoint{})
	model.RegisterPipeJoint(joint.FetchJoint{})
	model.RegisterPipeJoint(joint.UrlNormalizationJoint{})
	model.RegisterPipeJoint(joint.SaveTaskJoint{})
	model.RegisterPipeJoint(joint.HtmlToTextJoint{})
	model.RegisterPipeJoint(joint.IgnoreTimeoutJoint{})
	model.RegisterPipeJoint(joint.LoadMetadataJoint{})
	model.RegisterPipeJoint(joint.ParsePageJoint{})
	model.RegisterPipeJoint(joint.SaveSnapshotToDBJoint{})
	model.RegisterPipeJoint(joint.SaveSnapshotToFileSystemJoint{})
	model.RegisterPipeJoint(joint.InitTaskJoint{})
	model.RegisterPipeJoint(joint.UrlFilterJoint{})
	model.RegisterPipeJoint(joint.HashJoint{})
	model.RegisterPipeJoint(joint.IndexJoint{})
	model.RegisterPipeJoint(joint.TaskDeduplicationJoint{})
	model.RegisterPipeJoint(joint.ContentDeduplicationJoint{})
	model.RegisterPipeJoint(joint.UpdateCheckTimeJoint{})
	model.RegisterPipeJoint(joint.LanguageDetectJoint{})
	model.RegisterPipeJoint(joint.ChromeFetchJoint{})
	model.RegisterPipeJoint(joint.ExtractJoint{})
	log.Debug("end register joints")

}
