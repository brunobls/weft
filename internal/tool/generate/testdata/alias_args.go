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
// size += serviceweaver_size_outer_
// func serviceweaver_size_fixed_
// M(ctx context.Context, a0 P, a1 PP, a2 S, a3 M, a4 *P, a5 Bytes, a6 ...P) (r0 P, r1 S, r2 Bytes, err error) {

// Arguments and results of component methods may be aliases of structs,
// pointers, slices, maps and bytes, and of fixed-size structs.
package foo

import (
	"context"

	"github.com/ServiceWeaver/weaver"
)

type pair struct {
	weaver.AutoMarshal
	A int
	B string
}

type fixed struct {
	weaver.AutoMarshal
	A int32
	B bool
}

type P = pair
type PP = *pair
type S = []P
type M = map[string]P
type Byte = byte
type Bytes = []Byte
type Fixed = fixed
type FixedArray = [2]Fixed

type outer struct {
	weaver.AutoMarshal
	F Fixed
	A FixedArray
	S []Fixed
	M map[Fixed]Fixed
}

type foo interface {
	M(context.Context, P, PP, S, M, *P, Bytes, ...P) (P, S, Bytes, error)
	N(context.Context, outer, FixedArray) (outer, error)
}

type impl struct {
	weaver.Implements[foo]
}

func (impl) M(context.Context, P, PP, S, M, *P, Bytes, ...P) (P, S, Bytes, error) {
	return P{}, nil, nil, nil
}

func (impl) N(context.Context, outer, FixedArray) (outer, error) { return outer{}, nil }
