package game

import (
	"app/internal/modules/maptile"
	"app/internal/modules/player"
	"log"

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
	g.updateMap()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if err := g.renderMap(screen); err != nil {
		log.Println("error rendering map:", err)
	}
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

func (g *Game) updateMap() {
	g.Map.Update(g.Player.Position, g.Player.PreviousPosition)
}
