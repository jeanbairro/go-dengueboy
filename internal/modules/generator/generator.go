package generator

import (
	"app/internal/modules/building"
	"app/internal/modules/maptile"
	"app/internal/modules/room"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const maxConsecutiveFailures = 200

func Generate(roomCount int, rng *rand.Rand, images map[maptile.MapTile]*ebiten.Image) *building.Building {
	doors := map[building.Coord]map[room.Side]bool{}
	placed := []building.Coord{{X: 0, Y: 0}}

	sides := []room.Side{room.North, room.South, room.East, room.West}

	failures := 0
	for len(placed) < roomCount && failures < maxConsecutiveFailures {
		origin := placed[rng.Intn(len(placed))]
		side := sides[rng.Intn(len(sides))]
		neighborCoord := building.Neighbor(origin, side)

		alreadyPlaced := false
		for _, c := range placed {
			if c == neighborCoord {
				alreadyPlaced = true
				break
			}
		}
		if alreadyPlaced {
			failures++
			continue
		}

		placed = append(placed, neighborCoord)

		if doors[origin] == nil {
			doors[origin] = map[room.Side]bool{}
		}
		doors[origin][side] = true

		if doors[neighborCoord] == nil {
			doors[neighborCoord] = map[room.Side]bool{}
		}
		doors[neighborCoord][room.Opposite(side)] = true

		failures = 0
	}

	rooms := map[building.Coord]*room.Room{}
	for _, coord := range placed {
		rooms[coord] = &room.Room{
			Map:     &maptile.Map{Tiles: room.NewBlankRoomTiles(doors[coord]), Images: images},
			Doors:   doors[coord],
			Cleared: true,
		}
	}

	return &building.Building{Rooms: rooms, Current: building.Coord{X: 0, Y: 0}}
}
