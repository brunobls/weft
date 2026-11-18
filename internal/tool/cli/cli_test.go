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

package cli

import (
	"errors"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestUsageUsesInvokedName(t *testing.T) {
	for _, name := range []string{"weft", "weaver"} {
		u := Usage(name)
		for _, want := range []string{
			name + " generate", name + " version", name + " single",
			name + " multi", name + " ssh", name + " gke-local", name + " kube",
			`"` + name + `-<deployer>"`,
		} {
			if !strings.Contains(u, want) {
				t.Errorf("Usage(%q) lacks %q:\n%s", name, want, u)
			}
		}
		other := map[string]string{"weft": "weaver ", "weaver": "weft "}[name]
		if strings.Contains(u, other) {
			t.Errorf("Usage(%q) mentions the other name %q:\n%s", name, other, u)
		}
	}
}

func TestUsageAliasLine(t *testing.T) {
	if !strings.Contains(Usage("weaver"), "weaver is an alias of weft") {
		t.Errorf("weaver usage lacks the alias line")
	}
	u := Usage("weaver")
	i, j := strings.Index(u, `"weft-<deployer>"`), strings.Index(u, `"weaver-<deployer>"`)
	if i < 0 || j < 0 || i > j {
		t.Errorf("weaver usage does not document weft-<deployer> before weaver-<deployer>:\n%s", u)
	}
	if strings.Contains(Usage("weft"), "weaver") {
		t.Errorf("weft usage mentions weaver:\n%s", Usage("weft"))
	}
}

func TestPluginCandidates(t *testing.T) {
	got := pluginCandidates("gke")
	want := []string{"weft-gke", "weaver-gke"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("pluginCandidates(gke) = %v, want %v", got, want)
	}
}

func TestFindPlugin(t *testing.T) {
	has := func(names ...string) func(string) (string, error) {
		return func(n string) (string, error) {
			for _, h := range names {
				if h == n {
					return "/bin/" + n, nil
				}
			}
			return "", errors.New("not found")
		}
	}
	for _, tc := range []struct {
		name    string
		have    []string
		want    string
		wantErr bool
	}{
		{"weft only", []string{"weft-gke"}, "weft-gke", false},
		{"weaver only (fallback)", []string{"weaver-gke"}, "weaver-gke", false},
		{"both prefers weft", []string{"weaver-gke", "weft-gke"}, "weft-gke", false},
		{"none", nil, "", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := findPlugin("gke", has(tc.have...))
			if (err != nil) != tc.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tc.wantErr)
			}
			if got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
			if tc.wantErr {
				for _, n := range []string{"weft-gke", "weaver-gke"} {
					if !strings.Contains(err.Error(), n) {
						t.Errorf("error %q does not name %q", err, n)
					}
				}
			}
		})
	}
}

// TestBinaries builds both entry points and checks that the invoked name shows
// in the help and that generate is reachable through both.
func TestBinaries(t *testing.T) {
	if testing.Short() {
		t.Skip("builds binaries")
	}
	dir := t.TempDir()
	for _, name := range []string{"weft", "weaver"} {
		bin := filepath.Join(dir, name)
		build := exec.Command("go", "build", "-buildvcs=false", "-o", bin, "github.com/brunobls/weft/cmd/"+name)
		if out, err := build.CombinedOutput(); err != nil {
			t.Fatalf("build %s: %v\n%s", name, err, out)
		}
		out, err := exec.Command(bin, "help").CombinedOutput()
		if err != nil {
			t.Fatalf("%s help: %v\n%s", name, err, out)
		}
		if !strings.Contains(string(out), name+" generate") {
			t.Errorf("%s help does not name %q:\n%s", name, name, out)
		}
		out, err = exec.Command(bin, "help", "generate").CombinedOutput()
		if err != nil || !strings.Contains(string(out), name+" generate") {
			t.Errorf("%s help generate: err=%v\n%s", name, err, out)
		}
		// The group title and its lines name the same command ("weft multi").
		out, _ = exec.Command(bin, "multi", "--help").CombinedOutput()
		if !strings.Contains(string(out), "weft multi") || strings.Contains(string(out), "weaver multi") {
			t.Errorf("%s multi --help mixes names:\n%s", name, out)
		}
		out, _ = exec.Command(bin, "nosuch").CombinedOutput()
		if !strings.Contains(string(out), `"`+name+" nosuch\"") {
			t.Errorf("%s nosuch error does not use invoked name:\n%s", name, out)
		}
	}
}
