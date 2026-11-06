package minimap

import (
	"image/color"

	"app/internal/modules/building"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	CellSize = 10.0
	CellGap  = 2.0
	MarginX  = 12.0
	MarginY  = 12.0
)

type Cell struct {
	Coord   building.Coord
	ScreenX float64
	ScreenY float64
	Current bool
}

func LayoutCells(b *building.Building) []Cell {
	minX, minY, found := 0, 0, false
	for c := range b.Rooms {
		if !found || c.X < minX {
			minX = c.X
		}
		if !found || c.Y < minY {
			minY = c.Y
		}
		found = true
	}

	cells := make([]Cell, 0, len(b.Visited))
	for c, visited := range b.Visited {
		if !visited {
			continue
		}
		cells = append(cells, Cell{
			Coord:   c,
			ScreenX: MarginX + float64(c.X-minX)*(CellSize+CellGap),
			ScreenY: MarginY + float64(c.Y-minY)*(CellSize+CellGap),
			Current: c == b.Current,
		})
	}
	return cells
}

func Draw(screen *ebiten.Image, b *building.Building) {
	for _, cell := range LayoutCells(b) {
		clr := color.RGBA{150, 150, 150, 255}
		if cell.Current {
			clr = color.RGBA{255, 255, 255, 255}
		}
		vector.DrawFilledRect(
			screen,
			float32(cell.ScreenX),
			float32(cell.ScreenY),
			float32(CellSize),
			float32(CellSize),
			clr,
			false,
		)
	}
}
