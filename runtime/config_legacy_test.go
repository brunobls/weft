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

package runtime

import (
	"path/filepath"
	"strings"
	"testing"
)

func noValidation(string, string) error { return nil }

// TestLegacyConfigSections checks that a config section named after the old
// module path is an error that names the new section, instead of being
// silently ignored (which would, for example, start an app without mTLS when
// the old multi section had mtls = true).
func TestLegacyConfigSections(t *testing.T) {
	old, new := legacyModulePath, modulePath
	for _, c := range []struct {
		name, cfg string
		want      string // the new name the error must mention
	}{
		{"multi mtls", "[\"" + old + "/multi\"]\nmtls = true\n", new + "/multi"},
		{"single", "[\"" + old + "/single\"]\n", new + "/single"},
		{"ssh", "[\"" + old + "/ssh\"]\n", new + "/ssh"},
		{"app", "[\"" + old + "\"]\nbinary = \"./x\"\n", new},
		{"component", "[\"" + old + "/examples/chat/SQLStore\"]\nx = 1\n", new + "/examples/chat/SQLStore"},
		{"framework component", "[\"" + old + "/Main\"]\n", new + "/Main"},
		{"colocate", "[serviceweaver]\ncolocate = [[\"" + old + "/Main\", \"a/B\"]]\n", new + "/Main"},
		{"wins over other config errors", "[\"" + old + "/multi\"]\nmtls = true\n[serviceweaver]\nenv = [\"NOT_AN_ASSIGNMENT\"]\n", new + "/multi"},
		{"colocate with the new app section", "[\"" + new + "\"]\ncolocate = [[\"" + old + "/Main\"]]\n", new + "/Main"},
	} {
		t.Run(c.name, func(t *testing.T) {
			_, err := ParseConfig("weaver.toml", c.cfg, noValidation)
			if err == nil {
				t.Fatalf("ParseConfig succeeded; want an error that names %q", c.want)
			}
			if got := err.Error(); !strings.Contains(got, "no longer supported") || !strings.Contains(got, c.want) {
				t.Fatalf("error %q does not say the section is no longer supported and name %q", got, c.want)
			}
		})
	}
}

// TestNewConfigSectionsAccepted is the control of TestLegacyConfigSections:
// the same sections under the new names (and the short key) are accepted.
func TestNewConfigSectionsAccepted(t *testing.T) {
	new := modulePath
	for _, c := range []struct{ name, cfg string }{
		{"multi mtls", "[\"" + new + "/multi\"]\nmtls = true\n"},
		{"single", "[\"" + new + "/single\"]\n"},
		{"app long key", "[\"" + new + "\"]\nbinary = \"./x\"\n"},
		{"app short key", "[serviceweaver]\nbinary = \"./x\"\n"},
		{"component", "[\"" + new + "/examples/chat/SQLStore\"]\nx = 1\n"},
		{"colocate", "[serviceweaver]\ncolocate = [[\"" + new + "/Main\", \"a/B\"]]\n"},
		{"unrelated module with the same prefix", "[\"" + legacyModulePath + "-kube\"]\nx = 1\n"},
	} {
		t.Run(c.name, func(t *testing.T) {
			if _, err := ParseConfig("weaver.toml", c.cfg, noValidation); err != nil {
				t.Fatalf("ParseConfig: unexpected error %v", err)
			}
		})
	}
}

// TestWeftShortSection checks the [weft] short key of the application section,
// next to the kept [serviceweaver] one.
func TestWeftShortSection(t *testing.T) {
	for _, c := range []struct {
		name, cfg string
		wantBin   string // base name of the resolved binary
		wantErr   string
	}{
		{"weft", "[weft]\nbinary = \"./x\"\n", "x", ""},
		{"serviceweaver kept", "[serviceweaver]\nbinary = \"./y\"\n", "y", ""},
		{"both", "[weft]\nbinary = \"./x\"\n[serviceweaver]\nbinary = \"./y\"\n", "", "conflicting sections"},
		{"weft and long key", "[weft]\nbinary = \"./x\"\n[\"" + modulePath + "\"]\nbinary = \"./y\"\n", "", "conflicting sections"},
		{"weft unknown key", "[weft]\nnope = 1\n", "", "unknown keys"},
	} {
		t.Run(c.name, func(t *testing.T) {
			cfg, err := ParseConfig("weft.toml", c.cfg, noValidation)
			if c.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), c.wantErr) {
					t.Fatalf("err = %v, want containing %q", err, c.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if filepath.Base(cfg.Binary) != c.wantBin {
				t.Errorf("Binary = %q, want %q", cfg.Binary, c.wantBin)
			}
		})
	}
}
