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

package persist

import "errors"

type ORM interface {
	Save(o interface{}) error

	Update(o interface{}) error

	Delete(o interface{}) error

	Search(t interface{}, to interface{}, q *Query) (error, Result)

	Get(o interface{}) error

	GetBy(field string, value interface{}, t interface{}, to interface{}) (error, Result)

	Count(o interface{}) (int, error)

	GroupBy(o interface{},field string)(error,map[string]interface{})
}

type Sort struct {
	Field    string
	SortType SortType
}

type SortType string

const ASC SortType = "asc"
const DESC SortType = "desc"

type Query struct {
	Sort     *[]Sort
	From     int
	Size     int
	Conds    []*Cond
	RawQuery string
}

type Cond struct {
	Field       string
	SQLOperator string
	QueryType   QueryType
	BoolType    BoolType
	Value       interface{}
}

type BoolType string
type QueryType string

const Must BoolType = "must"
const MustNot BoolType = "must_not"
const Should BoolType = "should"

const Match QueryType = "match"
const RangeGt QueryType = "gt"
const RangeGte QueryType = "gte"
const RangeLt QueryType = "lt"
const RangeLte QueryType = "lte"

func Eq(field string, value interface{}) *Cond {
	c := Cond{}
	c.Field = field
	c.Value = value
	c.SQLOperator = " = "
	c.QueryType = Match
	c.BoolType = Must
	return &c
}

func NotEq(field string, value interface{}) *Cond {
	c := Cond{}
	c.Field = field
	c.Value = value
	c.SQLOperator = " != "
	c.QueryType = Match
	c.BoolType = MustNot
	return &c
}

func Gt(field string, value interface{}) *Cond {
	c := Cond{}
	c.Field = field
	c.Value = value
	c.SQLOperator = " > "
	c.QueryType = RangeGt
	c.BoolType = Must
	return &c
}

func Lt(field string, value interface{}) *Cond {
	c := Cond{}
	c.Field = field
	c.Value = value
	c.SQLOperator = " < "
	c.QueryType = RangeLt
	c.BoolType = Must
	return &c
}

func Ge(field string, value interface{}) *Cond {
	c := Cond{}
	c.Field = field
	c.Value = value
	c.SQLOperator = " >= "
	c.QueryType = RangeGte
	c.BoolType = Must
	return &c
}

func Le(field string, value interface{}) *Cond {
	c := Cond{}
	c.Field = field
	c.Value = value
	c.SQLOperator = " <= "
	c.QueryType = RangeLte
	c.BoolType = Must
	return &c
}

func Combine(conds ...[]*Cond) []*Cond {
	t := []*Cond{}
	for _, cs := range conds {
		for _, c := range cs {
			t = append(t, c)
		}
	}
	return t
}

func And(conds ...*Cond) []*Cond {
	t := []*Cond{}
	for _, c := range conds {
		t = append(t, c)
	}
	return t
}

type Result struct {
	Total  int
	Result interface{}
}

func GetBy(field string, value interface{}, t interface{}, to interface{}) (error, Result) {

	return getHandler().GetBy(field, value, t, to)
}

func Get(o interface{}) error {
	return getHandler().Get(o)
}

func Save(o interface{}) error {

	return getHandler().Save(o)
}

func Update(o interface{}) error {
	return getHandler().Update(o)
}

func Delete(o interface{}) error {
	return getHandler().Delete(o)
}

func Count(o interface{}) (int, error) {
	return getHandler().Count(o)
}

func Search(t interface{}, to interface{}, q *Query) (error, Result) {
	return getHandler().Search(t, to, q)
}

func GroupBy(o interface{},field string)(error,map[string]interface{}) {
	return getHandler().GroupBy(o, field)
}

var handler ORM

func getHandler() ORM {
	if handler == nil {
		panic(errors.New("store handler is not registered"))
	}
	return handler
}

func Register(h ORM) {
	handler = h
}
