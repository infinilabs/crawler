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
	"github.com/infinitbyte/framework/core/pipeline"
	"github.com/infinitbyte/gopa/pipeline/joint"
)

var inited bool = false

// InitJoints register crawler joint types in order to create pipeline with these joints on the fly
func InitJoints() {

	if inited {
		return
	}
	inited = true

	log.Debug("start register joints")
	pipeline.RegisterPipeJoint(joint.EmptyJoint{})
	pipeline.RegisterPipeJoint(joint.FilterCheckJoint{})
	pipeline.RegisterPipeJoint(joint.FetchJoint{})
	pipeline.RegisterPipeJoint(joint.UrlNormalizationJoint{})
	pipeline.RegisterPipeJoint(joint.SaveTaskJoint{})
	pipeline.RegisterPipeJoint(joint.HtmlToTextJoint{})
	pipeline.RegisterPipeJoint(joint.IgnoreTimeoutJoint{})
	pipeline.RegisterPipeJoint(joint.ParsePageJoint{})
	pipeline.RegisterPipeJoint(joint.SaveSnapshotToDBJoint{})
	pipeline.RegisterPipeJoint(joint.SaveSnapshotToFileSystemJoint{})
	pipeline.RegisterPipeJoint(joint.InitTaskJoint{})
	pipeline.RegisterPipeJoint(joint.UrlFilterJoint{})
	pipeline.RegisterPipeJoint(joint.HashJoint{})
	pipeline.RegisterPipeJoint(joint.IndexJoint{})
	pipeline.RegisterPipeJoint(joint.TaskDeduplicationJoint{})
	pipeline.RegisterPipeJoint(joint.ContentDeduplicationJoint{})
	pipeline.RegisterPipeJoint(joint.UpdateCheckTimeJoint{})
	pipeline.RegisterPipeJoint(joint.LanguageDetectJoint{})
	pipeline.RegisterPipeJoint(joint.ExtractJoint{})
	log.Debug("end register joints")

}
