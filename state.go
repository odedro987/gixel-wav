package main

import (
	"encoding/json"
	"fmt"
	ic "image/color"
	"math"

	"github.com/GixelEngine/gixel-engine/gixel"
	"github.com/GixelEngine/gixel-engine/gixel/cache"
	"github.com/odedro987/gixel-wav/internal/wav"
)

type PlayState struct {
	gixel.BaseGxlState
	signal   *wav.Signal
	format   *wav.WavMetadata
	samples  gixel.GxlGroup
	idx      int
	duration float32
}

func (s *PlayState) drawSample() {
	x := (float64(s.idx) / float64(len(s.signal.Samples)-1)) * GAME_WIDTH
	y := s.signal.Samples[s.idx] * GAME_HEIGHT / 2 * -1
	t := gixel.NewSprite(float64(x), float64(y+GAME_HEIGHT/2))
	t.ApplyGraphic(s.Game().Graphics().MakeGraphic(1, 1, ic.RGBA{0, 255, 30, 255}, cache.CacheOptions{Key: "sample"}))
	s.samples.Add(t)
}

func (s *PlayState) drawWaveFormSample() {
	x := (float64(s.idx) / float64(len(s.signal.Samples)-1)) * GAME_WIDTH
	y := s.signal.Samples[s.idx] * (GAME_HEIGHT / 2) * -1
	t := gixel.NewSprite(float64(x), float64(GAME_HEIGHT/2-math.Abs(float64(y))))
	t.ApplyGraphic(s.Game().Graphics().MakeGraphic(1, int(math.Abs(float64(y))*2)+1, ic.RGBA{0, 255, 30, 255}, cache.CacheOptions{}))
	s.samples.Add(t)
}

func (s *PlayState) Init(game *gixel.GxlGame) {
	s.BaseGxlState.Init(game)

	format, signal, err := wav.Decode("assets/440.wav")
	if err != nil {
		panic(err)
	}

	s.signal = signal
	s.format = format

	s.duration = (float32(s.signal.SampleRate) / 60.0)

	marshalFormat, _ := json.Marshal(format)
	fmt.Println(string(marshalFormat))

	s.samples = gixel.NewGroup(0)
	s.Add(s.samples)
}

func (s *PlayState) Update(elapsed float64) error {
	err := s.BaseGxlState.Update(elapsed)
	if err != nil {
		return err
	}

	for i := 0; i < int(s.duration); i++ {
		if s.idx >= len(s.signal.Samples) {
			break
		}
		s.drawSample()
		s.idx += 4
	}

	return nil
}
