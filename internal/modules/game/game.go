package game

import (
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
		InputSystem input.System
	}
)

func New(maptile *maptile.Map, player *player.Player, inputSystem input.System) *Game {
	return &Game{
		Map:         maptile,
		Player:      player,
		InputSystem: inputSystem,
	}
}

func (g *Game) Update() error {
	g.updateInput()
	g.updatePlayer()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Map.Draw(screen)
	g.Player.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) updateInput() {
	g.InputSystem.Update()
}

func (g *Game) updatePlayer() {
	g.Player.Update(g.Map)
}
