package minimap_test

import (
	"app/internal/modules/building"
	"app/internal/modules/minimap"
	"app/internal/modules/room"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLayoutCellsOnlyIncludesVisitedRooms(t *testing.T) {
	b := &building.Building{
		Rooms: map[building.Coord]*room.Room{
			{X: 0, Y: 0}: {}, {X: 1, Y: 0}: {},
		},
		Current: building.Coord{X: 0, Y: 0},
		Visited: map[building.Coord]bool{{X: 0, Y: 0}: true},
	}

	cells := minimap.LayoutCells(b)

	require.Len(t, cells, 1)
	require.Equal(t, building.Coord{X: 0, Y: 0}, cells[0].Coord)
}

func TestLayoutCellsFlagsExactlyTheCurrentRoom(t *testing.T) {
	b := &building.Building{
		Rooms: map[building.Coord]*room.Room{
			{X: 0, Y: 0}: {}, {X: 1, Y: 0}: {},
		},
		Current: building.Coord{X: 1, Y: 0},
		Visited: map[building.Coord]bool{{X: 0, Y: 0}: true, {X: 1, Y: 0}: true},
	}

	cells := minimap.LayoutCells(b)

	currentCount := 0
	for _, c := range cells {
		if c.Current {
			currentCount++
			require.Equal(t, building.Coord{X: 1, Y: 0}, c.Coord)
		}
	}
	require.Equal(t, 1, currentCount)
}

func TestLayoutCellsPositionsRelativeToFullBoundingBox(t *testing.T) {
	b := &building.Building{
		Rooms: map[building.Coord]*room.Room{
			{X: 0, Y: 0}: {}, {X: 1, Y: 0}: {}, {X: 2, Y: 0}: {},
		},
		Current: building.Coord{X: 2, Y: 0},
		// {1,0} exists in Rooms but is NOT visited — the bounding box
		// must still account for it.
		Visited: map[building.Coord]bool{{X: 0, Y: 0}: true, {X: 2, Y: 0}: true},
	}

	cells := minimap.LayoutCells(b)
	require.Len(t, cells, 2)

	byCoord := map[building.Coord]minimap.Cell{}
	for _, c := range cells {
		byCoord[c.Coord] = c
	}

	require.Equal(t, minimap.MarginX, byCoord[building.Coord{X: 0, Y: 0}].ScreenX)
	require.Equal(t, minimap.MarginY, byCoord[building.Coord{X: 0, Y: 0}].ScreenY)
	require.Equal(t, minimap.MarginX+2*(minimap.CellSize+minimap.CellGap), byCoord[building.Coord{X: 2, Y: 0}].ScreenX)
	require.Equal(t, minimap.MarginY, byCoord[building.Coord{X: 2, Y: 0}].ScreenY)
}

func TestLayoutCellsWithSingleRoomBuilding(t *testing.T) {
	b := &building.Building{
		Rooms:   map[building.Coord]*room.Room{{X: 0, Y: 0}: {}},
		Current: building.Coord{X: 0, Y: 0},
		Visited: map[building.Coord]bool{{X: 0, Y: 0}: true},
	}

	cells := minimap.LayoutCells(b)

	require.Len(t, cells, 1)
	require.Equal(t, minimap.MarginX, cells[0].ScreenX)
	require.Equal(t, minimap.MarginY, cells[0].ScreenY)
	require.True(t, cells[0].Current)
}
