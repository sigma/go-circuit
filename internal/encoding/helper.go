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

func readByteReader(r io.ByteReader, p []byte) (n int, err error) {
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

func readerToByteReader(r io.Reader) io.ByteReader {
	if br, ok := r.(io.ByteReader); ok {
		return br
	}
	return bufio.NewReader(r)
}
