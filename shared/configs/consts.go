package configs

const (
	RarityID_Unknown RarityID = iota // Sentinel value for the unknown rarity. This is not a valid rarity ID.
	RarityID_Common
	RarityID_Unusual
	RarityID_Rare
	RarityID_Epic
	RarityID_Legendary
	// RarityID_Mythic
	// RarityID_Ultra
	// RarityID_Super
	RarityID_Unique
	RarityID_SENTINEL // Sentinel value for the end of the rarity list. This is not a valid rarity ID.
)

const (
	PetalID_Unknown PetalID = iota // Sentinel value for the unknown petal. This is not a valid petal ID.
	PetalID_Basic
	PetalID_Fast
	PetalID_Heavy
	PetalID_Rock
	PetalID_Square
	PetalID_Rose
	PetalID_Stinger
	PetalID_Iris
	PetalID_Wing
	PetalID_Missile
	PetalID_Peas
	PetalID_Cactus
	PetalID_Faster
	PetalID_Rose2
	PetalID_Twin
	PetalID_Triplet
	PetalID_Bubble
	PetalID_Pollen
	PetalID_Dandelion
	PetalID_Cactus2
	PetalID_Egg
	PetalID_Egg2
	PetalID_Dahlia
	PetalID_Antennae
	PetalID_Heaviest
	PetalID_Stinger2
	PetalID_YinYang
	PetalID_Web
	PetalID_Web2
	PetalID_Honey
	PetalID_Leaf
	PetalID_Cactus3
	PetalID_Peas2
	PetalID_Basic2
	PetalID_Salt
	PetalID_Rice
	PetalID_Corn
	PetalID_SENTINEL
)

const (
	MobID_Unknown MobID = iota // Sentinel value for the unknown mob. This is not a valid mob ID.
	MobID_Rock
	MobID_Boulder
	MobID_Cactus
	MobID_Ladybug
	MobID_MassiveLadybug
	MobID_Bee
	MobID_BabyAnt
	MobID_WorkerAnt
	MobID_SoldierAnt
	MobID_QueenAnt
	MobID_AntHole
	MobID_Beetle
	MobID_MassiveBeetle
	MobID_Hornet
	MobID_Centipede
	MobID_Centipede2
	MobID_EvilCentipede
	MobID_EvilCentipede2
	MobID_Square
	MobID_Ladybug2
	MobID_LadybugShiny
	MobID_Spider
	MobID_SENTINEL // Sentinel value for the end of the mob list. This is not a valid mob ID.
)

var (
	Rarities map[RarityID]*Rarity = make(map[RarityID]*Rarity)
	Petals   map[PetalID]*Petal   = make(map[PetalID]*Petal)
	Mobs     map[MobID]*Mob       = make(map[MobID]*Mob)
)

var DefaultPetalOrbit PetalOrbit = PetalOrbit{
	Neutral: 65,
	Attack:  130,
	Defend:  40,
}

const (
	AttackMode_Neutral AttackMode = iota
	AttackMode_Attack
	AttackMode_Defend
)

const (
	MobAIMode_Passive MobAIMode = iota
	MobAIMode_FleeWhenHurt
	MobAIMode_FleeWhenPlayerNearby
	MobAIMode_Neutral
	MobAIMode_Aggressive
)

const (
	MobMovementMode_Stationary MobMovementMode = iota
	MobMovementMode_Direct
	MobMovementMode_Sine
	MobMovementMode_Jumps
	MobMovementMode_SlowTurning
	MobMovementMode_Chained
)
