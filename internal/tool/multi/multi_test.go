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

package multi

import (
	"regexp"
	"testing"
)

func TestPurgeKillMatchesBothNames(t *testing.T) {
	re := regexp.MustCompile(purgeSpec.Kill)
	for _, line := range []string{
		"/usr/bin/weft multi deploy weft.toml",
		"/usr/bin/weaver multi deploy weaver.toml",
		"weft multi logs --follow x",
		"weaver multi dashboard",
		"weft multi profile --duration=30s abc",
	} {
		if !re.MatchString(line) {
			t.Errorf("Kill %q does not match %q", purgeSpec.Kill, line)
		}
	}
	for _, line := range []string{
		"weft multi purge", "weaver multi purge", "weft multi status", "weaver multi status", "weft multi version",
	} {
		if re.MatchString(line) {
			t.Errorf("Kill %q matches %q, which must not be killed by purge", purgeSpec.Kill, line)
		}
	}
	if re.MatchString("weft single dashboard") {
		t.Errorf("Kill %q matches a single process", purgeSpec.Kill)
	}
}
