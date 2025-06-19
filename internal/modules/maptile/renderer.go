package maptile

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	imageSize = 32
)

func (m *Map) Draw(screen *ebiten.Image) error {
	for y := range m.Tiles {
		for x := range m.Tiles[y] {
			tile := MapTile(m.Tiles[y][x])
			op := &ebiten.DrawImageOptions{GeoM: ebiten.GeoM{}}
			op.GeoM.Translate(float64(imageSize*x), float64(imageSize*y))
			screen.DrawImage(m.Images[tile], op)
		}
	}
	return nil
}

func loadImageTiles() (map[MapTile]*ebiten.Image, error) {
	jungle, _, error := ebitenutil.NewImageFromFile("assets/images/jungle.png")
	if error != nil {
		return nil, error
	}
	emptyImage := jungle.SubImage(image.Rect(0, 0, imageSize, imageSize)).(*ebiten.Image)
	wallImage := jungle.SubImage(image.Rect(imageSize, 0, jungle.Bounds().Max.X, imageSize)).(*ebiten.Image)
	return map[MapTile]*ebiten.Image{
		Empty: emptyImage,
		Wall:  wallImage,
	}, nil
}
