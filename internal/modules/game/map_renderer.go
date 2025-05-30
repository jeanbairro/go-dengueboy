package game

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	ImageSize = 32
)

func (g *Game) renderMap(screen *ebiten.Image) {
	mapImages, err := loadMapImageTiles()
	if err != nil {
		log.Fatal(err)
		return
	}
	for y := range g.Map.Tiles {
		for x := range g.Map.Tiles[y] {
			tile := MapTile(g.Map.Tiles[y][x])
			op := &ebiten.DrawImageOptions{GeoM: ebiten.GeoM{}}
			op.GeoM.Translate(float64(ImageSize*x), float64(ImageSize*y))
			screen.DrawImage(mapImages[tile], op)
		}
	}
}

func loadMapImageTiles() (map[MapTile]*ebiten.Image, error) {
	jungle, _, error := ebitenutil.NewImageFromFile("assets/images/jungle.png", ebiten.FilterDefault)
	if error != nil {
		log.Fatal(error)
		return nil, error
	}
	empty := jungle.SubImage(image.Rect(0, 0, ImageSize, ImageSize)).(*ebiten.Image)
	wall := jungle.SubImage(image.Rect(ImageSize, 0, jungle.Bounds().Max.X, ImageSize)).(*ebiten.Image)
	return map[MapTile]*ebiten.Image{
		Empty:  empty,
		Wall:   wall,
		Player: wall,
	}, nil
}
