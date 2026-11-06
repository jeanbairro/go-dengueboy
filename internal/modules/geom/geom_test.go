package geom_test

import (
	"app/internal/modules/geom"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeDiagonalCardinalUnchanged(t *testing.T) {
	result := geom.NormalizeDiagonal(1, 0, 2.56)
	require.Equal(t, geom.Position{X: 2.56, Y: 0}, result)
}

func TestNormalizeDiagonalScalesDownDiagonalMovement(t *testing.T) {
	result := geom.NormalizeDiagonal(1, -1, 2.56)
	require.InDelta(t, 1.8101933598, result.X, 0.0001)
	require.InDelta(t, -1.8101933598, result.Y, 0.0001)
}

func TestNormalizeDiagonalOpposingInputCancelsToZero(t *testing.T) {
	result := geom.NormalizeDiagonal(0, 0, 2.56)
	require.Equal(t, geom.Position{X: 0, Y: 0}, result)
}
