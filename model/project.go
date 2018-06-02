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

package model

import (
	"github.com/infinitbyte/framework/core/errors"
	"github.com/infinitbyte/framework/core/persist"
	"github.com/infinitbyte/framework/core/util"
	"time"
)

// Project is a definition, include a collection of Host
type Project struct {
	ID          string    `storm:"id,unique" json:"id,omitempty" gorm:"not null;unique;primary_key" index:"id"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Enabled     bool      `json:"enabled"`
	Created     time.Time `json:"created,omitempty"`
	Updated     time.Time `json:"updated,omitempty"`
}

func CreateProject(project *Project) error {
	time := time.Now().UTC()
	project.ID = util.GetUUID()
	project.Created = time
	project.Updated = time
	return persist.Save(project)
}

func UpdateProject(project *Project) error {
	time := time.Now().UTC()
	project.Updated = time
	return persist.Update(project)
}

func DeleteProject(id string) error {
	project := Project{ID: id}
	return persist.Delete(&project)
}

func GetProject(id string) (Project, error) {
	project := Project{}
	project.ID = id
	err := persist.Get(&project)
	if err != nil {
		return project, err
	}
	if len(project.ID) == 0 || project.Updated.IsZero() {
		return project, errors.New("not found," + id)
	}

	return project, err
}

func GetProjectList(from, size int) (int, []Project, error) {
	var projects []Project
	sort := []persist.Sort{}
	sort = append(sort, persist.Sort{Field: "created", SortType: persist.ASC})
	queryO := persist.Query{Sort: &sort, From: from, Size: size}
	err, result := persist.Search(Project{}, &projects, &queryO)
	if err != nil {
		return 0, projects, err
	}
	if result.Result != nil && projects == nil || len(projects) == 0 {
		convertProject(result, &projects)
	}
	return result.Total, projects, err
}

func convertProject(result persist.Result, projects *[]Project) {
	if result.Result == nil {
		return
	}

	t, ok := result.Result.([]interface{})
	if ok {
		for _, i := range t {
			js := util.ToJson(i, false)
			t := Project{}
			util.FromJson(js, &t)
			*projects = append(*projects, t)
		}
	}
}
