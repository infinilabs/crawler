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

/**
config for each joint
*/
type JointConfig struct {
	JointName  string                 `json:"joint" config:"joint"`           //the joint name
	Parameters map[string]interface{} `json:"parameters" config:"parameters"` //kv parameters for this joint
	Enabled    bool                   `json:"enabled" config:"enabled"`
}

/**
config for each pipeline, a pipeline have more than one joints
*/
type PipelineConfig struct {
	Name          string         `json:"name" config:"name"`
	Context       *Context       `json:"context" config:"context"`
	StartJoint    *JointConfig   `json:"start" config:"start"`
	ProcessJoints []*JointConfig `json:"process" config:"process"`
	EndJoint      *JointConfig   `json:"end" config:"end"`
}
