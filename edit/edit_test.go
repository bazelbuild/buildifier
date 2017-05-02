/*
Copyright 2016 Google Inc. All Rights Reserved.
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
package edit

import (
	"reflect"
	"strings"
	"testing"

	"github.com/bazelbuild/buildifier/build"
)

var parseLabelTests = []struct {
	in   string
	pkg  string
	rule string
}{
	{"//devtools/buildozer:rule", "devtools/buildozer", "rule"},
	{"devtools/buildozer:rule", "devtools/buildozer", "rule"},
	{"//devtools/buildozer", "devtools/buildozer", "buildozer"},
	{"//base", "base", "base"},
	{"//base:", "base", "base"},
}

func TestParseLabel(t *testing.T) {
	for i, tt := range parseLabelTests {
		pkg, rule := ParseLabel(tt.in)
		if pkg != tt.pkg || rule != tt.rule {
			t.Errorf("%d. ParseLabel(%q) => (%q, %q), want (%q, %q)",
				i, tt.in, pkg, rule, tt.pkg, tt.rule)
		}
	}
}

var shortenLabelTests = []struct {
	in     string
	pkg    string
	result string
}{
	{"//devtools/buildozer:rule", "devtools/buildozer", ":rule"},
	{"//devtools/buildozer:rule", "devtools", "//devtools/buildozer:rule"},
	{"//base:rule", "devtools", "//base:rule"},
	{"//base:base", "devtools", "//base"},
	{"//base", "base", ":base"},
	{":local", "", ":local"},
	{"something else", "", "something else"},
	{"/path/to/file", "path/to", "/path/to/file"},
}

func TestShortenLabel(t *testing.T) {
	for i, tt := range shortenLabelTests {
		result := ShortenLabel(tt.in, tt.pkg)
		if result != tt.result {
			t.Errorf("%d. ShortenLabel(%q, %q) => %q, want %q",
				i, tt.in, tt.pkg, result, tt.result)
		}
	}
}

var labelsEqualTests = []struct {
	label1   string
	label2   string
	pkg      string
	expected bool
}{
	{"//devtools/buildozer:rule", "rule", "devtools/buildozer", true},
	{"//devtools/buildozer:rule", "rule:jar", "devtools", false},
}

func TestLabelsEqual(t *testing.T) {
	for i, tt := range labelsEqualTests {
		if got := LabelsEqual(tt.label1, tt.label2, tt.pkg); got != tt.expected {
			t.Errorf("%d. LabelsEqual(%q, %q, %q) => %v, want %v",
				i, tt.label1, tt.label2, tt.pkg, got, tt.expected)
		}
	}
}

var splitOnSpacesTests = []struct {
	in  string
	out []string
}{
	{"a", []string{"a"}},
	{"  abc def ", []string{"abc", "def"}},
	{`  abc\ def `, []string{"abc def"}},
}

func TestSplitOnSpaces(t *testing.T) {
	for i, tt := range splitOnSpacesTests {
		result := SplitOnSpaces(tt.in)
		if !reflect.DeepEqual(result, tt.out) {
			t.Errorf("%d. SplitOnSpaces(%q) => %q, want %q",
				i, tt.in, result, tt.out)
		}
	}
}

func TestInsertLoad(t *testing.T) {
	tests := []struct{ input, expected string }{
		{``, `load("location", "symbol")`},
		{`load("location", "symbol")`, `load("location", "symbol")`},
		{`load("location", "other", "symbol")`, `load("location", "other", "symbol")`},
		{`load("location", "other")`, `load("location", "other", "symbol")`},
		{
			`load("other loc", "symbol")`,
			`load("location", "symbol")
load("other loc", "symbol")`,
		},
	}

	for _, tst := range tests {
		bld, err := build.Parse("BUILD", []byte(tst.input))
		if err != nil {
			t.Error(err)
			continue
		}
		bld.Stmt = InsertLoad(bld.Stmt, []string{"location", "symbol"})
		got := strings.TrimSpace(string(build.Format(bld)))
		if got != tst.expected {
			t.Errorf("maybeInsertLoad(%s): got %s, expected %s", tst.input, got, tst.expected)
		}
	}
}
