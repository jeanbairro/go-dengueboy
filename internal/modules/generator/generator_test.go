package generator_test

import (
	"app/internal/modules/building"
	"app/internal/modules/generator"
	"app/internal/modules/room"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateProducesRequestedRoomCount(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	b := generator.Generate(6, rng, nil)
	require.Len(t, b.Rooms, 6)
}

func TestGenerateStartsAtOriginAndSpawnIsCurrent(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	b := generator.Generate(6, rng, nil)
	require.Contains(t, b.Rooms, building.Coord{X: 0, Y: 0})
	require.Equal(t, building.Coord{X: 0, Y: 0}, b.Current)
}

func TestGenerateSingleRoomDoesNotHang(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	b := generator.Generate(1, rng, nil)
	require.Len(t, b.Rooms, 1)
}

func TestGenerateAllRoomsAreClearedByDefault(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	b := generator.Generate(6, rng, nil)
	for _, r := range b.Rooms {
		require.True(t, r.Cleared)
	}
}

func allSides() []room.Side {
	return []room.Side{room.North, room.South, room.East, room.West}
}

func TestGenerateDoorsAreReciprocal(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	b := generator.Generate(10, rng, nil)
	for coord, r := range b.Rooms {
		for _, side := range allSides() {
			if !r.HasDoor(side) {
				continue
			}
			neighborCoord := building.Neighbor(coord, side)
			neighbor, exists := b.Rooms[neighborCoord]
			require.True(t, exists, "room at %+v has a door %v with no neighbor room", coord, side)
			require.True(t, neighbor.HasDoor(room.Opposite(side)), "neighbor at %+v does not reciprocate the door", neighborCoord)
		}
	}
}

func TestGenerateAllRoomsReachableFromSpawn(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	b := generator.Generate(10, rng, nil)

	visited := map[building.Coord]bool{{X: 0, Y: 0}: true}
	queue := []building.Coord{{X: 0, Y: 0}}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		r := b.Rooms[current]
		for _, side := range allSides() {
			if !r.HasDoor(side) {
				continue
			}
			next := building.Neighbor(current, side)
			if !visited[next] {
				visited[next] = true
				queue = append(queue, next)
			}
		}
	}
	require.Len(t, visited, len(b.Rooms), "every generated room must be reachable from spawn")
}

func TestGenerateIsDeterministicForTheSameSeed(t *testing.T) {
	a := generator.Generate(10, rand.New(rand.NewSource(42)), nil)
	b := generator.Generate(10, rand.New(rand.NewSource(42)), nil)
	require.Equal(t, a.Rooms, b.Rooms)
}
