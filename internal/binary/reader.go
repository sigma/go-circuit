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

func (r *Reader) read(n int) []byte {
	b := (*r)[:n]
	*r = (*r)[n:]
	return b
}

func (r *Reader) Uint8() uint8 {
	return r.read(1)[0]
}

type ByteOrder struct {
	reader *Reader
	order  binary.ByteOrder
}

func (b *ByteOrder) Uint32() uint32 {
	return b.order.Uint32(b.reader.read(4))
}

func (r *Reader) BigEndian() *ByteOrder {
	return &ByteOrder{
		reader: r,
		order:  binary.BigEndian,
	}
}

func (r *Reader) LittleEndian() *ByteOrder {
	return &ByteOrder{
		reader: r,
		order:  binary.LittleEndian,
	}
}

func (r *Reader) Section(n int) *Reader {
	r2 := Reader(r.read(n))
	return &r2
}

func (r *Reader) CheckCRC(crc uint32) error {
	c := crc32.ChecksumIEEE(*r)
	if c != crc {
		return fmt.Errorf("wrong CRC: want %x, got %x", c, crc)
	}
	return nil
}
