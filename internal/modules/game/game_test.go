package game_test

import (
	"app/internal/modules/building"
	"app/internal/modules/collision"
	"app/internal/modules/entity"
	"app/internal/modules/game"
	"app/internal/modules/geom"
	"app/internal/modules/maptile"
	"app/internal/modules/player"
	"app/internal/modules/room"
	"testing"

	input "github.com/quasilyte/ebitengine-input"
	"github.com/stretchr/testify/require"
)

func TestNewRegistersPlayerInEntities(t *testing.T) {
	p := &player.Player{
		Position:        geom.Position{X: 0, Y: 0},
		CollisionSystem: &collision.CollisionSystem{},
	}
	spawn := &room.Room{Map: &maptile.Map{}, Doors: map[room.Side]bool{}}
	b := &building.Building{Rooms: map[building.Coord]*room.Room{{0, 0}: spawn}, Current: building.Coord{0, 0}}

	g := game.New(b, p, input.System{})

	require.Equal(t, 1, g.Entities.Len())
	require.Equal(t, p, g.Player)
}

func TestCheckRoomTransitionMovesPlayerAndRebuildsCollision(t *testing.T) {
	spawnTiles := room.NewBlankRoomTiles(map[room.Side]bool{room.East: true})
	eastTiles := room.NewBlankRoomTiles(map[room.Side]bool{room.West: true})
	spawn := &room.Room{Map: &maptile.Map{Tiles: spawnTiles}, Doors: map[room.Side]bool{room.East: true}, Cleared: true}
	east := &room.Room{Map: &maptile.Map{Tiles: eastTiles}, Doors: map[room.Side]bool{room.West: true}, Cleared: true}
	b := &building.Building{
		Rooms:   map[building.Coord]*room.Room{{0, 0}: spawn, {1, 0}: east},
		Current: building.Coord{0, 0},
	}

	cs := &collision.CollisionSystem{}
	cs.Rebuild(spawn.Map)
	p := &player.Player{Position: room.DoorZone(room.East).Position, CollisionSystem: cs}

	g := &game.Game{Building: b, Player: p, Entities: entity.NewManager()}
	g.Entities.Add(p)

	g.CheckRoomTransition()

	require.Equal(t, building.Coord{1, 0}, b.Current)
	require.Equal(t, room.EntryPosition(room.West), p.Position)
}

func TestCheckRoomTransitionDoesNothingWithoutADoor(t *testing.T) {
	windowless := &room.Room{Map: &maptile.Map{Tiles: room.NewBlankRoomTiles(nil)}, Doors: map[room.Side]bool{}, Cleared: true}
	b := &building.Building{Rooms: map[building.Coord]*room.Room{{0, 0}: windowless}, Current: building.Coord{0, 0}}
	cs := &collision.CollisionSystem{}
	cs.Rebuild(windowless.Map)
	// Player standing exactly where a door zone would be if one existed:
	p := &player.Player{Position: room.DoorZone(room.North).Position, CollisionSystem: cs}

	g := &game.Game{Building: b, Player: p, Entities: entity.NewManager()}
	g.Entities.Add(p)

	g.CheckRoomTransition()

	require.Equal(t, building.Coord{0, 0}, b.Current)
	require.Equal(t, room.DoorZone(room.North).Position, p.Position)
}
