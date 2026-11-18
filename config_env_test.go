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

import (
	"os"
	"path/filepath"
	"testing"
)

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

// TestLocalOptionsReadsConfig checks the point of use of the env variables: the
// file named by either of them is the one read.
func TestLocalOptionsReadsConfig(t *testing.T) {
	dir := t.TempDir()
	write := func(name, contents string) string {
		f := filepath.Join(dir, name)
		if err := os.WriteFile(f, []byte(contents), 0o600); err != nil {
			t.Fatal(err)
		}
		return f
	}
	newF, oldF := write("weft.toml", "[weft]\n"), write("old.toml", "[serviceweaver]\n")
	for _, tc := range []struct{ weft, legacy, wantFile, wantCfg string }{
		{newF, "", newF, "[weft]\n"},
		{"", oldF, oldF, "[serviceweaver]\n"},
		{newF, oldF, newF, "[weft]\n"},
		{"", "", "", ""},
	} {
		t.Setenv("WEFT_CONFIG", tc.weft)
		t.Setenv("SERVICEWEAVER_CONFIG", tc.legacy)
		opts, err := localOptions()
		if err != nil {
			t.Fatal(err)
		}
		if opts.ConfigFilename != tc.wantFile || opts.Config != tc.wantCfg {
			t.Errorf("WEFT_CONFIG=%q SERVICEWEAVER_CONFIG=%q: got (%q, %q), want (%q, %q)", tc.weft, tc.legacy, opts.ConfigFilename, opts.Config, tc.wantFile, tc.wantCfg)
		}
	}
	t.Setenv("WEFT_CONFIG", filepath.Join(dir, "missing.toml"))
	if _, err := localOptions(); err == nil {
		t.Error("missing config file: no error")
	}
}
