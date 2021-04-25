package system

import (
	"log"
	"strconv"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo/common"
	"github.com/project-midgard/midgarts/internal/component"
	"github.com/project-midgard/midgarts/internal/entity"
	"github.com/project-midgard/midgarts/pkg/common/character"
	"github.com/project-midgard/midgarts/pkg/common/character/actionindex"
	"github.com/project-midgard/midgarts/pkg/common/character/statetype"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/grf"
)

type CharacterActionable interface {
	common.BasicFace
	component.CharacterStateComponentFace
	component.CharacterSpriteRenderInfoComponentFace
}

type CharacterActionSystem struct {
	grfFile *grf.File

	characters map[string]*entity.Character
}

func NewCharacterActionSystem(grfFile *grf.File) *CharacterActionSystem {
	return &CharacterActionSystem{
		grfFile,
		map[string]*entity.Character{},
	}
}

func (s *CharacterActionSystem) Add(char *entity.Character) {
	cmp, e := component.NewCharacterAttachmentComponent(s.grfFile, char.Gender, char.JobSpriteID, char.HeadIndex)
	if e != nil {
		log.Fatal(e)
	}

	char.SetCharacterAttachmentComponent(cmp)
	s.characters[strconv.Itoa(int(char.ID()))] = char
}

func (s CharacterActionSystem) AddByInterface(o ecs.Identifier) {
	char := o.(*entity.Character)
	s.Add(char)
}

func (s CharacterActionSystem) Update(dt float32) {
	now := time.Now()

	for _, c := range s.characters {
		previousAnimationHasEnded := now.After(c.AnimationEndsAt)
		var previousAnimationMustStopAtEnd bool

		if c.PreviousState == statetype.Walking {
			previousAnimationMustStopAtEnd = true
		}

		if (c.State != c.PreviousState && c.State != statetype.Idle) ||
			(c.State == statetype.Idle && previousAnimationHasEnded) ||
			(c.State == statetype.Idle && previousAnimationMustStopAtEnd) {
			c.AnimationStartedAt = now

			// TODO: treat special case when attacking
			var forcedDuration time.Duration
			c.ForcedDuration = forcedDuration

			if c.State == statetype.Walking {
				const ConstantMovementSpeed = 2.0
				c.FPSMultiplier = ConstantMovementSpeed
			} else {
				c.FPSMultiplier = 1.0
			}

			c.ActionIndex = actionindex.GetActionIndex(c.State)
			action := c.Files[character.AttachmentBody].ACT.Actions[c.ActionIndex]
			c.AnimationEndsAt = now.Add(time.Duration(action.DurationMilliseconds) * time.Millisecond)
		} else {
			// TODO:
		}
	}
}

func (s CharacterActionSystem) Remove(e ecs.BasicEntity) {
	delete(s.characters, strconv.Itoa(int(e.ID())))
}
