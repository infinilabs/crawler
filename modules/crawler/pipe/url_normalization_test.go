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
	"testing"
	"github.com/medcl/gopa/core/env"
	"github.com/medcl/gopa/core/pipeline"
	"github.com/stretchr/testify/assert"
)

func TestNormailzeLinks(t *testing.T) {

	context:= &pipeline.Context{Env:env.EmptyEnv()}
	context.data =map[pipeline.ContextKey]interface{}{}
	context.Set(CONTEXT_URL,"http://elasticsearch.cn/")
	context.Set(CONTEXT_REFERENCE_URL,"http://elasticsearch.cn/")
	context.Set(CONTEXT_DEPTH,1)
	parse:=UrlNormalizationJoint{}
	parse.Process(context)
	assert.Equal(t,"http://elasticsearch.cn/",context.MustGetString(CONTEXT_URL))

	context.Set(CONTEXT_URL,"index.html")
	parse.Process(context)
	assert.Equal(t,"http://elasticsearch.cn/index.html",context.MustGetString(CONTEXT_URL))


	context.Set(CONTEXT_URL,"/index.html")
	parse.Process(context)
	assert.Equal(t,"http://elasticsearch.cn/index.html",context.MustGetString(CONTEXT_URL))
}

func TestNormailzeLinks1(t *testing.T) {

	context:= &pipeline.Context{Env:env.EmptyEnv()}
	context.data =map[pipeline.ContextKey]interface{}{}
	context.Set(CONTEXT_URL,"http://localhost/")
	context.Set(CONTEXT_DEPTH,1)
	parse:=UrlNormalizationJoint{}
	parse.Process(context)
	assert.Equal(t,"http://localhost/",context.MustGetString(CONTEXT_URL))

	context.Set(CONTEXT_URL,"http://localhost/index.html")
	parse.Process(context)
	assert.Equal(t,"http://localhost/index.html",context.MustGetString(CONTEXT_URL))

	context.Set(CONTEXT_URL,"http://localhost:8080/index.html")
	parse.Process(context)
	assert.Equal(t,"http://localhost:8080/index.html",context.MustGetString(CONTEXT_URL))

	context.Set(CONTEXT_URL,"phpliteadmin.php?table=groupes&action=row_editordelete&pk=3&type=edit")
	context.Set(CONTEXT_REFERENCE_URL,"http://localhost:8080/index.html")
	parse.Process(context)
	assert.Equal(t,"http://localhost:8080/phpliteadmin.php?table=groupes&action=row_editordelete&pk=3&type=edit",context.MustGetString(CONTEXT_URL))

}

func TestNormailzeLinks2(t *testing.T) {

	context:= &pipeline.Context{Env:env.EmptyEnv()}
	context.data =map[pipeline.ContextKey]interface{}{}
	context.Set(CONTEXT_URL,"http://127.0.0.1:8080/modeling-your-data.html")
	context.Set(CONTEXT_DEPTH,1)
	parse:=UrlNormalizationJoint{}
	parse.Process(context)

	assert.Equal(t,"http://127.0.0.1:8080/modeling-your-data.html",context.MustGetString(CONTEXT_URL))
	assert.Equal(t,"",context.MustGetString(CONTEXT_SAVE_PATH))
	assert.Equal(t,"/modeling-your-data.html",context.MustGetString(CONTEXT_SAVE_FILENAME))

	context.Set(CONTEXT_URL,"http://127.0.0.1:8080/video")
	parse.Process(context)
	assert.Equal(t,"",context.MustGetString(CONTEXT_SAVE_PATH))
	assert.Equal(t,"/video",context.MustGetString(CONTEXT_SAVE_FILENAME))
}
