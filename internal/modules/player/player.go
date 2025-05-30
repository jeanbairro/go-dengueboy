package player

import "app/internal/modules/geom"

const (
	MoveLeft  CurrentAction = "move_left"
	MoveRight CurrentAction = "move_right"
	MoveUp    CurrentAction = "move_up"
	MoveDown  CurrentAction = "move_down"
	Interact  CurrentAction = "interact"
)

const (
	MoveValue int = 1
)

type (
	CurrentAction string

	Player struct {
		Position      geom.Position
		CurrentAction CurrentAction
	}
)

func (p *Player) Update(key string) {
	p.setCurrentAction(key)

	switch p.CurrentAction {
	case MoveLeft:
		p.Position.Move(MoveValue*-1, 0)
	case MoveRight:
		p.Position.Move(MoveValue, 0)
	case MoveUp:
		p.Position.Move(0, MoveValue)
	case MoveDown:
		p.Position.Move(0, MoveValue*-1)
	case Interact:
		//
	}
}

func (p *Player) setCurrentAction(key string) {
	p.CurrentAction = MoveDown
}
