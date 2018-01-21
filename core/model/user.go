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
	"github.com/emirpasic/gods/sets/hashset"
)

const (
	ROLE_GUEST string = "guest"
	ROLE_ADMIN string = "admin"
)

const (
	//GUEST
	PERMISSION_SNAPSHOT_VIEW string = "view_snapshot"

	//ADMIN
	PERMISSION_ADMIN_MINIMAL string = "admin_minimal"
)

func GetPermissionsByRole(role string) (*hashset.Set, error) {
	initRolesMap()
	return rolesMap[role], nil
}

var rolesMap = map[string]*hashset.Set{}

func initRolesMap() {
	if rolesMap != nil {
		return
	}
	set := hashset.New()
	set.Add(PERMISSION_SNAPSHOT_VIEW)
	rolesMap[ROLE_GUEST] = set
}
