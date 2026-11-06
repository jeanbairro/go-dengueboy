package room

import (
	"app/internal/modules/collision"
	"app/internal/modules/geom"
	"app/internal/modules/maptile"
)

const (
	North Side = iota
	South
	East
	West
)

const (
	DoorColumnNS   = 12
	DoorRowNorth   = 0
	DoorRowSouth   = 19
	DoorRowEW      = 9
	DoorColumnWest = 0
	DoorColumnEast = 24
)

type (
	Side int

	Room struct {
		Map     *maptile.Map
		Doors   map[Side]bool
		Cleared bool
	}
)

func Opposite(side Side) Side {
	switch side {
	case North:
		return South
	case South:
		return North
	case East:
		return West
	case West:
		return East
	default:
		return side
	}
}

func (r *Room) HasDoor(side Side) bool {
	return r.Doors[side]
}

func NewBlankRoomTiles(doors map[Side]bool) [][]int {
	const width, height = 25, 20

	tiles := make([][]int, height)
	for y := range tiles {
		row := make([]int, width)
		for x := range row {
			if y == 0 || y == height-1 || x == 0 || x == width-1 {
				row[x] = int(maptile.Wall)
			} else {
				row[x] = int(maptile.Empty)
			}
		}
		tiles[y] = row
	}

	if doors[North] {
		tiles[DoorRowNorth][DoorColumnNS] = int(maptile.Door)
	}
	if doors[South] {
		tiles[DoorRowSouth][DoorColumnNS] = int(maptile.Door)
	}
	if doors[West] {
		tiles[DoorRowEW][DoorColumnWest] = int(maptile.Door)
	}
	if doors[East] {
		tiles[DoorRowEW][DoorColumnEast] = int(maptile.Door)
	}

	return tiles
}

func doorTile(side Side) (column, row int) {
	switch side {
	case North:
		return DoorColumnNS, DoorRowNorth
	case South:
		return DoorColumnNS, DoorRowSouth
	case East:
		return DoorColumnEast, DoorRowEW
	default:
		return DoorColumnWest, DoorRowEW
	}
}

func DoorZone(side Side) collision.Object {
	column, row := doorTile(side)
	return collision.Object{
		Position: geom.Position{X: float64(column * maptile.ImageSize), Y: float64(row * maptile.ImageSize)},
		Width:    maptile.ImageSize,
		Height:   maptile.ImageSize,
	}
}

func EntryPosition(side Side) geom.Position {
	column, row := doorTile(side)
	switch side {
	case North:
		row++
	case South:
		row--
	case East:
		column--
	default:
		column++
	}
	return geom.Position{X: float64(column * maptile.ImageSize), Y: float64(row * maptile.ImageSize)}
}
