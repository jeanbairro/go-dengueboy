package building_test

import (
	"app/internal/modules/building"
	"app/internal/modules/room"
	"testing"

	"github.com/stretchr/testify/require"
)

// twoConnectedRooms builds a spawn room with:
//   - a door East, reciprocated by a room at {1,0} with a door West (the
//     "successful move" and "neighbor has no matching door" paths use this).
//   - a door South, with NO room placed at {0,1}. This isolates the
//     "neighbor does not exist" case: the door check passes, so only the
//     exists check can fail TryMove(South).
//   - NO door North, but WITH a room placed at {0,-1} that has a
//     reciprocal door South. This isolates the "current room has no door"
//     case: the neighbor genuinely exists (and would reciprocate), so only
//     the HasDoor check can fail TryMove(North).
func twoConnectedRooms() *building.Building {
	spawn := &room.Room{Doors: map[room.Side]bool{room.East: true, room.South: true}, Cleared: true}
	east := &room.Room{Doors: map[room.Side]bool{room.West: true}, Cleared: true}
	north := &room.Room{Doors: map[room.Side]bool{room.South: true}, Cleared: true}
	return &building.Building{
		Rooms: map[building.Coord]*room.Room{
			{0, 0}:  spawn,
			{1, 0}:  east,
			{0, -1}: north,
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
	// North: spawn has no door, but a neighbor genuinely exists at {0,-1}
	// (with a reciprocal door). Only the HasDoor check can fail this.
	b := twoConnectedRooms()
	ok := b.TryMove(room.North)
	require.False(t, ok)
	require.Equal(t, building.Coord{0, 0}, b.Current)
}

func TestTryMoveFailsWhenNeighborRoomDoesNotExist(t *testing.T) {
	// South: spawn has a door, but no room exists at {0,1}. Only the
	// neighbor-exists check can fail this.
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

func TestNewMarksStartRoomAsVisited(t *testing.T) {
	spawn := &room.Room{Cleared: true}
	b := building.New(map[building.Coord]*room.Room{{0, 0}: spawn}, building.Coord{0, 0})
	require.True(t, b.Visited[building.Coord{0, 0}])
	require.Len(t, b.Visited, 1)
}

func TestTryMoveMarksDestinationAsVisited(t *testing.T) {
	b := twoConnectedRooms() // já existe em building_test.go
	b.Visited = map[building.Coord]bool{{0, 0}: true}

	b.TryMove(room.East)

	require.True(t, b.Visited[building.Coord{1, 0}])
}
