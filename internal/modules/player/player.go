package player

import (
	"app/internal/modules/collision"
	"app/internal/modules/entity"
	"app/internal/modules/geom"

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
	initialPositionX float64 = 64
	initialPositionY float64 = 64
	speed            float64 = 2.56
)

var _ entity.Entity = (*Player)(nil)

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

func (p *Player) Update() {
	movement := p.readMovementVector()
	p.Move(movement)
	p.setSprite()
}

// Move applies movement to p.Position one axis at a time (X then Y),
// checking p.CollisionSystem.HasCollision independently per axis so the
// player slides along a wall instead of stopping dead when only one axis
// is blocked. It always updates p.PreviousPosition to the position before
// the move and p.CurrentAction to reflect the direction implied by
// movement, regardless of whether the move was blocked.
func (p *Player) Move(movement geom.Position) {
	p.updateFacing(movement)
	p.PreviousPosition = p.Position

	wantedX := p.Position.X + movement.X
	if !p.CollisionSystem.HasCollision(toObject(geom.Position{X: wantedX, Y: p.Position.Y})) {
		p.Position.X = wantedX
	}

	wantedY := p.Position.Y + movement.Y
	if !p.CollisionSystem.HasCollision(toObject(geom.Position{X: p.Position.X, Y: wantedY})) {
		p.Position.Y = wantedY
	}
}

func toObject(position geom.Position) collision.Object {
	return collision.Object{
		Position: position,
		Width:    float64(spriteSize),
		Height:   float64(spriteSize),
	}
}

// movementIntent turns the 4 pressed-direction booleans into raw per-axis
// intent in {-1, 0, 1}, summing opposing directions to 0 when both are
// pressed.
func movementIntent(left, right, up, down bool) (dx, dy float64) {
	if left {
		dx--
	}
	if right {
		dx++
	}
	if up {
		dy--
	}
	if down {
		dy++
	}
	return dx, dy
}

func (p *Player) readMovementVector() geom.Position {
	dx, dy := movementIntent(
		p.InputHandler.ActionIsPressed(MoveLeft),
		p.InputHandler.ActionIsPressed(MoveRight),
		p.InputHandler.ActionIsPressed(MoveUp),
		p.InputHandler.ActionIsPressed(MoveDown),
	)
	return geom.NormalizeDiagonal(dx, dy, speed)
}

// updateFacing sets p.CurrentAction from the sign of movement. MoveLeft and
// MoveRight take priority over MoveUp and MoveDown when both axes are
// non-zero, since the spritesheet only has 4 cardinal rows.
func (p *Player) updateFacing(movement geom.Position) {
	switch {
	case movement.X < 0:
		p.CurrentAction = MoveLeft
	case movement.X > 0:
		p.CurrentAction = MoveRight
	case movement.Y < 0:
		p.CurrentAction = MoveUp
	case movement.Y > 0:
		p.CurrentAction = MoveDown
	default:
		p.CurrentAction = NoAction
	}
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
