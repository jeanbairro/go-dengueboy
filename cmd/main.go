package main

import (
	"app/internal/modules/building"
	"app/internal/modules/collision"
	"app/internal/modules/game"
	"app/internal/modules/maptile"
	"app/internal/modules/player"
	"app/internal/modules/room"
	"fmt"
	"log"

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
	initialBuilding := &building.Building{
		Rooms:   buildDemoRooms(mapTile.Images),
		Current: building.Coord{X: 0, Y: 0},
	}
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

func buildDemoRooms(images map[maptile.MapTile]*ebiten.Image) map[building.Coord]*room.Room {
	newRoom := func(doors map[room.Side]bool) *room.Room {
		return &room.Room{
			Map: &maptile.Map{
				Tiles:  room.NewBlankRoomTiles(doors),
				Images: images,
			},
			Doors:   doors,
			Cleared: true,
		}
	}

	return map[building.Coord]*room.Room{
		{X: 0, Y: 0}:  newRoom(map[room.Side]bool{room.East: true}),
		{X: 1, Y: 0}:  newRoom(map[room.Side]bool{room.West: true, room.North: true}),
		{X: 1, Y: -1}: newRoom(map[room.Side]bool{room.South: true}),
	}
}
