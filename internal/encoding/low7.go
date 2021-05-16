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
	"math"
)

func Low7Decode(data []byte) ([]byte, error) {
	overhead := int(math.Ceil(float64(len(data)) / 8))
	size := len(data) - overhead
	res := make([]byte, size)

	var (
		loop     byte = 7
		highBits byte = 0
		idx      int  = 0
	)
	for _, b := range data {
		if loop < 7 {
			if (highBits & (1 << loop)) != 0 {
				b += 0x80
			}
			res[idx] = b
			loop++
			idx++
		} else {
			highBits = b
			loop = 0
		}
	}

	return res, nil
}
