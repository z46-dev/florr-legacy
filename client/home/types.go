package home

import (
	"github.com/z46-dev/florr-legacy/client/render"
	"github.com/z46-dev/florr-legacy/shared/configs"
	"github.com/z46-dev/spiffy"
)

type (
	MenuKind uint8

	ChangelogEntry struct {
		Date    string
		Changes []string
	}

	MenuState struct {
		Open                   MenuKind
		MouseX, MouseY, Scroll float32
		Width, Height          float32
		ScrollbarGrabOffset    float32
		DraggingScrollbar      bool
		Petals                 []configs.PetalID
		Mobs                   []configs.MobID
		Changelog              []ChangelogEntry
	}

	PetalAnimationPetal struct {
		ID         configs.PetalID
		X, Y, S, R float32
		Time, Seed float32
	}

	PetalAnimation struct {
		Ready          bool
		PossiblePetals []configs.PetalID
		Petals         []*PetalAnimationPetal
		Width, Height  float32
		time           float64
	}

	Home struct {
		RenderState *render.RendererState
		UI          *spiffy.UserInterface
		title       *spiffy.Label
		input       struct {
			prompt, submitHint, errorLabel *spiffy.Label
			nickname                       *spiffy.TextInput
		}

		tips struct {
			howTo, move, attack, defend *spiffy.Label
		}

		menus struct {
			settings, petals, mobs, changelog *spiffy.Button
			reddit, discord                   *spiffy.Button
			keyboardMovement, movementHelper  *spiffy.Checkbox
		}

		menu           MenuState
		petalAnimation *PetalAnimation
		OnSubmit       func(string)

		// modalBlocker                       *spiffy.Button
		// modalPanel                         *spiffy.Panel
		// modalTitle                         *spiffy.Label
		// modalLines                         []*spiffy.Label
		// modalClose                         *spiffy.Button
		// backgroundTime                     float64
	}
)
