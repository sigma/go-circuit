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
	"strings"

	"github.com/sigma/go-circuit/model"
)

var (
	manufacturerID = []byte{0x00, 0x20, 0x29}
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

type Patch struct {
	data []byte
}

func patchKind(sysex []byte) *model.Flavor {
	for _, m := range []*model.Flavor{model.Circuit, model.CircuitTracks} {
		if bytes.Equal(m.SysExPrefix, sysex[:len(m.SysExPrefix)]) {
			return m
		}
	}
	return nil
}

func NewPatch(sysex []byte) *Patch {
	if k := patchKind(sysex); k != nil {
		if len(sysex) == k.SysExSize {
			return &Patch{
				data: sysex[len(sysex)-340:],
			}
		}
		return &Patch{}
	}
	return nil
}

func (p *Patch) Name() string {
	if p.data == nil {
		return "Initial Patch"
	}
	return strings.TrimSpace(string(p.data[0:16]))
}

func (p *Patch) Category() Category {
	if p.data == nil {
		return CategoryNone
	}
	return Category(p.data[16])
}

func (p *Patch) Genre() Genre {
	if p.data == nil {
		return GenreNone
	}
	return Genre(p.data[17])
}

type PatchConfig struct {
	Flavor *model.Flavor
	Index  byte
}

func (p *Patch) Format(cfg *PatchConfig) []byte {
	prelude := append(
		manufacturerID,
		0x01,
		cfg.Flavor.ID,
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

	res := append(
		prelude,
		p.data...,
	)
	return res
}
