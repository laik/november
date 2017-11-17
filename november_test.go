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

package november_test

import "testing"
import "sync"
import "fmt"
import . "github.com/laik/november"

import "reflect"

type moduleTest struct {
	Id    int
	Name  string
	Mux   sync.Mutex
	Bills struct {
		BillId   int
		BillName string
		Amount   float64
	}
	GetValue func() string
}

func newModuleTest() (mt *moduleTest) {
	_mt := new(moduleTest)
	_mt.Id = 1
	_mt.Name = "mongotb"
	_mt.Bills.BillId = 99999
	_mt.Bills.BillName = "mongotb bills"
	_mt.GetValue = func() (tmpString string) {
		tmpString = fmt.Sprintf("[ id:%d,name:%s ]", mt.Id, mt.Name)
		return
	}
	return _mt
}

func (mt moduleTest) Get(key string) (value interface{}) {
	value, ok := Xget(mt, key)
	if !ok {
		fmt.Printf("the object not has attr: %s\n", value)
		return value
	}
	return
}

func (mt *moduleTest) Set(args map[string]interface{}) bool {
	for key, value := range args {
		if ok := Xset(mt, key, value); !ok {
			fmt.Printf("can not set value: %s\n", value)
			return false
		}
	}
	return true
}

func TestCoreInterface(t *testing.T) {
	mt := newModuleTest()
	var xint Xinterface = mt
	t.Logf("type struct %#v", xint)
	if x := xint.Get("Name"); x != nil {
		t.Logf("%#v", x)
	} else {
		t.Fatal(x)
	}

	if x := xint.Get("Id"); x != nil {
		t.Log(x)
	} else {
		t.Fatal(x)
	}

	t.Logf("%#v", mt.GetValue())
	xlist, ok := Xlist(mt)
	if ok {
		t.Logf("%#v", xlist)
	}
}

type Speaker interface {
	Speak() string
	Set(string)
}

type Teacher struct {
	Name string
}

func (this *Teacher) Speak() string {
	return this.Name
}

func (this *Teacher) Set(name string) {
	this.Name = name
}

func TestXnew(t *testing.T) {
	t.Log("test Xnew")
	var s Speaker
	s = &Teacher{"wocao"}
	t.Log(reflect.TypeOf(s).Kind())
	t.Log(Xcall("Speak", s))

}
