// Copyright © 2018 The Things Network Foundation, The Things Industries B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package random implements goroutine-safe utilities on top of a secure random source.
package random

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"math/big"
)

// Interface for random.
type Interface interface {
	Intn(n int) int
	String(n int) string
	Bytes(n int) []byte
	FillBytes(p []byte)
}

// TTNRandom is used as a wrapper around crypto/rand.
type TTNRandom struct {
	Source io.Reader
}

// New returns a new Random, in most cases you can also just use the global funcs.
func New() Interface {
	return &TTNRandom{
		Source: rand.Reader,
	}
}

var global = New()

// Intn returns a random number in the range [0,n). This func uses the global TTNRandom.
func Intn(n int) int { return global.Intn(n) }

// Intn returns a random number in the range [0,n).
func (r *TTNRandom) Intn(n int) int {
	i, err := rand.Int(r.Source, big.NewInt(int64(n)))
	if err != nil {
		panic(err) // r.Source is (very) broken.
	}
	return int(i.Int64())
}

// Bytes generates a random byte slice of length n. This func uses the global TTNRandom.
func Bytes(n int) []byte { return global.Bytes(n) }

// Bytes generates a random byte slice of length n.
func (r *TTNRandom) Bytes(n int) []byte {
	p := make([]byte, n)
	r.FillBytes(p)
	return p
}

// FillBytes fills the byte slice with random bytes. This func uses the global TTNRandom.
func FillBytes(p []byte) { global.FillBytes(p) }

// FillBytes fills the byte slice with random bytes.
func (r *TTNRandom) FillBytes(p []byte) {
	_, err := r.Source.Read(p)
	if err != nil {
		panic(fmt.Errorf("random.Bytes: %s", err))
	}
}

// String returns a random string of length n, it uses the characters of base64.URLEncoding.
// This func uses the global TTNRandom.
func String(n int) string { return global.String(n) }

// String returns a random string of length n, it uses the characters of base64.URLEncoding.
func (r *TTNRandom) String(n int) string {
	b := r.Bytes(n * 6 / 8)
	return base64.RawURLEncoding.EncodeToString(b)
}
