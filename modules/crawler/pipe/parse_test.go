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

package pipe

import (
	"github.com/medcl/gopa/core/pipeline"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProcessLinks(t *testing.T) {
	body := "<!DOCTYPE html> <html> <head> <meta content=\"text/html;charset=utf-8\" http-equiv=\"Content-Type\" /> <meta content=\"width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no\" name=\"viewport\" /> <meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge,Chrome=1\" /> <meta name=\"renderer\" content=\"webkit\" /> <title>Elastic中文社区</title> <meta name=\"keywords\" content=\"Elasticsearch中文社区，实时数据分析，实时数据检索, Elastic Stack，ELK，elasticsearch、logstash、kibana、beats等相关技术交流探讨\" /> <meta name=\"description\" content=\"Elasticsearch中文社区，elasticsearch、logstash、kibana,beats等相关技术交流探讨\" /> <base href=\"http://elasticsearch.cn/\" /><!--[if IE]></base><![endif]--> <link href=\"http://elasticsearch.cn/static/css/default/img/favicon.ico?v=20151125\" rel=\"shortcut icon\" type=\"image/x-icon\" /> <link rel=\"stylesheet\" type=\"text/css\" href=\"http://elasticsearch.cn/static/css/bootstrap.css\" /> <link rel=\"stylesheet\" type=\"text/css\" href=\"http://elasticsearch.cn/static/css/icon.css\" /> <link href=\"http://elasticsearch.cn/static/css/default/common.css?v=20151125\" rel=\"stylesheet\" type=\"text/css\" /> <link href=\"http://elasticsearch.cn/static/css/default/link.css?v=20151125\" rel=\"stylesheet\" type=\"text/css\" /> <link href=\"http://elasticsearch.cn/static/js/plug_module/style.css?v=20151125\" rel=\"stylesheet\" type=\"text/css\" /> </head> <body> <div style=\"display:none;\" id=\"__crond\"><a href=\"google.com\">myLink</a>" +
		"<a href=\"//baidu.com\">baidu</a>" +
		"<a href=\"/wiki/Marking/Users\">/wiki/Marking/Users</a>" +
		" </div> </body> </html>"

	context := pipeline.Context{}
	context.Init()
	context.Set(CONTEXT_URL, "http://elasticsearch.cn/")
	context.Set(CONTEXT_DEPTH, 1)
	context.Set(CONTEXT_PAGE_BODY_BYTES, []byte(body))
	parse := ParserJoint{}
	parse.Process(&context)

	links := context.MustGetMap(CONTEXT_PAGE_LINKS)
	println(links["google.com"])
	assert.Equal(t, "baidu", links["//baidu.com"])
	assert.Equal(t, "/wiki/Marking/Users", links["/wiki/Marking/Users"])
	assert.Equal(t, "myLink", links["google.com"])
}
