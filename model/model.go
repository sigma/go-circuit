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

package model

var (
	VendorID      = []byte{0x00, 0x20, 0x29}
	SampleID byte = 0x00
	SynthID  byte = 0x01
)

type Flavor struct {
	Name           string
	ID             byte
	SysExSize      int
	NumberProjects int
	NumberSamples  int
	NumberPatches  int
}

func (f *Flavor) SysExSamplePrefix() []byte {
	return append(VendorID, SampleID)
}

func (f *Flavor) SysExPatchPrefix() []byte {
	return append(VendorID, SynthID, f.ID)
}
