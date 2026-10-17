package entity

import "github.com/hajimehoshi/ebiten/v2"

type (
	Entity interface {
		Update()
		Draw(screen *ebiten.Image)
	}

	Manager struct {
		entities []Entity
	}
)

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) Add(e Entity) {
	m.entities = append(m.entities, e)
}

func (m *Manager) Update() {
	for _, e := range m.entities {
		e.Update()
	}
}

func (m *Manager) Draw(screen *ebiten.Image) {
	for _, e := range m.entities {
		e.Draw(screen)
	}
}

func (m *Manager) Len() int {
	return len(m.entities)
}
