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
	CircuitTracks = &Flavor{
		Name: "Circuit Tracks",
		ID:   0x64,
		SysExPrefix: []byte{
			0x00, 0x20, 0x29,
			0x01,
			0x64,
		},
		SysExSize:      350,
		NumberProjects: 64,
		NumberSamples:  64,
		NumberPatches:  128,
	}
)
