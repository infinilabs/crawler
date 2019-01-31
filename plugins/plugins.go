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

package plugins

import (
	"github.com/infinitbyte/framework/core/module"
	"github.com/infinitbyte/framework/core/pipeline"
	"github.com/infinitbyte/gopa/plugins/chrome"
	"github.com/infinitbyte/gopa/plugins/dispatch"
	"github.com/infinitbyte/gopa/plugins/lang_detect"
)

func Register() {
	module.RegisterUserPlugin(dispatch.DispatchModule{})
	module.RegisterUserPlugin(chrome.ChromePlugin{})
	//module.RegisterUserPlugin(tools_generator.GeneratorPlugin{})
	pipeline.RegisterPipeJoint(lang_detect.LanguageDetectJoint{})
}
