package main

import (
	"app/internal/modules/game"
	"log"

	"github.com/hajimehoshi/ebiten"
)

type Game struct{}

func main() {
	ebiten.SetWindowSize(800, 640)
	ebiten.SetWindowTitle("Hello, World!")
	game := game.Game{}
	game.Map.SetInitialMap()
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
