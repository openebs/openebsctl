/*
Copyright 2020-2021 The OpenEBS Authors

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

/*
A hacky set of tests to ensure cas-types don't get missed out in either
the CasList or the CasMap. This makes logical errors fail builds, to catch
more mistakes in build time than in run-time.
*/
package volume

import (
	"testing"
)

const supportedCasTypeCount = 2

// TestCasList is a dummy test which ensures that each cas-type volumes can be
// listed individually as well as collectively
func TestCasList(t *testing.T) {
	if got := CasList(); len(got) != supportedCasTypeCount {
		t.Fatalf("mismatched number of supported cas-types in the function list, got: %d, expected: %d",
			len(got), supportedCasTypeCount)
	}
}

// TestCasMap is a dummy test which ensures that each cas-type volumes can be
// listed individually as well as collectively
func TestCasMap(t *testing.T) {
	if got := CasMap(); len(got) != supportedCasTypeCount {
		t.Fatalf("mismatched number of supported cas-types in the function map, got: %d, expected: %d",
			len(got), supportedCasTypeCount)
	}
}

// TestCasMap_EmptyAbsent checks if "" cas-type hasn't been added, doing so
// will introduce a wrong-implementation to fetch all volumes across cas-types
func TestCasMap_EmptyAbsent(t *testing.T) {
	if _, ok := CasMap()[""]; ok {
		t.Fatalf("\"\" is not a valid cas-type, please remove it, it'll break some logic")
	}
}
