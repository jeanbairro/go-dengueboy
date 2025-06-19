package geom

import (
	"math"
)

type Position struct {
	X, Y float64
}

func (p *Position) Move(dx, dy float64) {
	p.X += dx
	p.Y += dy
}

func (p *Position) DistanceTo(otherPosition Position) float64 {
	distanceX := math.Abs(p.X - otherPosition.X)
	distanceY := math.Abs(p.Y - otherPosition.Y)
	return distanceX + distanceY
}
