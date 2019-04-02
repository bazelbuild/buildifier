/*
Copyright 2018 Bazel contributors. All Rights Reserved.

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

package build

import (
	"testing"
)

func TestIsBuildFilename(t *testing.T) {
	cases := map[string]FileType{
		"BUILD":           TypeBuild,
		"build":           TypeBuild,
		"bUIld":           TypeBuild,
		"BUILD.oss":       TypeBuild,
		"BUILD.bazel":     TypeBuild,
		"BUILD.bazel.oss": TypeBuild,
		"build.bzl":       TypeBzl,
		"build.sky":       TypeDefault,
		"WORKSPACE":       TypeWorkspace,
		"external.BUILD":  TypeBuild,
		"BUILD.external":  TypeBuild,
		"aBUILD":          TypeDefault,
		"thing.sky":       TypeDefault,
		"my.WORKSPACE":    TypeWorkspace,
		"thing.bzl":       TypeBzl,
		"thing.bzl.oss":   TypeBzl,
		"thing.bzl.exe":   TypeDefault,
		"workspace.bazel": TypeWorkspace,
		"workspace.bzl":   TypeBzl,
		"foo.bar":         TypeDefault,
	}
	for name, fileType := range cases {
		res := getFileType(name)
		if res != fileType {
			t.Errorf("isBuildFilename(%q) should be %v but was %v", name, fileType, res)
		}
	}
}
