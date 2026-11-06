package building

import "app/internal/modules/room"

type (
	Coord struct {
		X, Y int
	}

	Building struct {
		Rooms   map[Coord]*room.Room
		Current Coord
		Visited map[Coord]bool
	}
)

func New(rooms map[Coord]*room.Room, start Coord) *Building {
	return &Building{Rooms: rooms, Current: start, Visited: map[Coord]bool{start: true}}
}

func Neighbor(c Coord, side room.Side) Coord {
	switch side {
	case room.North:
		return Coord{c.X, c.Y - 1}
	case room.South:
		return Coord{c.X, c.Y + 1}
	case room.East:
		return Coord{c.X + 1, c.Y}
	case room.West:
		return Coord{c.X - 1, c.Y}
	default:
		return c
	}
}

func (b *Building) CurrentRoom() *room.Room {
	return b.Rooms[b.Current]
}

func (b *Building) TryMove(side room.Side) bool {
	current := b.CurrentRoom()
	if !current.HasDoor(side) {
		return false
	}
	if !current.Cleared {
		return false
	}

	neighborCoord := Neighbor(b.Current, side)
	neighbor, exists := b.Rooms[neighborCoord]
	if !exists {
		return false
	}
	if !neighbor.HasDoor(room.Opposite(side)) {
		return false
	}

	b.Current = neighborCoord
	if b.Visited == nil {
		b.Visited = map[Coord]bool{}
	}
	b.Visited[neighborCoord] = true
	return true
}
