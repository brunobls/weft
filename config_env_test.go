// Copyright 2026 The Weft Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package weft

import "testing"

func TestConfigFileFromEnv(t *testing.T) {
	for _, tc := range []struct{ weft, legacy, want string }{
		{"", "", ""},
		{"a.toml", "", "a.toml"},
		{"", "old.toml", "old.toml"},
		{"a.toml", "old.toml", "a.toml"}, // WEFT_CONFIG wins
	} {
		t.Setenv("WEFT_CONFIG", tc.weft)
		t.Setenv("SERVICEWEAVER_CONFIG", tc.legacy)
		if got := configFileFromEnv(); got != tc.want {
			t.Errorf("WEFT_CONFIG=%q SERVICEWEAVER_CONFIG=%q: got %q, want %q", tc.weft, tc.legacy, got, tc.want)
		}
	}
}
