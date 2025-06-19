package game

import (
	"app/internal/modules/maptile"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	ImageSize = 32
)

func (g *Game) renderMap(screen *ebiten.Image) error {
	mapImages, err := loadMapImageTiles()
	if err != nil {
		return err
	}
	for y := range g.Map.Tiles {
		for x := range g.Map.Tiles[y] {
			tile := maptile.MapTile(g.Map.Tiles[y][x])
			op := &ebiten.DrawImageOptions{GeoM: ebiten.GeoM{}}
			op.GeoM.Translate(float64(ImageSize*x), float64(ImageSize*y))
			screen.DrawImage(mapImages[tile], op)
		}
	}
	return nil
}

func loadMapImageTiles() (map[maptile.MapTile]*ebiten.Image, error) {
	jungle, _, error := ebitenutil.NewImageFromFile("assets/images/jungle.png")
	if error != nil {
		log.Fatal(error)
		return nil, error
	}
	empty := jungle.SubImage(image.Rect(0, 0, ImageSize, ImageSize)).(*ebiten.Image)
	wall := jungle.SubImage(image.Rect(ImageSize, 0, jungle.Bounds().Max.X, ImageSize)).(*ebiten.Image)
	return map[maptile.MapTile]*ebiten.Image{
		maptile.Empty:  empty,
		maptile.Wall:   wall,
		maptile.Player: empty,
	}, nil
}
