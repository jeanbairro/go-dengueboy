package main

import (
	"app/internal/modules/collision"
	"app/internal/modules/game"
	"app/internal/modules/generator"
	"app/internal/modules/maptile"
	"app/internal/modules/player"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
)

type Game struct{}

func main() {
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	game, error := initializeGame()
	if error != nil {
		log.Fatal(error)
		return
	}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func initializeGame() (*game.Game, error) {
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
	mapTile, err := maptile.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create map tile: %w", err)
	}
	initialBuilding := generator.Generate(6, rand.New(rand.NewSource(time.Now().UnixNano())), mapTile.Images)
	collisionSystem, err := collision.New(initialBuilding.CurrentRoom().Map)
	if err != nil {
		return nil, fmt.Errorf("failed to create collision system: %w", err)
	}
	player, err := player.New(inputSystem.NewHandler(0, keymap), collisionSystem)
	if err != nil {
		return nil, fmt.Errorf("failed to create player: %w", err)
	}
	return game.New(initialBuilding, player, inputSystem), nil
}
