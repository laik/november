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
import "strconv"
import "strings"

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
	if rv.Type().Kind() == reflect.Ptr && !rv.IsNil() {
		rv = rv.Elem()
		if !rv.CanSet() {
			fmt.Errorf(`type can not set field:"%s",value:"%s" new struct:%#v`+"\n", key, value, t)
			return
		}
	}
	field := rv.FieldByName(key)
	field.Set(reflect.ValueOf(value))
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

//unmarsha
func XunmarshaText(obj interface{}, data string, _func func(s string) ([]string, error)) (ok bool) {
	text, err := _func(data)
	if err != nil {
		fmt.Errorf("can not split text")
		return false
	}
	rv := reflect.ValueOf(obj)
	if rv.Type().Kind() == reflect.Ptr && !rv.IsNil() {
		rv = rv.Elem()
		if !rv.CanSet() {
			return false
		}
	}
	if len(text) != rv.NumField() {
		return false
	}
	for i := 0; i < rv.NumField(); i++ {
		kind := rv.Field(i).Type().Kind()
		value := text[i]
		switch kind {
		case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
			switch kind {
			case reflect.Int8:
				out, err := strconv.ParseInt(value, 10, 8)
				if err != nil {
					fmt.Errorf("string[%s] covert Int8 fail. %s", value, err)
				}
				rv.Field(i).Set(reflect.ValueOf(int8(out)))
			case reflect.Int16:
				out, err := strconv.ParseInt(value, 10, 16)
				if err != nil {
					fmt.Errorf("string[%s] covert Int16 fail. %s", value, err)
				}
				rv.Field(i).Set(reflect.ValueOf(int16(out)))
			case reflect.Int32:
				out, err := strconv.ParseInt(value, 10, 32)
				if err != nil {
					fmt.Errorf("string[%s] covert Int32 fail. %s", value, err)
				}
				rv.Field(i).Set(reflect.ValueOf(int32(out)))
			case reflect.Int64:
				out, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					fmt.Errorf("string[%s] covert int64 fail. %s", value, err)
				}
				rv.Field(i).SetInt(out)
			default:
				out, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					fmt.Errorf("string[%s] covert int fail. %s", value, err)
				}
				rv.Field(i).Set(reflect.ValueOf(int(out)))
			}
		case reflect.Float32, reflect.Float64:
			switch kind {
			case reflect.Float32:
				if val, err := strconv.ParseFloat(value, 32); err == nil {
					rv.Field(i).Set(reflect.ValueOf(float32(val)))
				}
			default:
				if val, err := strconv.ParseFloat(value, 64); err == nil {
					rv.Field(i).Set(reflect.ValueOf(val))
				}
			}
		case reflect.Bool:
			var tmp bool
			if value == "Y" || strings.ToUpper(value) == "YES" {
				tmp = true
			}
			rv.Field(i).SetBool(tmp)
		case reflect.String:
			rv.Field(i).SetString(value)
		}
	}
	return true
}
