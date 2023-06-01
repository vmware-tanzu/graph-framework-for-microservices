/*
Copyright 2015 The Kubernetes Authors.
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

package utils_test

import (
	"api-gw/pkg/utils"
	"testing"
)

func TestPlural(t *testing.T) {
	cases := []struct {
		typeName string
		expected string
	}{
		{
			"I",
			"I",
		},
		{
			"Pod",
			"Pods",
		},
		{
			"Entry",
			"Entries",
		},
		{
			"Bus",
			"Buses",
		},
		{
			"Fizz",
			"Fizzes",
		},
		{
			"Search",
			"Searches",
		},
		{
			"Autograph",
			"Autographs",
		},
		{
			"Dispatch",
			"Dispatches",
		},
		{
			"Earth",
			"Earths",
		},
		{
			"City",
			"Cities",
		},
		{
			"Ray",
			"Rays",
		},
		{
			"Fountain",
			"Fountains",
		},
		{
			"Life",
			"Lives",
		},
		{
			"Leaf",
			"Leaves",
		},
	}
	for _, c := range cases {
		if e, a := c.expected, utils.ToPlural(c.typeName); e != a {
			t.Errorf("Unexpected result from plural namer. Expected: %s, Got: %s", e, a)
		}
	}
}
