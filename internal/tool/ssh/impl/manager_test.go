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

import (
	"slices"
	"testing"
)

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

func TestShellQuote(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{"/tmp/d/weft", "'/tmp/d/weft'"},
		{"/tmp/d/weft;touch X", "'/tmp/d/weft;touch X'"},
		{"/tmp/d/my weft", "'/tmp/d/my weft'"},
		{"/tmp/d/it's", `'/tmp/d/it'\''s'`},
	} {
		if got := shellQuote(tc.in); got != tc.want {
			t.Errorf("shellQuote(%q) = %s, want %s", tc.in, got, tc.want)
		}
	}
}

// TestBabysitterCmd checks the command at its point of use: the binary is
// found under its real name and reaches the remote shell quoted.
func TestBabysitterCmd(t *testing.T) {
	for _, tc := range []struct{ exe, wantPath string }{
		{"/usr/bin/weft", "'/tmp/d/weft'"},
		{"/home/u/bin/weaver", "'/tmp/d/weaver'"},
		{"/x/weft;touch X", "'/tmp/d/weft;touch X'"},
		{"/x/my weft", "'/tmp/d/my weft'"},
	} {
		cmd := babysitterCmd("host", "K=v", "/tmp/d", tc.exe)
		want := []string{"ssh", "host", "K=v", tc.wantPath, "ssh", "babysitter"}
		if !slices.Equal(cmd.Args, want) {
			t.Errorf("exe %q: args %q, want %q", tc.exe, cmd.Args, want)
		}
	}
}
