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
	"github.com/infinitbyte/gopa/core/model"
	"github.com/infinitbyte/gopa/core/pipeline"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNormailzeLinks(t *testing.T) {

	context := &pipeline.Context{}
	context.Init()
	task := model.Task{}
	task.Url = "http://elasticsearch.cn/"
	task.OriginalUrl = "http://elasticsearch.cn/"
	task.Depth = 1
	task.Breadth = 1
	task.Host = "elasticsearch.cn"

	context.Set(CONTEXT_CRAWLER_TASK, &task)

	snapshot := model.Snapshot{}
	context.Set(CONTEXT_CRAWLER_SNAPSHOT, &snapshot)

	parse := UrlNormalizationJoint{}
	parse.Process(context)
	assert.Equal(t, "http://elasticsearch.cn/", task.Url)

	task = model.Task{}
	task.Url = "index.html"
	task.Depth = 1
	task.Breadth = 1
	task.Host = "elasticsearch.cn"

	context.Set(CONTEXT_CRAWLER_TASK, &task)
	snapshot = model.Snapshot{}
	context.Set(CONTEXT_CRAWLER_SNAPSHOT, &snapshot)

	parse.Process(context)
	assert.Equal(t, "http://elasticsearch.cn/index.html", task.Url)

	task = model.Task{}
	task.Url = "/index.html"
	task.Depth = 1
	task.Breadth = 1
	task.Host = "elasticsearch.cn"

	context.Set(CONTEXT_CRAWLER_TASK, &task)
	snapshot = model.Snapshot{}
	context.Set(CONTEXT_CRAWLER_SNAPSHOT, &snapshot)

	parse.Process(context)
	assert.Equal(t, "http://elasticsearch.cn/index.html", task.Url)
}

func TestNormailzeLinks1(t *testing.T) {

	context := &pipeline.Context{}
	context.Init()
	task := model.Task{}
	task.Url = "http://localhost/"
	task.Depth = 1
	task.Breadth = 1

	context.Set(CONTEXT_CRAWLER_TASK, &task)
	snapshot := model.Snapshot{}
	context.Set(CONTEXT_CRAWLER_SNAPSHOT, &snapshot)

	parse := UrlNormalizationJoint{}
	parse.Process(context)
	assert.Equal(t, "http://localhost/", task.Url)

	task = model.Task{}
	task.Url = "http://localhost/index.html"
	task.Depth = 1
	task.Breadth = 1

	context.Set(CONTEXT_CRAWLER_TASK, &task)
	snapshot = model.Snapshot{}
	context.Set(CONTEXT_CRAWLER_SNAPSHOT, &snapshot)

	parse.Process(context)
	assert.Equal(t, "http://localhost/index.html", task.Url)

	task = model.Task{}
	task.Url = "http://localhost:8080/index.html"
	task.Depth = 1
	task.Breadth = 1

	context.Set(CONTEXT_CRAWLER_TASK, &task)
	snapshot = model.Snapshot{}
	context.Set(CONTEXT_CRAWLER_SNAPSHOT, &snapshot)

	parse.Process(context)
	assert.Equal(t, "http://localhost:8080/index.html", task.Url)

	task = model.Task{}
	task.Url = "phpliteadmin.php?table=groupes&action=row_editordelete&pk=3&type=edit"
	task.Reference = "http://localhost:8080/index.html"
	task.Depth = 1
	task.Breadth = 1

	context.Set(CONTEXT_CRAWLER_TASK, &task)
	snapshot = model.Snapshot{}
	context.Set(CONTEXT_CRAWLER_SNAPSHOT, &snapshot)

	parse.Process(context)
	assert.Equal(t, "http://localhost:8080/phpliteadmin.php?table=groupes&action=row_editordelete&pk=3&type=edit", task.Url)

}

