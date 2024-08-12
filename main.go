package main

import (
	"embed"

	"github.com/GixelEngine/gixel-engine/gixel"
)

const GAME_WIDTH = 1280
const GAME_HEIGHT = 960

//go:embed assets
var assets embed.FS

func main() {
	gixel.NewGame(1280, 960, "Gixel WAV", &assets, &PlayState{}, 1).Run()

	// format, signal, err := wav.Decode("assets/440.wav")
	// if err != nil {
	// 	panic(err)
	// }

	// marshalFormat, _ := json.Marshal(format)
	// fmt.Println(string(marshalFormat))

	// // marshalSignal, _ := json.Marshal(signal)
	// fmt.Println(signal)
}
