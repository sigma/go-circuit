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

package mime

import (
	"bytes"

	"github.com/gabriel-vasile/mimetype"
)

const (
	sysexMime             = "application/vnd.novation.sysex"
	CircuitPackMime       = "application/vnd.novation.sysex.circuit.pack"
	CircuitTracksPackMime = "application/vnd.novation.circuit.tracks.pack+zip"
)

func init() {
	mimetype.Lookup("application/zip").Extend(circuitTracksPack, CircuitTracksPackMime, ".circuittrackspack")
	mimetype.Extend(novationSysex, sysexMime, ".syx")

	syx := mimetype.Lookup(sysexMime)
	syx.Extend(circuitPackSysex, CircuitPackMime, ".circuitpack")
}

func circuitTracksPack(raw []byte, _ uint32) bool {
	typ := "projects/PK"
	return (len(raw) >= 30+len(typ) &&
		string(raw[30:30+len(typ)]) == typ)
}

func novationSysex(raw []byte, limit uint32) bool {
	return len(raw) > 4 && bytes.Equal(raw[:4], []byte{0xf0, 0x00, 0x20, 0x29})
}

func circuitPackSysex(raw []byte, _ uint32) bool {
	return len(raw) > 6 && bytes.Equal(raw[4:6], []byte{0x00, 0x77})
}
