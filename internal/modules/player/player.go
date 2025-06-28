package player

import (
	"app/internal/modules/collision"
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
	initialPositionX float64 = 2
	initialPositionY float64 = 2
	speed            float64 = .08
)

type (
	Player struct {
		Position           geom.Position
		PreviousPosition   geom.Position
		CurrentAction      input.Action
		InputHandler       *input.Handler
		Sprites            map[input.Action][]*ebiten.Image
		CurrentSprite      *ebiten.Image
		CurrentSpriteIndex int
		CollisionSystem    *collision.CollisionSystem
	}
)

func New(inputHandler *input.Handler, collisionSystem *collision.CollisionSystem) (*Player, error) {
	sprites, err := loadSprites()
	if err != nil {
		return nil, err
	}
	return &Player{
		Position:           geom.Position{X: initialPositionX, Y: initialPositionY},
		PreviousPosition:   geom.Position{X: initialPositionX, Y: initialPositionY},
		CurrentAction:      NoAction,
		InputHandler:       inputHandler,
		Sprites:            sprites,
		CurrentSprite:      sprites[NoAction][0],
		CurrentSpriteIndex: 0,
		CollisionSystem:    collisionSystem,
	}, nil
}

func (p *Player) Update(mapTile *maptile.Map) {
	p.setCurrentAction()
	p.setSprite()
	wantedPosition := p.getWantedPosition()
	if !p.CollisionSystem.HasCollision(toObject(wantedPosition)) {
		p.PreviousPosition = p.Position
		p.Position = wantedPosition
	}
}

func toObject(position geom.Position) collision.Object {
	return collision.Object{
		Position: geom.Position{
			X: float64(position.X) * float64(spriteSize),
			Y: float64(position.Y) * float64(spriteSize),
		},
		Width:  float64(spriteSize),
		Height: float64(spriteSize),
	}
}

func (p *Player) getWantedPosition() geom.Position {
	wantedPosition := p.Position
	switch p.CurrentAction {
	case MoveLeft:
		wantedPosition.Move(speed*-1, 0)
	case MoveRight:
		wantedPosition.Move(speed, 0)
	case MoveUp:
		wantedPosition.Move(0, speed*-1)
	case MoveDown:
		wantedPosition.Move(0, speed)
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
	if sprites, exists := p.Sprites[p.CurrentAction]; exists {
		if p.CurrentSpriteIndex >= len(sprites) {
			p.CurrentSpriteIndex = 0
		}
		p.CurrentSprite = sprites[p.CurrentSpriteIndex]
		p.CurrentSpriteIndex++
	}
}
