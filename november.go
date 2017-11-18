// Copyright 2017 The laik Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations

// under the License.

package november

import "reflect"
import "fmt"

type Xinterface interface {
	Ixget
	Ixset
}

type Ixget interface {
	Get(key string) interface{}
}

type Ixset interface {
	Set(value map[string]interface{}) bool
}

func Xlist(t interface{}) (field []string, ok bool) {
	// enum instance field
	field = make([]string, 0)
	//reflect type to struct
	rt := reflect.TypeOf(t)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	for i := 0; i < rt.NumField(); i++ {
		field = append(field, rt.Field(i).Name)
	}
	ok = true
	return
}

func Xget(t interface{}, value string) (news interface{}, ok bool) {
	//check object hasattr
	if fields, xok := Xlist(t); xok {
		for idx, _ := range fields {
			if fields[idx] == value {
				goto TODO
			}
		}
		fmt.Errorf("not found %s in %#v", value, t)
		return
	}
TODO:
	//getter instance attr values; like python getattr
	rt := reflect.TypeOf(t)

	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	//get field value by reflect value
	for i := 0; i < rt.NumField(); i++ {
		if rt.Field(i).Name == value {
			rv := reflect.Indirect(reflect.ValueOf(t))
			return rv.FieldByName(value).Interface(), true
		}
	}
	return
}

func Xset(t interface{}, key string, value interface{}) (ok bool) {
	//check object hasattr
	if fields, xok := Xlist(t); xok {
		for idx, _ := range fields {
			if fields[idx] == key {
				goto TODO
			}
		}
		fmt.Errorf("not found %s in %#v", value, t)
		return false
	}
TODO:
	//setter instance attr values; like python Setattr
	rv := reflect.ValueOf(t)

	if rv.Type().Kind() == reflect.Ptr {
		rv = rv.Elem()
		//set field value
		if !rv.CanSet() {
			fmt.Errorf(`type can not set field:"%s",value:"%s" new struct:%#v`+"\n", key, value, t)
			return
		}
	}
	f := rv.FieldByName(key)
	f.Set(reflect.ValueOf(value))
	return true
}

func Xcall(method string, object interface{}, args ...interface{}) ([]reflect.Value, error) {
	rv := reflect.ValueOf(object)

	if len(args) > 0 {
		input := make([]reflect.Value, len(args))
		for i, _ := range args {
			input[i] = reflect.ValueOf(args[i])
		}
		return rv.MethodByName(method).Call(input), nil
	}
	return rv.MethodByName(method).Call([]reflect.Value{}), nil
}

func xType(object interface{},field string) reflect.Type{
	rv := reflect.ValueOf(object)

	if rv.Type().Kind() == reflect.Ptr {
		rv = rv.Elem()
		//set field value
		if !rv.CanSet() {
			fmt.Errorf(`type can not set field:"%s",value:"%s" new struct:%#v`+"\n", key, value, t)
			return
		}
	}
	return rv.FieldByName(field).Type()
}

func Xformat(object interface{}, data string, split func(string) ([]string, error)) bool {
	fields, ok := Xlist(object)
	if !ok {
		return false
	}
	values, err := split(data)
	if err != nil {
		fmt.Errorf("cant not split data %s", data)
		return false
	}
	if len(fields) != len(values) {
		return false
	}
	for idx, field := range fields {
		Xset(object, field, values[idx].(xType(object,field)))

	}
	return true
}
