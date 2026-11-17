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
// codegen.Register(codegen.Registration{
// Name:  "foo/foo",
// wEaVeReDgE:foo/foo→foo/bar⟧

// Component fields (Implements, Ref, WithRouter) may be declared through
// aliases, and the argument of Ref may itself be an alias.
package foo

import (
	"context"

	"github.com/brunobls/weft"
)

type bar interface {
	M(context.Context) error
}

type barImpl struct {
	weaver.Implements[bar]
}

func (barImpl) M(context.Context) error { return nil }

type B = bar
type RefB = weaver.Ref[B]

type foo interface {
	M(context.Context) error
}

type F = foo
type Impl = weaver.Implements[F]

type impl struct {
	Impl
	b RefB
}

func (impl) M(context.Context) error { return nil }
