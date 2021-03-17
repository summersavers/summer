package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type BaddieComponent struct {
	Name         string
	Phases       map[string]Phase
	CurrentPhase string
	Font         *common.Font
	Clip         *common.Player
	HP, MaxHP    int
	CastAt       int
	CastTime     int
	Spritesheet  *common.Spritesheet
}

type Phase struct {
	Animations []*common.Animation
	Abilities  []BaddieAbility
}

type BaddieAbility struct {
	Weight            int
	LogMessages       []string
	Times             []float32
	AbilityAnimations []int
}

func (c *BaddieComponent) GetBaddieComponent() *BaddieComponent {
	return c
}

type BaddieFace interface {
	GetBaddieComponent() *BaddieComponent
}

type BaddiePhaseAble interface {
	BaddieFace
	common.BasicFace
}

type NotBaddiePhaseComponent struct{}

func (n *NotBaddiePhaseComponent) GetNotBaddiePhaseComponent() *NotBaddiePhaseComponent {
	return n
}

type NotBaddiePhaseAble interface {
	GetNotBaddiePhaseComponent() *NotBaddiePhaseComponent
}

type baddiePhaseEntity struct {
	*ecs.BasicEntity
	*BaddieComponent
}

type BaddiePhaseSystem struct {
	entities      []baddiePhaseEntity
	baddiesprites []*sprite

	w *ecs.World
}

func (s *BaddiePhaseSystem) New(w *ecs.World) {
	s.w = w
}

func (s *BaddiePhaseSystem) reposition() {
	var positions []engo.Point
	switch len(s.entities) {
	case 1:
		positions = append(positions, engo.Point{X: 320, Y: 130})
	case 2:
		positions = append(positions, engo.Point{X: 320, Y: 130})
		positions = append(positions, engo.Point{X: 320, Y: 130})
	case 3:
		positions = append(positions, engo.Point{X: 320, Y: 130})
		positions = append(positions, engo.Point{X: 320, Y: 130})
		positions = append(positions, engo.Point{X: 320, Y: 130})
	case 4:
		positions = append(positions, engo.Point{X: 320, Y: 130})
		positions = append(positions, engo.Point{X: 320, Y: 130})
		positions = append(positions, engo.Point{X: 320, Y: 130})
		positions = append(positions, engo.Point{X: 320, Y: 130})
	}
	for k, position := range positions {
		s.baddiesprites[k].SetCenter(position)
	}
}

func (s *BaddiePhaseSystem) Add(basic *ecs.BasicEntity, bp *BaddieComponent) {
	s.entities = append(s.entities, baddiePhaseEntity{basic, bp})
	spr := sprite{BasicEntity: ecs.NewBasic()}
	spr.Drawable = bp.Spritesheet.Drawable(0)
	spr.Width = spr.Drawable.Width()
	spr.Height = spr.Drawable.Height()
	spr.SetZIndex(2)
	s.w.AddEntity(&spr)
	for _, sys := range s.w.Systems() {
		switch system := sys.(type) {
		case *TargetSystem:
			system.Add(&spr.BasicEntity, &spr.RenderComponent, bp)
		}
	}
	s.baddiesprites = append(s.baddiesprites, &spr)
	s.reposition()
}

func (s *BaddiePhaseSystem) AddByInterface(i ecs.Identifier) {
	a := i.(BaddiePhaseAble)
	s.Add(a.GetBasicEntity(), a.GetBaddieComponent())
}

func (s *BaddiePhaseSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, entity := range s.entities {
		if entity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		s.entities = append(s.entities[:delete], s.entities[delete+1:]...)
		s.baddiesprites = append(s.baddiesprites[:delete], s.baddiesprites[delete+1:]...)
	}
	s.reposition()
}

func (s *BaddiePhaseSystem) Update(dt float32) {}
