package configs

import "github.com/z46-dev/gctx2d"

type (
	RarityID        uint8
	PetalID         uint8
	MobID           uint8
	AttackMode      uint8
	MobAIMode       uint8
	MobMovementMode uint8

	Rarity struct {
		ID    RarityID     `json:"id"`    // The ID of the rarity
		Name  string       `json:"name"`  // The name of the rarity
		Color gctx2d.Color `json:"color"` // The color of the rarity
	}

	PetalMeta struct {
		Name        string   `json:"name"`                  // The name of the petal
		Description string   `json:"description,omitempty"` // The description of the petal
		Rarity      RarityID `json:"rarity"`                // The corresponding rarity of the petal
	}

	PetalOrbit struct {
		Neutral float64 `json:"neutral"`
		Attack  float64 `json:"attack"`
		Defend  float64 `json:"defend"`
	}

	PetalStats struct {
		ReloadMilliseconds uint32      `json:"reload_milliseconds"`
		Health             float64     `json:"health"`
		Damage             float64     `json:"damage"`
		Radius             float64     `json:"radius"`
		Orbit              *PetalOrbit `json:"orbit,omitempty"` // The orbit distances for this petal
	}

	PetalHeal struct {
		Amount float64 `json:"amount"` // The amount of health this petal heals
		Delay  float64 `json:"delay"`  // The delay in milliseconds between heals
	}

	PetalSpecial struct {
		Count                    int        `json:"count,omitempty"`                       // How many petals are spawned for this petal
		ClumpSize                int        `json:"clump_size,omitempty"`                  // If there are multiple petals, how do we divide clumps up? 1 = individual, 2 = pairs, and so forth.
		Passive                  bool       `json:"passive,omitempty"`                     // Is this petal unable to be in "attack mode"?
		Wearable                 bool       `json:"wearable,omitempty"`                    // Is this petal wearable by the player (aka not an entity)?
		Heal                     *PetalHeal `json:"heal,omitempty"`                        // Can this petal heal the player?
		RotationRadiansPerSecond float64    `json:"rotation_radians_per_second,omitempty"` // How much extra rotation does this petal gift to the flower?
		ConstantHealPerSecond    float64    `json:"constant_heal_per_second,omitempty"`    // How much constant healing does this petal provide to the flower?
	}

	Petal struct {
		ID      PetalID       `json:"id"`
		Meta    *PetalMeta    `json:"meta,omitempty"`
		Stats   *PetalStats   `json:"stats,omitempty"`
		Special *PetalSpecial `json:"special,omitempty"`
	}

	MobMeta struct {
		Name        string   `json:"name"`                  // The name of the mob
		Description string   `json:"description,omitempty"` // The description of the mob
		Rarity      RarityID `json:"rarity"`                // The corresponding rarity of the mob
	}

	MobStats struct {
		Health      float64 `json:"health"`
		Damage      float64 `json:"damage"`
		Radius      float64 `json:"radius"`
		Pushability float64 `json:"pushability"`
		Density     float64 `json:"density"`
	}

	MobSpecial struct{}

	MobDrop struct {
		PetalID PetalID `json:"petal_id"` // The ID of the petal that this mob drops
		Chance  float64 `json:"chance"`   // The chance that this mob will drop the petal (0.0 - 1.0)
	}

	// MobVariance{ Health: [2]float64{.8,1.2} } would mean that the mob's health is multiplied by a random number between 0.8 and 1.2.
	MobVariance struct {
		Health []float64 `json:"health"` // The min and max health multiplier of the mob
		Damage []float64 `json:"damage"` // The min and max damage multiplier of the mob
		Radius []float64 `json:"radius"` // The min and max radius multiplier of the mob
	}

	MobAI struct {
		Mode         MobAIMode       `json:"mode"`          // The AI mode of the mob
		MovementMode MobMovementMode `json:"movement_mode"` // The movement mode of the mob
		IdleSpeed    float64         `json:"idle_speed"`    // The speed of the mob when idle
		AltSpeed     float64         `json:"alt_speed"`     // The speed of the mob when in an alternate state (e.g. attacking, fleeing)
	}

	Mob struct {
		ID       MobID        `json:"id"`
		Meta     *MobMeta     `json:"meta,omitempty"`
		Stats    *MobStats    `json:"stats,omitempty"`
		Special  *MobSpecial  `json:"special,omitempty"`
		Variance *MobVariance `json:"variance,omitempty"`
		AI       *MobAI       `json:"ai,omitempty"`
		Drops    []*MobDrop   `json:"drops,omitempty"`
	}
)
