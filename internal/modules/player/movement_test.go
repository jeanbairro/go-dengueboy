package player

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMovementIntentSingleDirection(t *testing.T) {
	dx, dy := movementIntent(true, false, false, false)
	require.Equal(t, -1.0, dx)
	require.Equal(t, 0.0, dy)
}

func TestMovementIntentOpposingHorizontalCancels(t *testing.T) {
	dx, dy := movementIntent(true, true, false, false)
	require.Equal(t, 0.0, dx)
	require.Equal(t, 0.0, dy)
}

func TestMovementIntentOpposingVerticalCancels(t *testing.T) {
	dx, dy := movementIntent(false, false, true, true)
	require.Equal(t, 0.0, dx)
	require.Equal(t, 0.0, dy)
}

func TestMovementIntentDiagonal(t *testing.T) {
	dx, dy := movementIntent(false, true, true, false)
	require.Equal(t, 1.0, dx)
	require.Equal(t, -1.0, dy)
}
