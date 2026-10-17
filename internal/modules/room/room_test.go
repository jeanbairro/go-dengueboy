package room_test

import (
	"app/internal/modules/collision"
	"app/internal/modules/geom"
	"app/internal/modules/maptile"
	"app/internal/modules/room"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOppositeSide(t *testing.T) {
	require.Equal(t, room.South, room.Opposite(room.North))
	require.Equal(t, room.North, room.Opposite(room.South))
	require.Equal(t, room.West, room.Opposite(room.East))
	require.Equal(t, room.East, room.Opposite(room.West))
}

func TestRoomHasDoor(t *testing.T) {
	r := &room.Room{Doors: map[room.Side]bool{room.North: true}}
	require.True(t, r.HasDoor(room.North))
	require.False(t, r.HasDoor(room.South))
}

func TestNewBlankRoomTilesHasBorderWalls(t *testing.T) {
	tiles := room.NewBlankRoomTiles(map[room.Side]bool{})
	require.Len(t, tiles, 20)
	require.Len(t, tiles[0], 25)
	require.Equal(t, int(maptile.Wall), tiles[0][0], "top-left corner must stay a wall")
	require.Equal(t, int(maptile.Empty), tiles[10][12], "room interior must be empty")
}

func TestNewBlankRoomTilesOpensDoorOnRequestedSide(t *testing.T) {
	tiles := room.NewBlankRoomTiles(map[room.Side]bool{room.North: true})
	require.Equal(t, int(maptile.Door), tiles[0][12], "north door must open at the fixed column")
	require.Equal(t, int(maptile.Wall), tiles[19][12], "south side must stay closed when not requested")
}

func TestDoorZonePositions(t *testing.T) {
	require.Equal(t, collision.Object{Position: geom.Position{X: 384, Y: 0}, Width: 32, Height: 32}, room.DoorZone(room.North))
	require.Equal(t, collision.Object{Position: geom.Position{X: 384, Y: 608}, Width: 32, Height: 32}, room.DoorZone(room.South))
	require.Equal(t, collision.Object{Position: geom.Position{X: 768, Y: 288}, Width: 32, Height: 32}, room.DoorZone(room.East))
	require.Equal(t, collision.Object{Position: geom.Position{X: 0, Y: 288}, Width: 32, Height: 32}, room.DoorZone(room.West))
}

func TestEntryPositions(t *testing.T) {
	require.Equal(t, geom.Position{X: 384, Y: 32}, room.EntryPosition(room.North))
	require.Equal(t, geom.Position{X: 384, Y: 576}, room.EntryPosition(room.South))
	require.Equal(t, geom.Position{X: 736, Y: 288}, room.EntryPosition(room.East))
	require.Equal(t, geom.Position{X: 32, Y: 288}, room.EntryPosition(room.West))
}
