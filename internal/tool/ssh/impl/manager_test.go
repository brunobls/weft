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


package impl

import "testing"

// TestBabysitterPath checks that the remote binary is looked up under the name
// it was copied with, i.e. the real name of the deployer's executable.
func TestBabysitterPath(t *testing.T) {
	for _, tc := range []struct{ dir, exe, want string }{
		{"/tmp/d", "/usr/bin/weft", "/tmp/d/weft"},
		{"/tmp/d", "/home/u/go/bin/weaver", "/tmp/d/weaver"},
		{"/tmp/d", "/opt/custom-name", "/tmp/d/custom-name"},
	} {
		if got := babysitterPath(tc.dir, tc.exe); got != tc.want {
			t.Errorf("babysitterPath(%q, %q) = %q, want %q", tc.dir, tc.exe, got, tc.want)
		}
	}
}
