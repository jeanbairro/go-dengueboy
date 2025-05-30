package game

const (
	Empty     MapTile = 0
	Wall      MapTile = 1
	WaterTire MapTile = 2
	WaterVase MapTile = 3
	EmptyTire MapTile = 4
	Player    MapTile = 5
)

type (
	MapTile int

	Map struct {
		Tiles [][]int
	}
)

func (m *Map) SetInitialMap() {
	m.Tiles = [][]int{
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 5, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 1, 0, 0, 0, 0, 3, 0, 4, 1},
		{1, 0, 0, 0, 3, 0, 0, 0, 0, 0, 1},
		{1, 2, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 2, 0, 0, 1},
		{1, 0, 2, 0, 0, 0, 0, 0, 2, 0, 1},
		{1, 0, 0, 0, 3, 0, 0, 0, 0, 0, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	}
}
