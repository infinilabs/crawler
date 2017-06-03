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
	JointName  string                 `json:"joint_name"` //the joint name
	Parameters map[string]interface{} `json:"parameters"` //kv parameters for this joint
}

/**
config for each pipeline, a pipeline have more than one joints
*/
type PipelineConfig struct {
	Name          string         `json:"name"`
	Context       *Context       `json:"context"`
	InputJoint    *JointConfig   `json:"input_joint"`
	ProcessJoints []*JointConfig `json:"process_joints"`
	OutputJoint   *JointConfig   `json:"output_joint"`
}
