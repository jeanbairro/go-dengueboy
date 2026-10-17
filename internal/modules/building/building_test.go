package building_test

import (
	"app/internal/modules/building"
	"app/internal/modules/room"
	"testing"

	"github.com/stretchr/testify/require"
)

func twoConnectedRooms() *building.Building {
	spawn := &room.Room{Doors: map[room.Side]bool{room.East: true}, Cleared: true}
	east := &room.Room{Doors: map[room.Side]bool{room.West: true}, Cleared: true}
	return &building.Building{
		Rooms: map[building.Coord]*room.Room{
			{0, 0}: spawn,
			{1, 0}: east,
		},
		Current: building.Coord{0, 0},
	}
}

func TestTryMoveSucceedsWhenBothSidesHaveMatchingDoors(t *testing.T) {
	b := twoConnectedRooms()
	ok := b.TryMove(room.East)
	require.True(t, ok)
	require.Equal(t, building.Coord{1, 0}, b.Current)
}

func TestTryMoveFailsWhenCurrentRoomHasNoDoorThatSide(t *testing.T) {
	b := twoConnectedRooms()
	ok := b.TryMove(room.North)
	require.False(t, ok)
	require.Equal(t, building.Coord{0, 0}, b.Current)
}

func TestTryMoveFailsWhenNeighborRoomDoesNotExist(t *testing.T) {
	b := twoConnectedRooms()
	ok := b.TryMove(room.South)
	require.False(t, ok)
}

func TestTryMoveFailsWhenNeighborHasNoMatchingDoor(t *testing.T) {
	spawn := &room.Room{Doors: map[room.Side]bool{room.East: true}, Cleared: true}
	eastWithoutWestDoor := &room.Room{Doors: map[room.Side]bool{}, Cleared: true}
	b := &building.Building{
		Rooms: map[building.Coord]*room.Room{
			{0, 0}: spawn,
			{1, 0}: eastWithoutWestDoor,
		},
		Current: building.Coord{0, 0},
	}
	ok := b.TryMove(room.East)
	require.False(t, ok, "neighbor must reciprocate the door on the opposite side")
}

func TestTryMoveFailsWhenCurrentRoomIsNotCleared(t *testing.T) {
	spawn := &room.Room{Doors: map[room.Side]bool{room.East: true}, Cleared: false}
	east := &room.Room{Doors: map[room.Side]bool{room.West: true}, Cleared: true}
	b := &building.Building{
		Rooms: map[building.Coord]*room.Room{
			{0, 0}: spawn,
			{1, 0}: east,
		},
		Current: building.Coord{0, 0},
	}
	ok := b.TryMove(room.East)
	require.False(t, ok)
}
