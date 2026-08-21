package game

import (
	"sync"
	"time"

	"github.com/z46-dev/arctic"
	"github.com/z46-dev/florr-legacy/shared"
	"github.com/z46-dev/florr-legacy/shared/configs"
	"github.com/z46-dev/florr-legacy/shared/util"
	"github.com/z46-dev/golog"
)

type (
	GameMode   uint8
	EntityType uint8

	GameOptions struct {
		Mode                   GameMode
		Seed                   uint64
		NumTeams               int
		PopulatePreservedWorld bool
	}

	Game struct {
		Time                       uint64
		Options                    GameOptions
		NumTeams                   int
		clients                    *shared.Collection[*Client]
		flowers                    *shared.Collection[*Flower]
		petals                     *shared.Collection[*Petal]
		mobs                       *shared.Collection[*Mob]
		flowersSpatialHash         *util.SpatialHash[*Flower]
		petalsSpatialHash          *util.SpatialHash[*Petal]
		mobsSpatialHash            *util.SpatialHash[*Mob]
		flowersCache               *shared.Collection[*FlowerCache]
		petalsCache                *shared.Collection[*PetalCache]
		mobsCache                  *shared.Collection[*MobCache]
		Width, Height              float64
		nextEntityID, nextClientID uint64
		stopChan, doneChan         chan struct{}
		runOnce, stopOnce          sync.Once
		lastUpdate                 time.Time
		L                          *golog.Logger
		mspts                      []float64
		mspt                       float64
		fps, framesNow             uint64
	}

	Camera struct {
		Client                            *Client
		Position                          *util.Vec2
		FOV                               float64
		FlowersSeen, PetalsSeen, MobsSeen map[uint64]uint64
	}

	Client struct {
		ID     uint64
		Peer   *arctic.ServerClient
		Game   *Game
		Camera *Camera
		Body   *Flower
		TeamID int
	}

	HealthComponent struct {
		Current, Maximum float64
		LastDamageTick   uint64
		Regeneration     float64
	}

	GenericGameEntity struct {
		Game                 *Game
		Type                 EntityType
		ID                   uint64
		TeamID               uint64
		Position, Velocity   *util.Vec2
		Size, Facing         float64
		Density, Pushability float64
		Health               *HealthComponent
		Damage               float64
		AABB                 *util.AABB
		isDead               bool
		collisionIDs         []uint64
		refFlower            *Flower
		refPetal             *Petal
		refMob               *Mob
	}

	LivePetal struct {
		Slot     *Slot
		Timer    float64
		Petal    *Petal
		SubIndex int
	}

	Slot struct {
		Flower               *Flower
		PetalConfig          *configs.Petal
		Petals               []*LivePetal
		SlotIndex, NumClumps int
	}

	Flower struct {
		*GenericGameEntity
		Name                              string
		NumSlots                          int
		Slots                             []*Slot
		OrbitAngle, OrbitRadiansPerSecond float64
		AttackMode                        configs.AttackMode
		Movement                          *util.Vec2
	}

	Petal struct {
		*GenericGameEntity
		Parent     *Flower
		Config     *configs.Petal
		HealTicker float64
	}

	Mob struct {
		*GenericGameEntity
		Summoner *Flower
		Config   *configs.Mob
		AI       *MobAI
	}

	MobAI struct {
		Mob                         *Mob
		Cfg                         *configs.MobAI
		Target                      *GenericGameEntity
		TargetTimer, Timer1, Timer2 float64
		MayRetarget                 bool
	}

	// Comms types
	GenericObjectCache struct {
		ID, AsOf                                                uint64
		New, Old                                                *shared.Writer
		X, Y, Size, Rotation                                    float64
		StatusFlags                                             uint8
		PosChanged, SizeChanged, RotChanged, StatusFlagsChanged bool
	}

	FlowerCache struct {
		GenericObjectCache
		HealthRatio                      float64
		AttackMode                       configs.AttackMode
		HealthChanged, AttackModeChanged bool
	}

	PetalCache struct {
		GenericObjectCache
		HealthRatio float64
	}

	MobCache struct {
		GenericObjectCache
		HealthRatio   float64
		HealthChanged bool
	}
)

const (
	GameModeOldFFA GameMode = iota
	GameModeOldBattleRoyale
	GameModeOldWaves
	GameModeLineMap
	GameModeWaves
	GameModeWorldMap
)

const (
	EntityTypeGeneric EntityType = iota
	EntityTypeFlower
	EntityTypePetal
	EntityTypeMob
)
