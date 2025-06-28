package player

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	input "github.com/quasilyte/ebitengine-input"
)

const (
	spriteSize        = 32
	spriteCount       = 5
	moveDownImageRow  = 0
	moveLeftImageRow  = 1
	moveUpImageRow    = 2
	moveRightImageRow = 3
)

func (p *Player) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(p.Position.X*spriteSize), float64(p.Position.Y*spriteSize))
	screen.DrawImage(p.CurrentSprite, op)
}

func loadSprites() (map[input.Action][]*ebiten.Image, error) {
	charImage, _, err := ebitenutil.NewImageFromFile("assets/images/char.png")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return map[input.Action][]*ebiten.Image{
		NoAction:  {charImage.SubImage(image.Rect(0, 0, spriteSize, spriteSize)).(*ebiten.Image)},
		MoveDown:  loadActionSprites(charImage, moveDownImageRow),
		MoveLeft:  loadActionSprites(charImage, moveLeftImageRow),
		MoveUp:    loadActionSprites(charImage, moveUpImageRow),
		MoveRight: loadActionSprites(charImage, moveRightImageRow),
	}, nil
}

func loadActionSprites(charImage *ebiten.Image, imageRow int) []*ebiten.Image {
	sprites := make([]*ebiten.Image, spriteCount)
	for i := range spriteCount {
		sprites[i] = charImage.SubImage(image.Rect(spriteSize*i, spriteSize*imageRow, spriteSize*(i+1), spriteSize*(imageRow+1))).(*ebiten.Image)
	}
	return sprites
}
