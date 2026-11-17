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

// EXPECTED
// M(ctx context.Context) (err error) {

// A component method may use aliases of context.Context and error.
package foo

import (
	"context"

	"github.com/brunobls/weft"
)

type Ctx = context.Context
type E = error

type foo interface {
	M(Ctx) E
}

type impl struct {
	weft.Implements[foo]
}

func (impl) M(Ctx) E { return nil }
