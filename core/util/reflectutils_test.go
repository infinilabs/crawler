package util

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestCloneValue(t *testing.T) {
	data := map[string]interface{}{}
	v := map[string]interface{}{}
	v["key"] = "value"
	data["name"] = "name-1"
	data["data"] = v

	type T struct {
		Name string `json:"name"`
		Data struct {
			Key string `json:"key"`
		}
	}

	js := ToJson(data, false)
	o := T{}
	FromJson(js, &o)

	assert.Equal(t, "name-1", o.Name)
	assert.Equal(t, "value", o.Data.Key)
}

func TestGetTag(t *testing.T) {

	type S struct {
		F string `species:"gopher" color:"blue"`
	}

	ts := S{F: "test F"}
	st := reflect.TypeOf(ts)

	field := st.Field(0)
	fmt.Println(field.Tag.Get("color"), field.Tag.Get("species"))

	fmt.Println(field.Name)

	fmt.Println(reflect.Indirect(reflect.ValueOf(ts)).FieldByName(field.Name).String())

	v := GetFieldValueByTagName(ts, "color", "blue")
	fmt.Println(v)
	assert.Equal(t, v, "test F")

	vs := &S{"flower"}

	fmt.Println(reflect.TypeOf(vs))
	fmt.Println(reflect.ValueOf(vs))
	fmt.Println(reflect.Indirect(reflect.ValueOf(vs)).Type().Name())

	se := reflect.TypeOf(vs).Elem()
	for i := 0; i < se.NumField(); i++ {
		fmt.Println(se.Field(i).Name)
		fmt.Println(se.Field(i).Type)
		fmt.Println(se.Field(i).Tag)
	}

	v1 := GetFieldValueByTagName(vs, "color", "blue")
	fmt.Println(v1)
	assert.Equal(t, v1, "flower")

	fmt.Println(reflect.TypeOf(ts))
	fmt.Println(reflect.TypeOf(vs))
	fmt.Println(reflect.ValueOf(ts))
	fmt.Println(reflect.ValueOf(vs))

	assert.Equal(t, GetTypeName(ts, false), "S")
	assert.Equal(t, GetTypeName(vs, false), "S")

}
