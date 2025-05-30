package game

import (
	"app/internal/modules/player"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type (
	Game struct {
		Map    Map
		Player player.Player
	}
)

func (g *Game) Update(*ebiten.Image) error {
	keys := inpututil.JustPressedTouchIDs()
	g.Player.Update(keys[0])
	g.Map.Update(g.Player.Position)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.renderMap(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 640
}
