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
// EncodeBinaryMarshaler
// DecodeBinaryUnmarshaler
// RegisterSerializable[*sm]()

// Types that implement encoding.BinaryMarshaler, the weaver encoding methods
// or proto.Message through aliases in their method signatures are recognized,
// and so are aliases of pointers to such types.
package foo

import (
	"context"

	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/runtime/codegen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Byte = byte
type Bytes = []Byte
type Str = string
type Enc = codegen.Encoder
type Dec = codegen.Decoder
type EP = *codegen.Encoder
type DP = *codegen.Decoder
type Msg = protoreflect.Message

// bm implements encoding.BinaryMarshaler and BinaryUnmarshaler.
type bm struct{ notSerializable chan int }

func (*bm) MarshalBinary() (Bytes, error) { return nil, nil }
func (*bm) UnmarshalBinary(Bytes) error   { return nil }

type BP = *bm

// we implements the weaver encoding methods by hand.
type we struct{ notSerializable chan int }

func (we) WeaverMarshal(*Enc)    {}
func (*we) WeaverUnmarshal(*Dec) {}

// we2 does the same, with aliases of the pointer types.
type we2 struct{ notSerializable chan int }

func (we2) WeaverMarshal(EP)    {}
func (*we2) WeaverUnmarshal(DP) {}

// pm looks like a proto message.
type pm struct{ notSerializable chan int }

func (*pm) ProtoReflect() Msg { return nil }

type PMP = *pm

// sm implements error with an alias result.
type sm struct {
	weaver.AutoMarshal
	A int
}

func (sm) Error() Str { return "" }

type foo interface {
	M(context.Context, BP, *bm) (BP, error)
	N(context.Context, we, we2, *pm, sm) error
	O(context.Context, PMP) (PMP, error)
}

type impl struct {
	weaver.Implements[foo]
}

func (impl) M(context.Context, BP, *bm) (BP, error)    { return nil, nil }
func (impl) O(context.Context, PMP) (PMP, error)       { return nil, nil }
func (impl) N(context.Context, we, we2, *pm, sm) error { return nil }
