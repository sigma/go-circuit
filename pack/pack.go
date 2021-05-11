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
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/sigma/go-circuit/model"
	"gitlab.com/gomidi/midi/reader"
)

type Pack struct {
	Name     string
	Color    string
	Projects []*Project
	Samples  []*Sample
	Patches  []*Patch
}

func (p *Pack) Write(w io.Writer, f *model.Flavor) error {
	if n := len(p.Projects); n > f.NumberProjects {
		return fmt.Errorf("too many projects: %d", n)
	}
	if n := len(p.Samples); n > f.NumberSamples {
		return fmt.Errorf("too many samples: %d", n)
	}
	if n := len(p.Patches); n > f.NumberPatches {
		return fmt.Errorf("too many patches: %d", n)
	}

	if f == model.CircuitTracks {
		return p.writeCircuitTracks(w)
	}

	return fmt.Errorf("unsupported flavor: %s", f.Name)
}

func (p *Pack) Read(r io.Reader) error {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(r); err != nil {
		return err
	}

	c, err := buf.ReadByte()
	if err != nil {
		return err
	}
	buf.UnreadByte()

	if c == 0x50 {
		return fmt.Errorf("unsupported input format: %s", model.CircuitTracks.Name)
	}

	return p.readCircuit(buf)
}

type packIndex struct {
	Name     string        `json:"name"`
	Color    string        `json:"color"`
	Product  string        `json:"product"`
	Version  string        `json:"version"`
	Projects []*packObject `json:"projects"`
	Samples  []*packObject `json:"samples"`
	Patches  []*packObject `json:"patches"`
}

type packObject struct {
	Name string `json:"name"`
	Path string `json:"url,omitempty"`
}

func (p *Pack) writeCircuitTracks(w io.Writer) error {
	f := model.CircuitTracks
	zw := zip.NewWriter(w)
	defer zw.Close()

	idx := &packIndex{
		Name:    p.Name,
		Color:   p.Color,
		Product: "circuit-tracks",
		Version: "2.0",
	}

	zw.CreateHeader(&zip.FileHeader{
		Name: "projects/",
	})

	for i := 0; i < f.NumberProjects; i++ {
		fname := fmt.Sprintf("projects/project_%d.ncs", i)

		project := &Project{}

		idx.Projects = append(idx.Projects, &packObject{
			Name: "",
			Path: fname,
		})

		w, err := zw.Create(fname)
		if err != nil {
			return err
		}
		w.Write(project.Format(&ProjectConfig{Flavor: f}))
	}

	for i := 0; i < f.NumberSamples; i++ {
		fname := fmt.Sprintf("samples/sample_%d.wav", i)

		idx.Samples = append(idx.Samples, &packObject{
			Name: "",
			Path: fname,
		})
	}

	zw.CreateHeader(&zip.FileHeader{
		Name: "patches/",
	})

	for i := 0; i < f.NumberPatches; i++ {
		fname := fmt.Sprintf("patches/patch_%d.syx", i)

		patch := &Patch{}
		if i < len(p.Patches) {
			patch = p.Patches[i]
		}

		idx.Patches = append(idx.Patches, &packObject{
			Name: patch.Name(),
			Path: fname,
		})

		w, err := zw.Create(fname)
		if err != nil {
			return err
		}
		w.Write([]byte{0xf0})
		w.Write(patch.Format(&PatchConfig{Flavor: f, Index: byte(i)}))
		w.Write([]byte{0xf7})
	}

	iw, err := zw.Create("index.json")
	if err != nil {
		return err
	}
	body, err := json.Marshal(idx)
	if err != nil {
		return err
	}
	iw.Write(body)

	return nil
}

func (p *Pack) readCircuit(r io.Reader) error {
	prefix := append(
		manufacturerID,
		0x01,
		0x60,
	)

	syxReader := func(_ *reader.Position, data []byte) {
		if bytes.Equal(prefix, data[:len(prefix)]) {
			p.Patches = append(p.Patches, &Patch{data: data[8:]})
		}
	}

	midiReader := reader.New(reader.SysEx(syxReader), reader.NoLogger())

	if err := reader.ReadAllFrom(midiReader, r); err != nil && err != io.EOF {
		return err
	}

	return nil
}