func TestNormailzeLinks2(t *testing.T) {

	context := &pipeline.Context{}
	context.Init()

	task := model.Task{}
	task.Url = "http://127.0.0.1:8080/modeling-your-data.html"
	task.Depth = 1
	task.Breadth = 1

	context.Set(CONTEXT_CRAWLER_TASK, &task)
	snapshot := model.Snapshot{}
	context.Set(CONTEXT_CRAWLER_SNAPSHOT, &snapshot)

	parse := UrlNormalizationJoint{}
	parse.Process(context)

	assert.Equal(t, "http://127.0.0.1:8080/modeling-your-data.html", task.Url)
	assert.Equal(t, "", snapshot.Path)
	assert.Equal(t, "/modeling-your-data.html", snapshot.File)

	task = model.Task{}
	task.Url = "http://127.0.0.1:8080/video"
	task.Depth = 1
	task.Breadth = 1

	context.Set(CONTEXT_CRAWLER_TASK, &task)
	snapshot = model.Snapshot{}
	context.Set(CONTEXT_CRAWLER_SNAPSHOT, &snapshot)

	parse.Process(context)
	assert.Equal(t, "/video", snapshot.Path)
	assert.Equal(t, "default.html", snapshot.File)
}

func TestNormailzeLinks3(t *testing.T) {

	context := &pipeline.Context{}
	context.Init()
	parse := UrlNormalizationJoint{}
	task := model.Task{}
	task.Url = "http://conf.elasticsearch.cn/2015/../2015/../2015/../2015/../2015/../2015/../2015/../2015/../2015/../2015/../2015/../2015/../2015/../2015/../2015/../2015/../2015/../2015/../2015/../2015/../2015/../2015/../2015/../2015/beijing.html?c=3&a=1&b=9&c=0#targetsa"
	task.Depth = 1
	task.Breadth = 1

	context.Set(CONTEXT_CRAWLER_TASK, &task)
	snapshot := model.Snapshot{}
	context.Set(CONTEXT_CRAWLER_SNAPSHOT, &snapshot)

	parse.Process(context)

	assert.Equal(t, "http://conf.elasticsearch.cn/2015/beijing.html?c=3&a=1&b=9&c=0", task.Url)
	assert.Equal(t, "/2015", snapshot.Path)
	assert.Equal(t, "/beijing_a_1_b_9_c_3_0.html", snapshot.File)
}

func TestNormailzeLinks4(t *testing.T) {

	context := &pipeline.Context{}
	context.Init()
	parse := UrlNormalizationJoint{}

	task := model.Task{}
	task.Url = "../2015/chengdu.html"
	task.Reference = "http://conf.elasticsearch.cn/2015/chengdu.html"
	task.Depth = 1
	task.Breadth = 1

	context.Set(CONTEXT_CRAWLER_TASK, &task)
	snapshot := model.Snapshot{}
	context.Set(CONTEXT_CRAWLER_SNAPSHOT, &snapshot)

	parse.Process(context)

	assert.Equal(t, "http://conf.elasticsearch.cn/2015/chengdu.html", task.Url)
	assert.Equal(t, "/2015", snapshot.Path)
	assert.Equal(t, "/chengdu.html", snapshot.File)
}

func TestNormailzeLinks5(t *testing.T) {

	context := &pipeline.Context{}
	context.Init()
	parse := UrlNormalizationJoint{}

	task := model.Task{}
	task.Url = "../../../../articles/%E4%BF%A1/%E4%B9%89/%E5%AE%97/%E4%BF%A1%E4%B9%89%E5%AE%97.html"
	task.Reference = "http://wiki.example.org/articles/%E8%82%AF/%E5%A1%94/%E5%9F%BA/%E8%82%AF%E5%A1%94%E5%9F%BA%E5%B7%9E.html"
	task.Depth = 1
	task.Breadth = 1

	context.Set(CONTEXT_CRAWLER_TASK, &task)
	snapshot := model.Snapshot{}
	context.Set(CONTEXT_CRAWLER_SNAPSHOT, &snapshot)

	parse.Process(context)

	assert.Equal(t, "http://wiki.example.org/articles/%E4%BF%A1/%E4%B9%89/%E5%AE%97/%E4%BF%A1%E4%B9%89%E5%AE%97.html", task.Url)
}
