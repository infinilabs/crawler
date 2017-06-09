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
	"github.com/medcl/gopa/core/pipeline"
	"github.com/medcl/gopa/modules/crawler/pipe"
)

func InitJoints() {
	log.Debug("start register joints")
	pipeline.Register(pipe.Empty, pipe.EmptyJoint{})
	pipeline.Register(pipe.UrlCheckedFilter, pipe.UrlCheckedFilterJoint{})
	pipeline.Register(pipe.Fetch, pipe.FetchJoint{})
	pipeline.Register(pipe.UrlNormalization, pipe.UrlNormalizationJoint{})
	pipeline.Register(pipe.SaveTask, pipe.SaveTaskJoint{})
	pipeline.Register(pipe.HtmlToText, pipe.HtmlToTextJoint{})
	pipeline.Register(pipe.IgnoreTimeout, pipe.IgnoreTimeoutJoint{})
	pipeline.Register(pipe.LoadMetadata, pipe.LoadMetadataJoint{})
	pipeline.Register(pipe.ParsePage, pipe.ParsePageJoint{})
	pipeline.Register(pipe.SaveSnapshotToDB, pipe.SaveSnapshotToDBJoint{})
	pipeline.Register(pipe.SaveSnapshotToFileSystem, pipe.SaveSnapshotToFileSystemJoint{})
	pipeline.Register(pipe.InitTask, pipe.InitTaskJoint{})
	pipeline.Register(pipe.UrlExtFilter, pipe.UrlExtFilterJoint{})
	pipeline.Register(pipe.Hash, pipe.HashJoint{})
	log.Debug("end register joints")

}
