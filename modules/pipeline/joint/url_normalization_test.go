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

package joint

import (
	"github.com/infinitbyte/gopa/core/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNormailzeLinks(t *testing.T) {

	context := &model.Context{}
	context.Init()
	task := model.Task{}
	task.Url = "http://elasticsearch.cn/"
	task.OriginalUrl = "http://elasticsearch.cn/"
	task.Depth = 1
	task.Breadth = 1
	task.Host = "elasticsearch.cn"

	context.Set(model.CONTEXT_TASK_URL, task.Url)
	context.Set(model.CONTEXT_TASK_Host, task.Host)
	context.Set(model.CONTEXT_TASK_OriginalUrl, task.OriginalUrl)

	snapshot := model.Snapshot{}
	context.Set(model.CONTEXT_SNAPSHOT, &snapshot)

	parse := UrlNormalizationJoint{}
	parse.Process(context)
	assert.Equal(t, "http://elasticsearch.cn/", context.MustGetString(model.CONTEXT_TASK_URL))

	task = model.Task{}
	task.Url = "index.html"
	task.Depth = 1
	task.Breadth = 1
	task.Host = "elasticsearch.cn"

	context = &model.Context{}
	context.Set(model.CONTEXT_TASK_URL, task.Url)
	context.Set(model.CONTEXT_TASK_Host, task.Host)
	snapshot = model.Snapshot{}
	context.Set(model.CONTEXT_SNAPSHOT, &snapshot)

	parse.Process(context)
	assert.Equal(t, "http://elasticsearch.cn/index.html", context.MustGetString(model.CONTEXT_TASK_URL))

	task = model.Task{}
	task.Url = "/index.html"
	task.Depth = 1
	task.Breadth = 1
	task.Host = "elasticsearch.cn"

	context = &model.Context{}
	context.Set(model.CONTEXT_TASK_URL, task.Url)
	context.Set(model.CONTEXT_TASK_Host, task.Host)
	snapshot = model.Snapshot{}
	context.Set(model.CONTEXT_SNAPSHOT, &snapshot)

	parse.Process(context)
	assert.Equal(t, "http://elasticsearch.cn/index.html", context.MustGetString(model.CONTEXT_TASK_URL))
}

func TestNormailzeLinks1(t *testing.T) {

	context := &model.Context{}
	context.Init()
	task := model.Task{}
	task.Url = "http://localhost/"
	task.Depth = 1
	task.Breadth = 1

	context.Set(model.CONTEXT_TASK_URL, task.Url)
	context.Set(model.CONTEXT_TASK_Host, task.Host)
	snapshot := model.Snapshot{}
	context.Set(model.CONTEXT_SNAPSHOT, &snapshot)

	parse := UrlNormalizationJoint{}
	parse.Process(context)
	assert.Equal(t, "http://localhost/", context.MustGetString(model.CONTEXT_TASK_URL))

	task = model.Task{}
	task.Url = "http://localhost/index.html"
	task.Depth = 1
	task.Breadth = 1

	context.Set(model.CONTEXT_TASK_URL, task.Url)
	context.Set(model.CONTEXT_TASK_Host, task.Host)
	snapshot = model.Snapshot{}
	context.Set(model.CONTEXT_SNAPSHOT, &snapshot)

	parse.Process(context)
	assert.Equal(t, "http://localhost/index.html", context.MustGetString(model.CONTEXT_TASK_URL))

	task = model.Task{}
	task.Url = "http://localhost:8080/index.html"
	task.Depth = 1
	task.Breadth = 1

	context.Set(model.CONTEXT_TASK_URL, task.Url)
	context.Set(model.CONTEXT_TASK_Host, task.Host)
	snapshot = model.Snapshot{}
	context.Set(model.CONTEXT_SNAPSHOT, &snapshot)

	parse.Process(context)
	assert.Equal(t, "http://localhost:8080/index.html", context.MustGetString(model.CONTEXT_TASK_URL))

	task = model.Task{}
	task.Url = "phpliteadmin.php?table=groupes&action=row_editordelete&pk=3&type=edit"
	task.Reference = "http://localhost:8080/index.html"
	task.Depth = 1
	task.Breadth = 1

	context.Set(model.CONTEXT_TASK_URL, task.Url)
	context.Set(model.CONTEXT_TASK_Reference, task.Reference)
	context.Set(model.CONTEXT_TASK_Host, task.Host)
	snapshot = model.Snapshot{}
	context.Set(model.CONTEXT_SNAPSHOT, &snapshot)

	parse.Process(context)
	assert.Equal(t, "http://localhost:8080/phpliteadmin.php?table=groupes&action=row_editordelete&pk=3&type=edit", context.MustGetString(model.CONTEXT_TASK_URL))

}

