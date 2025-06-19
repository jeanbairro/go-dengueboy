package player

import (
	"app/internal/modules/geom"
	"app/internal/modules/maptile"

	"github.com/hajimehoshi/ebiten/v2"
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

func New(inputHandler *input.Handler) (*Player, error) {
	sprites, err := loadSprites()
	if err != nil {
		return nil, err
	}
	return &Player{
		Position:         geom.Position{X: 1, Y: 1},
		PreviousPosition: geom.Position{X: 1, Y: 1},
		CurrentAction:    NoAction,
		InputHandler:     inputHandler,
		Sprites:          sprites,
		CurrentSprite:    sprites[NoAction],
	}, nil
}

func (p *Player) Update(mapTile *maptile.Map) {
	p.setCurrentAction()
	p.setSprite()
	wantedPosition := p.getWantedPosition()
	if mapTile.GetTileAt(wantedPosition) == maptile.Empty {
		p.PreviousPosition = p.Position
		p.Position = wantedPosition
	}
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

func (p *Player) setSprite() {
	if sprite, exists := p.Sprites[p.CurrentAction]; exists {
		p.CurrentSprite = sprite
		return
	}
	p.CurrentSprite = p.Sprites[NoAction]

}
