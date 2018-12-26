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
	"github.com/infinitbyte/gopa/pipeline/joints/filter"
	"github.com/infinitbyte/gopa/pipeline/joints/input"
	"github.com/infinitbyte/gopa/pipeline/joints/output"
)

var inited bool = false

// InitJoints register crawler joint types in order to create pipeline with these joints on the fly
func InitJoints() {

	if inited {
		return
	}
	inited = true

	log.Debug("start register joints")
	pipeline.RegisterPipeJoint(filter.EmptyJoint{})
	pipeline.RegisterPipeJoint(filter.FilterCheckJoint{})
	pipeline.RegisterPipeJoint(filter.FetchJoint{})
	pipeline.RegisterPipeJoint(filter.UrlNormalizationJoint{})
	pipeline.RegisterPipeJoint(output.SaveTaskJoint{})
	pipeline.RegisterPipeJoint(filter.HtmlToTextJoint{})
	pipeline.RegisterPipeJoint(filter.IgnoreTimeoutJoint{})
	pipeline.RegisterPipeJoint(filter.ParsePageJoint{})
	pipeline.RegisterPipeJoint(filter.SaveSnapshotToDBJoint{})
	pipeline.RegisterPipeJoint(filter.SaveSnapshotToFileSystemJoint{})
	pipeline.RegisterPipeJoint(input.InitTaskJoint{})
	pipeline.RegisterPipeJoint(filter.UrlFilterJoint{})
	pipeline.RegisterPipeJoint(filter.HashJoint{})
	pipeline.RegisterPipeJoint(filter.IndexJoint{})
	pipeline.RegisterPipeJoint(filter.TaskDeduplicationJoint{})
	pipeline.RegisterPipeJoint(filter.ContentDeduplicationJoint{})
	pipeline.RegisterPipeJoint(filter.UpdateCheckTimeJoint{})
	pipeline.RegisterPipeJoint(filter.ExtractJoint{})
	log.Debug("end register joints")

}
