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
	"fmt"
	"io"
)

type NybbleReader struct {
	r io.ByteReader
}

func NewNybbleReader(r io.Reader) *NybbleReader {
	return &NybbleReader{r: readerToByteReader(r)}
}

func (r *NybbleReader) ReadByte() (byte, error) {
	high, err := r.r.ReadByte()
	if err != nil {
		return 0, err
	}

	low, err := r.r.ReadByte()
	if err != nil {
		if err != io.EOF {
			return 0, err
		}
		return 0, fmt.Errorf("NybbleReader expects an even number of nybbles")
	}

	return high<<4 | low, nil
}

func (r *NybbleReader) Read(p []byte) (n int, err error) {
	return readByteReader(r, p)
}

var _ io.ByteReader = (*NybbleReader)(nil)
var _ io.Reader = (*NybbleReader)(nil)
