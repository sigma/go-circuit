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

package binary

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
)

type Reader []byte

func (r *Reader) advance(n int) {
	*r = (*r)[n:]
}

func (r *Reader) Uint8() uint8 {
	v := (*r)[0]
	r.advance(1)
	return v
}

type ByteOrder struct {
	r *Reader
	o binary.ByteOrder
}

func (b *ByteOrder) Uint32() uint32 {
	v := b.o.Uint32((*b.r)[:4])
	b.r.advance(4)
	return v
}

func (r *Reader) BigEndian() *ByteOrder {
	return &ByteOrder{
		r: r,
		o: binary.BigEndian,
	}
}

func (r *Reader) LittleEndian() *ByteOrder {
	return &ByteOrder{
		r: r,
		o: binary.LittleEndian,
	}
}

func (r *Reader) Section(n int) *Reader {
	r2 := (*r)[:n]
	r.advance(n)
	return &r2
}

func (r *Reader) CheckCRC(crc uint32) error {
	c := crc32.ChecksumIEEE(*r)
	if c != crc {
		return fmt.Errorf("wrong CRC: want %x, got %x", c, crc)
	}
	return nil
}
