package collision_test

import (
	"app/internal/modules/collision"
	"app/internal/modules/geom"
	"app/internal/modules/maptile"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRebuildReplacesObjectsFromNewMap(t *testing.T) {
	c := &collision.CollisionSystem{
		Objects: []collision.Object{{Position: geom.Position{X: 999, Y: 999}, Width: 32, Height: 32}},
	}
	newMap := &maptile.Map{Tiles: [][]int{{1, 0}, {0, 1}}}

	c.Rebuild(newMap)

	require.Len(t, c.Objects, 2, "only the two Wall tiles in newMap become objects")
	require.NotContains(t, c.Objects, collision.Object{Position: geom.Position{X: 999, Y: 999}, Width: 32, Height: 32})
}
