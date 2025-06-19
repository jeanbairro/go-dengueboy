package player

import (
	"app/internal/modules/geom"
	"app/internal/modules/maptile"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	input "github.com/quasilyte/ebitengine-input"
)

const (
	NoAction input.Action = iota
	MoveLeft
	MoveRight
	MoveUp
	MoveDown
	Interact
)

const (
	MoveValue int = 1
)

type (
	Player struct {
		Position         geom.Position
		PreviousPosition geom.Position
		CurrentAction    input.Action
		InputHandler     *input.Handler
		Sprites          map[input.Action]*ebiten.Image
		CurrentSprite    *ebiten.Image
	}
)

func New(inputHandler *input.Handler) *Player {
	sprites := loadPlayerSprites()
	return &Player{
		Position:         geom.Position{X: 1, Y: 1},
		PreviousPosition: geom.Position{X: 1, Y: 1},
		CurrentAction:    NoAction,
		InputHandler:     inputHandler,
		Sprites:          sprites,
		CurrentSprite:    sprites[NoAction],
	}
}

func (p *Player) Update(mapTile *maptile.Map) {
	p.setCurrentAction()
	wantedPosition := p.getWantedPosition()
	if mapTile.GetTileAt(wantedPosition) == maptile.Empty {
		p.PreviousPosition = p.Position
		p.Position = wantedPosition
	}
	p.updateSprite()
}

func (p *Player) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(p.Position.X*32), float64(p.Position.Y*32))
	screen.DrawImage(p.CurrentSprite, op)
}

func (p *Player) getWantedPosition() geom.Position {
	wantedPosition := p.Position
	switch p.CurrentAction {
	case MoveLeft:
		wantedPosition.Move(MoveValue*-1, 0)
	case MoveRight:
		wantedPosition.Move(MoveValue, 0)
	case MoveUp:
		wantedPosition.Move(0, MoveValue*-1)
	case MoveDown:
		wantedPosition.Move(0, MoveValue)
	}
	return wantedPosition
}

func (p *Player) setCurrentAction() {
	if p.InputHandler.ActionIsPressed(MoveLeft) {
		p.CurrentAction = MoveLeft
		return
	}
	if p.InputHandler.ActionIsPressed(MoveRight) {
		p.CurrentAction = MoveRight
		return
	}
	if p.InputHandler.ActionIsPressed(MoveUp) {
		p.CurrentAction = MoveUp
		return
	}
	if p.InputHandler.ActionIsPressed(MoveDown) {
		p.CurrentAction = MoveDown
		return
	}
	if p.InputHandler.ActionIsPressed(Interact) {
		p.CurrentAction = Interact
		return
	}
	p.CurrentAction = NoAction
}

func (p *Player) updateSprite() {
	if sprite, exists := p.Sprites[p.CurrentAction]; exists {
		p.CurrentSprite = sprite
	} else {
		p.CurrentSprite = p.Sprites[NoAction]
	}
}

func loadPlayerSprites() map[input.Action]*ebiten.Image {
	charImage, _, err := ebitenutil.NewImageFromFile("assets/images/char.png")
	if err != nil {
		log.Fatal(err)
	}

	spriteSize := 32
	return map[input.Action]*ebiten.Image{
		NoAction:  charImage.SubImage(image.Rect(0, 0, spriteSize, spriteSize)).(*ebiten.Image),
		MoveLeft:  charImage.SubImage(image.Rect(spriteSize, 0, spriteSize*2, spriteSize)).(*ebiten.Image),
		MoveRight: charImage.SubImage(image.Rect(spriteSize*2, 0, spriteSize*3, spriteSize)).(*ebiten.Image),
		MoveUp:    charImage.SubImage(image.Rect(spriteSize*3, 0, spriteSize*4, spriteSize)).(*ebiten.Image),
		MoveDown:  charImage.SubImage(image.Rect(spriteSize*4, 0, spriteSize*5, spriteSize)).(*ebiten.Image),
	}
}