func TestNormailzeLinks2(t *testing.T) {

	context := &model.Context{}
	context.Init()

	task := model.Task{}
	task.Url = "http://127.0.0.1:8080/modeling-your-data.html"
	task.Depth = 1
	task.Breadth = 1

	context.Set(model.CONTEXT_TASK_URL, task.Url)
	context.Set(model.CONTEXT_TASK_Host, task.Host)
	snapshot := model.Snapshot{}
	context.Set(model.CONTEXT_SNAPSHOT, &snapshot)

	parse := UrlNormalizationJoint{}
	parse.Process(context)

	assert.Equal(t, "http://127.0.0.1:8080/modeling-your-data.html", context.MustGetString(model.CONTEXT_TASK_URL))
	assert.Equal(t, "", snapshot.Path)
	assert.Equal(t, "/modeling-your-data.html", snapshot.File)

	task = model.Task{}
	task.Url = "http://127.0.0.1:8080/video"
	task.Depth = 1
	task.Breadth = 1

	context.Set(model.CONTEXT_TASK_URL, task.Url)
	context.Set(model.CONTEXT_TASK_Host, task.Host)
	snapshot = model.Snapshot{}
	context.Set(model.CONTEXT_SNAPSHOT, &snapshot)

	parse.Process(context)
	assert.Equal(t, "/video", snapshot.Path)
	assert.Equal(t, "default.html", snapshot.File)
}

func TestNormailzeLinks3(t *testing.T) {

	context := &model.Context{}
	context.Init()
	parse := UrlNormalizationJoint{}
	task := model.Task{}
	task.Url = "http://conf.elasticsearch.cn/2015/../2015/../2015/../2015/../2015/../2015/../2015/../2015/../2015/../2015/../2015/../2015/../2015/../2015/../2015/../2015/../2015/../2015/../2015/../2015/../2015/../2015/../2015/../2015/beijing.html?c=3&a=1&b=9&c=0#targetsa"
	task.Depth = 1
	task.Breadth = 1

	context.Set(model.CONTEXT_TASK_URL, task.Url)
	context.Set(model.CONTEXT_TASK_Host, task.Host)
	snapshot := model.Snapshot{}
	context.Set(model.CONTEXT_SNAPSHOT, &snapshot)

	parse.Process(context)

	assert.Equal(t, "http://conf.elasticsearch.cn/2015/beijing.html?c=3&a=1&b=9&c=0", context.MustGetString(model.CONTEXT_TASK_URL))
	assert.Equal(t, "/2015", snapshot.Path)
	assert.Equal(t, "/beijing_a_1_b_9_c_3_0.html", snapshot.File)
}

func TestNormailzeLinks4(t *testing.T) {

	context := &model.Context{}
	context.Init()
	parse := UrlNormalizationJoint{}

	task := model.Task{}
	task.Url = "../2015/chengdu.html"
	task.Reference = "http://conf.elasticsearch.cn/2015/chengdu.html"
	task.Depth = 1
	task.Breadth = 1

	context.Set(model.CONTEXT_TASK_URL, task.Url)
	context.Set(model.CONTEXT_TASK_Reference, task.Reference)
	context.Set(model.CONTEXT_TASK_Host, task.Host)
	snapshot := model.Snapshot{}
	context.Set(model.CONTEXT_SNAPSHOT, &snapshot)

	parse.Process(context)

	assert.Equal(t, "http://conf.elasticsearch.cn/2015/chengdu.html", context.MustGetString(model.CONTEXT_TASK_URL))
	assert.Equal(t, "/2015", snapshot.Path)
	assert.Equal(t, "/chengdu.html", snapshot.File)
}

func TestNormailzeLinks5(t *testing.T) {

	context := &model.Context{}
	context.Init()
	parse := UrlNormalizationJoint{}

	task := model.Task{}
	task.Url = "../../../../articles/%E4%BF%A1/%E4%B9%89/%E5%AE%97/%E4%BF%A1%E4%B9%89%E5%AE%97.html"
	task.Reference = "http://wiki.example.org/articles/%E8%82%AF/%E5%A1%94/%E5%9F%BA/%E8%82%AF%E5%A1%94%E5%9F%BA%E5%B7%9E.html"
	task.Depth = 1
	task.Breadth = 1

	context.Set(model.CONTEXT_TASK_URL, task.Url)
	context.Set(model.CONTEXT_TASK_Reference, task.Reference)
	context.Set(model.CONTEXT_TASK_Host, task.Host)
	snapshot := model.Snapshot{}
	context.Set(model.CONTEXT_SNAPSHOT, &snapshot)

	parse.Process(context)

	assert.Equal(t, "http://wiki.example.org/articles/%E4%BF%A1/%E4%B9%89/%E5%AE%97/%E4%BF%A1%E4%B9%89%E5%AE%97.html", context.MustGetString(model.CONTEXT_TASK_URL))
}
