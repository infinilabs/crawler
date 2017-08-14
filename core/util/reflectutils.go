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

package util

import (
	"errors"
	"reflect"
	"strings"
)

// Invoke dynamic execute function via function name and parameters
func Invoke(any interface{}, name string, args ...interface{}) {
	inputs := make([]reflect.Value, len(args))
	for i := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	reflect.ValueOf(any).MethodByName(name).Call(inputs)
}

// GetFieldValueByTagName return the field value which field was tagged with this tagName, only support string field
func GetFieldValueByTagName(any interface{}, tagName string, tagValue string) string {

	t := reflect.TypeOf(any)
	if PrefixStr(t.String(), "*") {
		se := reflect.TypeOf(any).Elem()

		for i := 0; i < se.NumField(); i++ {
			v := se.Field(i).Tag.Get(tagName)
			if v != "" {
				if v == tagValue {
					return reflect.Indirect(reflect.ValueOf(any)).FieldByName(se.Field(i).Name).String()
				}
			}

		}
	}

	for i := 0; i < t.NumField(); i++ {
		v := t.Field(i).Tag.Get(tagName)
		if v != "" {
			if v == tagValue {
				return reflect.Indirect(reflect.ValueOf(any)).FieldByName(t.Field(i).Name).String()
			}
		}

	}

	panic(errors.New("tag was not found"))
}

func GetTypeName(any interface{}, lowercase bool) string {
	name := reflect.Indirect(reflect.ValueOf(any)).Type().Name()
	if lowercase {
		name = strings.ToLower(name)
	}
	return name
}
