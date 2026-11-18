// Copyright 2023 Google LLC
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
// NoRetry: []int{1}

// A non-retriable method may be named through an alias of the component
// interface.
package foo

import (
	"context"

	"github.com/brunobls/weft"
)

type foo interface {
	A(context.Context) error
	B(context.Context) error
}

type impl struct{ weft.Implements[foo] }

func (l *impl) A(context.Context) error { return nil }
func (l *impl) B(context.Context) error { return nil }

type F = foo

var _ weft.NotRetriable = F.B
