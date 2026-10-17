package game

import (
	"app/internal/modules/building"
	"app/internal/modules/collision"
	"app/internal/modules/entity"
	"app/internal/modules/geom"
	"app/internal/modules/maptile"
	"app/internal/modules/player"
	"app/internal/modules/room"

	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 640
)

type (
	Game struct {
		Building    *building.Building
		Player      *player.Player
		Entities    *entity.Manager
		InputSystem input.System
	}
)

func New(building *building.Building, player *player.Player, inputSystem input.System) *Game {
	entities := entity.NewManager()
	entities.Add(player)

	return &Game{
		Building:    building,
		Player:      player,
		Entities:    entities,
		InputSystem: inputSystem,
	}
}

func (g *Game) Update() error {
	g.updateInput()
	g.Entities.Update()
	g.CheckRoomTransition()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Building.CurrentRoom().Map.Draw(screen)
	g.Entities.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) updateInput() {
	g.InputSystem.Update()
}

func (g *Game) CheckRoomTransition() {
	currentRoom := g.Building.CurrentRoom()

	for _, side := range []room.Side{room.North, room.South, room.East, room.West} {
		if !currentRoom.HasDoor(side) {
			continue
		}

		doorZone := collision.CollisionSystem{Objects: []collision.Object{room.DoorZone(side)}}
		if !doorZone.HasCollision(playerObject(g.Player.Position)) {
			continue
		}

		if !g.Building.TryMove(side) {
			continue
		}

		g.Player.Position = room.EntryPosition(room.Opposite(side))
		g.Player.CollisionSystem.Rebuild(g.Building.CurrentRoom().Map)
		return
	}
}

func playerObject(position geom.Position) collision.Object {
	return collision.Object{
		Position: position,
		Width:    maptile.ImageSize,
		Height:   maptile.ImageSize,
	}
}
