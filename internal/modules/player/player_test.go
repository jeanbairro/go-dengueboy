package player_test

import (
	"app/internal/modules/collision"
	"app/internal/modules/geom"
	"app/internal/modules/maptile"
	"app/internal/modules/player"
	"testing"

	input "github.com/quasilyte/ebitengine-input"
	"github.com/stretchr/testify/require"
)

func TestPlayerMovement(t *testing.T) {
	mapTile, error := maptile.New()
	require.NoError(t, error)
	require.NotNil(t, mapTile)

	testPlayer, error := player.New(&input.Handler{}, &collision.CollisionSystem{})
	require.NoError(t, error)
	require.NotNil(t, testPlayer)

	testPlayer.Position = geom.Position{X: 2, Y: 2}
	testPlayer.PreviousPosition = geom.Position{X: 2, Y: 2}

	testPlayer.CurrentAction = player.MoveLeft
	testPlayer.Update(mapTile)

	require.Equal(t, geom.Position{X: 1, Y: 2}, testPlayer.Position)
}
