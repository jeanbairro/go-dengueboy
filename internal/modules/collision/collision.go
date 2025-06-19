package collision

import (
	"app/internal/modules/geom"
	"app/internal/modules/maptile"
)

type (
	Object struct {
		Width    float64
		Height   float64
		Position geom.Position
	}

	CollisionSystem struct {
		Objects []Object
	}
)

func New(mapTile *maptile.Map) (*CollisionSystem, error) {
	objects := createCollisionObjects(mapTile)
	return &CollisionSystem{Objects: objects}, nil
}

func (c *CollisionSystem) HasCollision(obj Object) bool {
	for _, other := range c.Objects {
		if isColliding(obj, other) {
			return true
		}
	}
	return false
}

func isColliding(a, b Object) bool {
	return a.Position.X < b.Position.X+b.Width &&
		a.Position.X+a.Width > b.Position.X &&
		a.Position.Y < b.Position.Y+b.Height &&
		a.Position.Y+a.Height > b.Position.Y
}

func createCollisionObjects(mapTile *maptile.Map) []Object {
	const imageSize = float64(maptile.ImageSize)
	objects := []Object{}

	for y, row := range mapTile.Tiles {
		for x, tileValue := range row {
			tile := maptile.MapTile(tileValue)
			if tile == maptile.Wall {
				objects = append(objects, Object{
					Width:  imageSize,
					Height: imageSize,
					Position: geom.Position{
						X: float64(x) * imageSize,
						Y: float64(y) * imageSize,
					},
				})
			}
		}
	}
	return objects
}
