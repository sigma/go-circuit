// Copyright 2021 Google LLC
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

package encoding

import (
	"bufio"
	"io"
)

type Low7Reader struct {
	r        io.ByteReader
	highBits byte
	nBits    byte
}

func NewLow7Reader(r io.Reader) *Low7Reader {
	if br, ok := r.(io.ByteReader); ok {
		return &Low7Reader{r: br}
	}
	return &Low7Reader{r: bufio.NewReader(r)}
}

func (r *Low7Reader) ReadByte() (byte, error) {
	if r.nBits == 0 {
		bits, err := r.r.ReadByte()
		if err != nil {
			return 0, err
		}
		r.highBits = bits
		r.nBits = 7
	}

	next, err := r.r.ReadByte()
	if err != nil {
		return 0, err
	}

	high := (r.highBits & 0x01) * 0x80
	r.highBits = r.highBits >> 1
	r.nBits--
	return high | next, nil
}

func (r *Low7Reader) Read(p []byte) (n int, err error) {
	for i := range p {
		next, e := r.ReadByte()
		if e != nil {
			err = e
			return
		}
		p[i] = next
		n++
	}
	return
}

var _ io.ByteReader = (*Low7Reader)(nil)
var _ io.Reader = (*Low7Reader)(nil)
