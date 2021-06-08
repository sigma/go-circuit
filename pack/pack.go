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
	"io/ioutil"

	"github.com/go-audio/wav"
	"github.com/orcaman/writerseeker"
	"github.com/sigma/go-circuit/internal/binary"
	"github.com/sigma/go-circuit/internal/encoding"
	"github.com/sigma/go-circuit/model"
	"gitlab.com/gomidi/midi/reader"
)

type Pack struct {
	Name     string
	Color    string
	Projects []*Project
	Samples  []*Sample
	Patches  []*Patch

	rawSamples []byte
	inSamples  bool
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

	if err := buf.UnreadByte(); err != nil {
		return err
	}

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

	_, _ = zw.CreateHeader(&zip.FileHeader{
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
		if _, err := w.Write(project.Format(&ProjectConfig{Flavor: f})); err != nil {
			return err
		}
	}

	_, _ = zw.CreateHeader(&zip.FileHeader{
		Name: "samples/",
	})

	for i := 0; i < f.NumberSamples; i++ {
		fname := fmt.Sprintf("samples/sample_%d.wav", i)

		sample := &Sample{}

		if i < len(p.Samples) {
			sample = p.Samples[i]
		}

		idx.Samples = append(idx.Samples, &packObject{
			Name: sample.Name,
			Path: fname,
		})

		// TODO: normalize wav format if needed
		if sample.Data != nil {
			w, err := zw.Create(fname)
			if err != nil {
				return err
			}
			data, err := ioutil.ReadAll(sample.Data)
			if err != nil {
				return err
			}
			if _, err := w.Write(data); err != nil {
				return err
			}
		}
	}

	_, _ = zw.CreateHeader(&zip.FileHeader{
		Name: "patches/",
	})

	for i := 0; i < f.NumberPatches; i++ {
		fname := fmt.Sprintf("patches/patch_%d.syx", i)

		var patch *Patch = nil
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
		if _, err := w.Write([]byte{0xf0}); err != nil {
			return err
		}
		if _, err := w.Write(patch.Format(&PatchConfig{Flavor: f, Index: byte(i)})); err != nil {
			return err
		}
		if _, err := w.Write([]byte{0xf7}); err != nil {
			return err
		}
	}

	iw, err := zw.Create("index.json")
	if err != nil {
		return err
	}
	body, err := json.Marshal(idx)
	if err != nil {
		return err
	}
	if _, err := iw.Write(body); err != nil {
		return err
	}

	return nil
}

func (p *Pack) readCircuit(r io.Reader) error {
	samplePrefix := model.Circuit.SysExSamplePrefix()
	patchPrefix := model.Circuit.SysExPatchPrefix()
	var readErr error

	syxReader := func(_ *reader.Position, data []byte) {
		if bytes.Equal(samplePrefix, data[:len(samplePrefix)]) {
			readErr = p.readSysexData(data[len(samplePrefix):])
		} else if bytes.Equal(patchPrefix, data[:len(patchPrefix)]) {
			p.Patches = append(p.Patches, NewPatch(data))
		}
	}

	midiReader := reader.New(reader.SysEx(syxReader), reader.NoLogger())

	if err := reader.ReadAllFrom(midiReader, r); err != nil && err != io.EOF {
		return err
	}

	return readErr
}

func (p *Pack) parseSamples(crc uint32) error {
	r := binary.Reader(p.rawSamples)
	if err := r.CheckCRC(crc); err != nil {
		return err
	}
	n := int(r.Uint8())
	for i := 0; i < n; i++ {
		channels := r.Uint8()
		bits := r.Uint8()
		rate := r.LittleEndian().Uint32()

		writer := &writerseeker.WriterSeeker{}

		e := wav.NewEncoder(writer,
			int(rate),
			int(bits),
			int(channels),
			1)
		defer e.Close()

		length := r.LittleEndian().Uint32()
		size := uint32(bits / 8)
		nframes := length / size
		s := r.Section(int(length))

		for f := 0; f < int(nframes); f++ {
			frame := make([]byte, size)
			for i := 0; i < int(size); i++ {
				frame[int(size)-i-1] = s.Uint8()
			}
			if err := e.WriteFrame(frame); err != nil {
				return err
			}
		}

		sample := &Sample{
			Name: "",
			Data: writer.Reader(),
		}
		p.Samples = append(p.Samples, sample)
	}
	return nil
}

func (p *Pack) readSysexData(data []byte) error {
	switch cmd := data[0]; cmd {
	case 0x77:
		// TODO: allocate the full unpacked slice here
		// We have potentially 2 sections sharing the same sysex command:
		// - the sessions one
		// - the samples one
		// The samples section is identified by the next 8 nybbles: 0x0023b000
		r := encoding.NewNybbleReader(bytes.NewBuffer(data[1:9]))
		chunk, err := io.ReadAll(r)
		if err != nil {
			return err
		}
		if bytes.Equal(chunk, []byte{0x00, 0x23, 0xb0, 0x00}) {
			p.inSamples = true
		}
	case 0x79:
		if p.inSamples {
			r := encoding.NewLow7Reader(bytes.NewBuffer(data[1:]))
			chunk, err := io.ReadAll(r)
			if err != nil {
				return err
			}
			p.rawSamples = append(p.rawSamples, chunk...)
		}
	case 0x7a:
		if p.inSamples {
			r := encoding.NewNybbleReader((bytes.NewBuffer(data[1:9])))
			body, err := io.ReadAll(r)
			if err != nil {
				return err
			}

			crc := binary.Reader(body)
			if err := p.parseSamples(crc.BigEndian().Uint32()); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("invalid sample sysex cmd: %v", cmd)
	}
	return nil
}
