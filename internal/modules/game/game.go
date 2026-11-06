package game

import (
	"app/internal/modules/entity"
	"app/internal/modules/maptile"
	"app/internal/modules/player"

	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 640
)

type (
	Game struct {
		Map         *maptile.Map
		Player      *player.Player
		Entities    *entity.Manager
		InputSystem input.System
	}
)

func New(maptile *maptile.Map, player *player.Player, inputSystem input.System) *Game {
	entities := entity.NewManager()
	entities.Add(player)

	return &Game{
		Map:         maptile,
		Player:      player,
		Entities:    entities,
		InputSystem: inputSystem,
	}
}

func (g *Game) Update() error {
	g.updateInput()
	g.Entities.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Map.Draw(screen)
	g.Entities.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) updateInput() {
	g.InputSystem.Update()
}
