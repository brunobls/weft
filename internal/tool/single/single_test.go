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

package single

import (
	"regexp"
	"slices"
	"testing"
)

func TestPurgeKillMatchesBothNames(t *testing.T) {
	re := regexp.MustCompile(purgeSpec.Kill)
	for _, line := range []string{
		"/usr/bin/weft single dashboard",
		"/usr/bin/weaver single dashboard",
		"weft single profile --duration=30s abc",
		"weaver single profile --duration=30s abc",
	} {
		if !re.MatchString(line) {
			t.Errorf("Kill %q does not match %q", purgeSpec.Kill, line)
		}
	}
	for _, line := range []string{
		"weft single purge", "weaver single purge", "weft single status", "weaver single status", "weft single version",
	} {
		if re.MatchString(line) {
			t.Errorf("Kill %q matches %q, which must not be killed by purge", purgeSpec.Kill, line)
		}
	}
	if re.MatchString("weft multi dashboard") {
		t.Errorf("Kill %q matches a multi process", purgeSpec.Kill)
	}
}

func TestAppEnvWritesBothVariables(t *testing.T) {
	got := appEnv([]string{"A=b"}, "/x/weft.toml")
	for _, want := range []string{"A=b", "WEFT_CONFIG=/x/weft.toml", "SERVICEWEAVER_CONFIG=/x/weft.toml"} {
		if !slices.Contains(got, want) {
			t.Errorf("env %v lacks %q", got, want)
		}
	}
}
