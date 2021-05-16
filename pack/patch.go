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

package pack

import (
	"bytes"
	"encoding/binary"
	"strings"

	"github.com/sigma/go-circuit/model"
)

type Genre byte

const (
	GenreNone Genre = iota
	GenreClassic
	GenreBreaks
	GenreHouse
	GenreIndustrial
	GenreJazz
	GenreHipHop
	GenrePopRock
	GenreTechno
	GenreDubStep
)

type Category byte

const (
	CategoryNone Category = iota
	CategoryArp
	CategoryBass
	CategoryBell
	CategoryClassic
	CategoryDrum
	CategoryKeyboard
	CategoryLead
	CategoryMotion
	CategoryPad
	CategoryPoly
	CategorySFX
	CategoryString
	CategoryUser
	CategoryVocal
)

type Voice struct {
	PolyphonyMode  byte
	PortamentoRate byte
	PreGlide       byte
	KeyboardOctave byte
}

type Mixer struct {
	Osc1Level       byte
	Osc2Level       byte
	RingModeLevel12 byte
	NoiseLevel      byte
	PreFXLevel      byte
	PostFXLevel     byte
}

type Filter struct {
	Routing    byte
	Drive      byte
	DriveType  byte
	Type       byte
	Frequency  byte
	Track      byte
	Resonance  byte
	QNormalize byte
	Env2ToFreq byte
}

type Oscillator struct {
	Wave             byte
	WaveInterpolate  byte
	PulseWidthIndex  byte
	VirtualSyncDepth byte
	Density          byte
	DensityDetune    byte
	Semitones        byte
	Cents            byte
	PitchBend        byte
}

type ADSR struct {
	Attack  byte
	Decay   byte
	Sustain byte
	Release byte
}

type VelocityEnvelope struct {
	Velocity byte
	ADSR
}

type DelayEnvelope struct {
	Delay byte
	ADSR
}

type LFO struct {
	WaveForm    byte
	PhaseOffset byte
	SlewRate    byte
	Delay       byte
	DelaySync   byte
	Rate        byte
	RateSync    byte
	Bits        byte
}

type Band struct {
	Frequency byte
	Level     byte
}

type Equalizer struct {
	Bass, Mid, Trebble Band
}

type Distortion struct {
	Type         byte
	Compensation byte
}

type Chorus struct {
	Type     byte
	Rate     byte
	RateSync byte
	Feedback byte
	ModDepth byte
	Delay    byte
}

type Mod struct {
	Source1, Source2 byte
	Depth            byte
	Destination      byte
}

type KnobTarget struct {
	Destination byte
	Start, End  byte
	Depth       byte
}

type Knob struct {
	Position   byte
	A, B, C, D KnobTarget
}

type Patch struct {
	PatchName                [16]byte
	Category                 Category
	Genre                    Genre
	Reserved                 [14]byte
	Voice                    Voice
	Osc1, Osc2               Oscillator
	Mixer                    Mixer
	Filter                   Filter
	Envelope1, Envelope2     VelocityEnvelope
	Envelope3                DelayEnvelope
	LFO1, LFO2               LFO
	DistortionLevel          byte
	FXReserved1              byte
	ChorusLevel              byte
	FXReserved2, FXReserved3 byte
	Equalizer                Equalizer
	FXReserved               [5]byte
	Distortion               Distortion
	Chorus                   Chorus
	ModMatrix                [20]Mod
	Macros                   [8]Knob
}

func patchKind(sysex []byte) *model.Flavor {
	for _, m := range []*model.Flavor{model.Circuit, model.CircuitTracks} {
		if prefix := m.SysExPatchPrefix(); bytes.Equal(prefix, sysex[:len(prefix)]) {
			return m
		}
	}
	return nil
}

func NewPatch(sysex []byte) *Patch {
	var p Patch
	if k := patchKind(sysex); k != nil {
		if len(sysex) == k.SysExSize {
			binary.Read(bytes.NewReader(sysex[len(sysex)-340:]), binary.LittleEndian, &p)
		}
		return &p
	}
	return nil
}

func (p *Patch) Name() string {
	if p == nil {
		return ""
	}
	name := strings.TrimSpace(string(p.PatchName[:]))
	if name == "" {
		name = "Initial Patch"
	}
	return name
}

type PatchConfig struct {
	Flavor *model.Flavor
	Index  byte
}

func (p *Patch) Format(cfg *PatchConfig) []byte {
	prelude := append(
		cfg.Flavor.SysExPatchPrefix(),
		0x01,
	)

	if cfg.Flavor == model.CircuitTracks {
		prelude = append(prelude, 0x7f, 0x7f)
	}

	prelude = append(
		prelude,
		cfg.Index,
		0x00,
	)

	data := new(bytes.Buffer)
	if p != nil {
		binary.Write(data, binary.LittleEndian, p)
	}

	res := append(
		prelude,
		data.Bytes()...,
	)
	return res
}
