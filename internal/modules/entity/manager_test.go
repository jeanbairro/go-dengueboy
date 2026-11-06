package entity_test

import (
	"app/internal/modules/entity"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/stretchr/testify/require"
)

type fakeEntity struct {
	updateCalls int
	drawCalls   int
}

func (f *fakeEntity) Update()                   { f.updateCalls++ }
func (f *fakeEntity) Draw(screen *ebiten.Image) { f.drawCalls++ }

func TestManagerUpdateCallsEveryRegisteredEntity(t *testing.T) {
	m := entity.NewManager()
	a, b := &fakeEntity{}, &fakeEntity{}
	m.Add(a)
	m.Add(b)

	m.Update()

	require.Equal(t, 1, a.updateCalls)
	require.Equal(t, 1, b.updateCalls)
}

func TestManagerDrawCallsEveryRegisteredEntity(t *testing.T) {
	m := entity.NewManager()
	a, b := &fakeEntity{}, &fakeEntity{}
	m.Add(a)
	m.Add(b)

	m.Draw(nil)

	require.Equal(t, 1, a.drawCalls)
	require.Equal(t, 1, b.drawCalls)
}

func TestManagerWithNoEntitiesIsNoop(t *testing.T) {
	m := entity.NewManager()

	require.NotPanics(t, func() {
		m.Update()
		m.Draw(nil)
	})
	require.Equal(t, 0, m.Len())
}

func TestManagerLenReflectsRegistrations(t *testing.T) {
	m := entity.NewManager()
	m.Add(&fakeEntity{})
	m.Add(&fakeEntity{})

	require.Equal(t, 2, m.Len())
}
