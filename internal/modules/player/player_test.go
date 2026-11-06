package player_test

import (
	"app/internal/modules/collision"
	"app/internal/modules/geom"
	"app/internal/modules/player"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPlayerMoveFreelyWithNoCollision(t *testing.T) {
	p := &player.Player{
		Position:        geom.Position{X: 100, Y: 100},
		CollisionSystem: &collision.CollisionSystem{},
	}

	p.Move(geom.Position{X: 2.56, Y: 0})

	require.Equal(t, geom.Position{X: 102.56, Y: 100}, p.Position)
	require.Equal(t, geom.Position{X: 100, Y: 100}, p.PreviousPosition)
	require.Equal(t, player.MoveRight, p.CurrentAction)
}

func TestPlayerMoveSlidesAlongWallWhenOneAxisBlocked(t *testing.T) {
	wall := collision.Object{Position: geom.Position{X: 132, Y: 100}, Width: 32, Height: 32}
	p := &player.Player{
		Position:        geom.Position{X: 100, Y: 100},
		CollisionSystem: &collision.CollisionSystem{Objects: []collision.Object{wall}},
	}

	p.Move(geom.Position{X: 5, Y: 5})

	require.Equal(t, geom.Position{X: 100, Y: 105}, p.Position)
}

func TestPlayerMoveBlockedOnBothAxesStaysInPlace(t *testing.T) {
	wall := collision.Object{Position: geom.Position{X: 100, Y: 116}, Width: 48, Height: 48}
	p := &player.Player{
		Position:        geom.Position{X: 100, Y: 100},
		CollisionSystem: &collision.CollisionSystem{Objects: []collision.Object{wall}},
	}

	p.Move(geom.Position{X: 5, Y: 5})

	require.Equal(t, geom.Position{X: 100, Y: 100}, p.Position)
}

func TestPlayerMoveZeroVectorSetsIdleFacing(t *testing.T) {
	p := &player.Player{
		Position:        geom.Position{X: 50, Y: 50},
		CurrentAction:   player.MoveRight,
		CollisionSystem: &collision.CollisionSystem{},
	}

	p.Move(geom.Position{X: 0, Y: 0})
	p.Move(geom.Position{X: 0, Y: 0})

	require.Equal(t, geom.Position{X: 50, Y: 50}, p.Position)
	require.Equal(t, player.NoAction, p.CurrentAction)
}
