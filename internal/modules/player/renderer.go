package player

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	input "github.com/quasilyte/ebitengine-input"
)

const (
	spriteSize = 32
)

func (p *Player) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(p.Position.X*spriteSize), float64(p.Position.Y*spriteSize))
	screen.DrawImage(p.CurrentSprite, op)
}

func loadSprites() (map[input.Action]*ebiten.Image, error) {
	charImage, _, err := ebitenutil.NewImageFromFile("assets/images/char.png")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return map[input.Action]*ebiten.Image{
		NoAction:  charImage.SubImage(image.Rect(0, 0, spriteSize, spriteSize)).(*ebiten.Image),
		MoveLeft:  charImage.SubImage(image.Rect(spriteSize, 0, spriteSize*2, spriteSize)).(*ebiten.Image),
		MoveRight: charImage.SubImage(image.Rect(spriteSize*2, 0, spriteSize*3, spriteSize)).(*ebiten.Image),
		MoveUp:    charImage.SubImage(image.Rect(spriteSize*3, 0, spriteSize*4, spriteSize)).(*ebiten.Image),
		MoveDown:  charImage.SubImage(image.Rect(spriteSize*4, 0, spriteSize*5, spriteSize)).(*ebiten.Image),
	}, nil
}
