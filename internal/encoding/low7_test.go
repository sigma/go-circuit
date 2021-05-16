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
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLow7Decode(t *testing.T) {
	data := []struct {
		name     string
		enc, dec []byte
	}{
		{
			name: "single block",
			enc:  []byte{0x18, 0x40, 0x01, 0x10, 0x00, 0x3b, 0x00, 0x00},
			dec:  []byte{ /**/ 0x40, 0x01, 0x10, 0x80, 0xbb, 0x00, 0x00},
		},
	}

	for _, tt := range data {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dec, err := Low7Decode(tt.enc)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tt.dec, dec); diff != "" {
				t.Errorf("Low7Decode() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
