package geom

import "math"

type Position struct {
	X, Y int
}

func (p *Position) Move(dx, dy int) {
	p.X += dx
	p.Y += dy
}

func (p *Position) DistanceTo(otherPosition Position) int {
	distanceX := math.Abs(float64(p.X) - float64(otherPosition.X))
	distanceY := math.Abs(float64(p.Y) - float64(otherPosition.Y))
	return int(distanceX + distanceY)
}
