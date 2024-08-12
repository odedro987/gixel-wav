package main

import (
	"encoding/json"
	"fmt"
	ic "image/color"

	"github.com/GixelEngine/gixel-engine/gixel"
	"github.com/GixelEngine/gixel-engine/gixel/cache"
	"github.com/odedro987/gixel-wav/internal/wav"
)

type PlayState struct {
	gixel.BaseGxlState
}

func (s *PlayState) Init(game *gixel.GxlGame) {
	s.BaseGxlState.Init(game)

	format, signal, err := wav.Decode("assets/440.wav")
	if err != nil {
		panic(err)
	}

	marshalFormat, _ := json.Marshal(format)
	fmt.Println(string(marshalFormat))

	for i := 0; i < len(signal.Samples)/10; i++ {
		fmt.Println(i, float64(i)/float64(format.SampleRate-1)*GAME_WIDTH, signal.Samples[i], float64(signal.Samples[i]*GAME_HEIGHT/2+GAME_HEIGHT/2))
		t := gixel.NewSprite(float64(i)/float64(format.SampleRate-1)*GAME_WIDTH, float64(signal.Samples[i]*GAME_HEIGHT/2+GAME_HEIGHT/2))
		t.ApplyGraphic(game.Graphics().MakeGraphic(1, 1, ic.RGBA{0, 255, 30, 255}, cache.CacheOptions{}))
		s.Add(t)
	}
}

func (s *PlayState) Update(elapsed float64) error {
	err := s.BaseGxlState.Update(elapsed)
	if err != nil {
		return err
	}

	return nil
}
