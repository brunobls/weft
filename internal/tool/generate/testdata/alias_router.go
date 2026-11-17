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
// Routed

// A router key may have fields whose types are aliases, and the router
// itself may be embedded through an alias of weaver.WithRouter.
package foo

import (
	"context"

	"github.com/brunobls/weft"
)

type K = int

type key struct {
	X K
}

type Routed interface {
	A(context.Context) error
}

type RT = router
type WR = weaver.WithRouter[RT]

type routed struct {
	weaver.Implements[Routed]
	WR
}

func (routed) A(context.Context) error { return nil }

type router struct{}

func (router) A(context.Context) key { return key{} }
