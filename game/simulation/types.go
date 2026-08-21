package simulation

import (
	"time"

	"github.com/z46-dev/arctic"
	"github.com/z46-dev/florr-legacy/client/render"
	"github.com/z46-dev/florr-legacy/shared"
	"github.com/z46-dev/florr-legacy/shared/configs"
	"github.com/z46-dev/golog"
)

type (
	Simulation struct {
		RenderState      *render.RendererState
		Width, Height    float32
		Time             float32
		WorldLerpAlpha   float32
		Camera           *Camera
		PlayerID         uint64
		Player           *Flower
		Flowers          *shared.Collection[*Flower]
		Petals           *shared.Collection[*Petal]
		Mobs             *shared.Collection[*Mob]
		Client           *arctic.Client
		WorldWidth       float32
		WorldHeight      float32
		LastWorldStateAt time.Time
		L                *golog.Logger
		serverStats      struct {
			numFlowers int
			numPetals  int
			numMobs    int
			mspt       float64
			fps        uint64
		}
	}

	Camera struct {
		X, Y, PrevX, PrevY, RealX, RealY              float32
		FieldOfView, PrevFieldOfView, RealFieldOfView float32
	}

	GenericWorldObject struct {
		ID                                   uint64
		Simulation                           *Simulation
		X, Y, PrevX, PrevY, RealX, RealY     float32
		Size, PrevSize, RealSize             float32
		Rotation, PrevRotation, RealRotation float32
		Damaged                              bool
		DamagedFlashCurrent                  float32
		DamagedFlashWasActive                bool
		IsDead                               bool
		DeathDelta                           float32
		Alpha                                float32
	}

	HasHealth struct {
		GenericWorldObject
		HealthRatio, RealHealthRatio float32
		SubHealthRatio               float32
		LastDamagedAt                float32
	}

	Flower struct {
		HasHealth
		Name       string
		AttackMode configs.AttackMode
		MouthCurve float32
		EyeClip    float32
	}

	Petal struct {
		GenericWorldObject
		PetalID  configs.PetalID
		RarityID configs.RarityID
	}

	Mob struct {
		HasHealth
		MobID     configs.MobID
		RarityID  configs.RarityID
		Attacking bool
	}
)

const (
	DEATH_DELTA_DURATION  float32 = 0.2 // seconds
	HEALTH_DELTA_DURATION float32 = 0.3 // seconds
)
