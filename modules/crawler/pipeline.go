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
	"github.com/infinitbyte/gopa/core/pipeline"
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
	pipeline.Register(pipe.EmptyJoint{})
	pipeline.Register(pipe.FilterCheckJoint{})
	pipeline.Register(pipe.FetchJoint{})
	pipeline.Register(pipe.UrlNormalizationJoint{})
	pipeline.Register(pipe.SaveTaskJoint{})
	pipeline.Register(pipe.HtmlToTextJoint{})
	pipeline.Register(pipe.IgnoreTimeoutJoint{})
	pipeline.Register(pipe.LoadMetadataJoint{})
	pipeline.Register(pipe.ParsePageJoint{})
	pipeline.Register(pipe.SaveSnapshotToDBJoint{})
	pipeline.Register(pipe.SaveSnapshotToFileSystemJoint{})
	pipeline.Register(pipe.InitTaskJoint{})
	pipeline.Register(pipe.UrlFilterJoint{})
	pipeline.Register(pipe.HashJoint{})
	pipeline.Register(pipe.IndexJoint{})
	pipeline.Register(pipe.TaskDeduplicationJoint{})
	pipeline.Register(pipe.ContentDeduplicationJoint{})
	pipeline.Register(pipe.UpdateCheckTimeJoint{})
	log.Debug("end register joints")

}
