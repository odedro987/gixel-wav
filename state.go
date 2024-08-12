package main

import (
	"github.com/GixelEngine/gixel-engine/gixel"
)

type PlayState struct {
	gixel.BaseGxlState
}

func (s *PlayState) Init(game *gixel.GxlGame) {
	s.BaseGxlState.Init(game)
}

func (s *PlayState) Update(elapsed float64) error {
	err := s.BaseGxlState.Update(elapsed)
	if err != nil {
		return err
	}

	return nil
}
