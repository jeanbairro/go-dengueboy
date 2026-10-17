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

func NormalizeDiagonal(dx, dy, speed float64) Position {
	if dx != 0 && dy != 0 {
		const invSqrt2 = 0.7071067811865476
		dx *= invSqrt2
		dy *= invSqrt2
	}
	return Position{X: dx * speed, Y: dy * speed}
}
