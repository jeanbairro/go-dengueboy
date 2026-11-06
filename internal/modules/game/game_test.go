package game_test

import (
	"app/internal/modules/collision"
	"app/internal/modules/game"
	"app/internal/modules/geom"
	"app/internal/modules/maptile"
	"app/internal/modules/player"
	"testing"

	input "github.com/quasilyte/ebitengine-input"
	"github.com/stretchr/testify/require"
)

func TestNewRegistersPlayerInEntities(t *testing.T) {
	p := &player.Player{
		Position:        geom.Position{X: 0, Y: 0},
		CollisionSystem: &collision.CollisionSystem{},
	}

	g := game.New(&maptile.Map{}, p, input.System{})

	require.Equal(t, 1, g.Entities.Len())
	require.Equal(t, p, g.Player)
}
