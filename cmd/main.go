package main

import (
	"app/internal/modules/game"
	"app/internal/modules/maptile"
	"app/internal/modules/player"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
)

type Game struct{}

func main() {
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	game := initializeGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func initializeGame() *game.Game {
	inputSystem := input.System{}
	inputSystem.Init(input.SystemConfig{
		DevicesEnabled: input.AnyDevice,
	})
	keymap := input.Keymap{
		player.MoveLeft:  {input.KeyGamepadLeft, input.KeyLeft, input.KeyA},
		player.MoveRight: {input.KeyGamepadRight, input.KeyRight, input.KeyD},
		player.MoveUp:    {input.KeyGamepadUp, input.KeyUp, input.KeyW},
		player.MoveDown:  {input.KeyGamepadDown, input.KeyDown, input.KeyS},
		player.Interact:  {input.KeyGamepadA, input.KeyEnter, input.KeySpace},
	}
	player := player.New(inputSystem.NewHandler(0, keymap))
	mapTile := maptile.New()
	return game.New(mapTile, player, inputSystem)
}
