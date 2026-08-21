package game

import (
	"math"
	"time"
)

const (
	TicksPerSecond uint64        = 20
	TickInterval   time.Duration = time.Second / time.Duration(TicksPerSecond)
)

// Physics Constants
const (
	TAU                       float64 = 2 * math.Pi
	velocityIntegral          float64 = 125
	velocityRetention         float64 = 0.992
	petalSpring               float64 = 0.007
	petalDistanceFactor       float64 = 0.0011111111111111111
	baseFlowerSize            float64 = 25
	baseFlowerHealth          float64 = 100
	baseFlowerDamage          float64 = 25
	baseFlowerRegeneration    float64 = 1.0 / 240000.0
	baseOrbitRadius           float64 = 65
	baseOrbitRadiansPerSecond float64 = 2.5
	SPEED_SCALE_FACTOR        float64 = 0.125
	AI_RETARGET_MS            float64 = 3500
)

const (
	WorldWidth  float64 = 50000
	WorldHeight float64 = 3000
)
