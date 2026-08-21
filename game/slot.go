package game

import (
	"math"

	"github.com/z46-dev/florr-legacy/shared/configs"
	"github.com/z46-dev/florr-legacy/shared/util"
)

func NewSlot(g *Game, f *Flower, idx int) (slot *Slot) {
	slot = &Slot{
		Flower:    f,
		SlotIndex: idx,
	}

	return
}

func (s *Slot) Define(cfg *configs.Petal) (self *Slot) {
	// Clean up after previous petal config
	if s.PetalConfig != nil {
		s.Flower.OrbitRadiansPerSecond -= s.PetalConfig.Special.RotationRadiansPerSecond
	}

	s.PetalConfig = cfg
	s.Petals = make([]*LivePetal, cfg.Special.Count)
	s.NumClumps = cfg.Special.Count / cfg.Special.ClumpSize
	if cfg.Special.Count%cfg.Special.ClumpSize != 0 {
		s.NumClumps++
	}

	for i := range s.Petals {
		s.Petals[i] = &LivePetal{
			Slot:     s,
			Timer:    -float64(i * int(TickInterval.Milliseconds())),
			Petal:    nil,
			SubIndex: i,
		}
	}

	s.Flower.OrbitRadiansPerSecond += cfg.Special.RotationRadiansPerSecond

	self = s
	return
}

func (s *Slot) Update(deltaMilliseconds float64, numClumps, currentClump int) {
	var (
		subClumpIdx int        = 0
		target      *util.Vec2 = util.Vec2_0()
	)

	for i, livePetal := range s.Petals {
		if i > 0 && i%s.PetalConfig.Special.ClumpSize == 0 {
			subClumpIdx++
		}

		if livePetal.Petal == nil {
			livePetal.Timer += deltaMilliseconds
			if livePetal.Timer >= float64(s.PetalConfig.Stats.ReloadMilliseconds) {
				livePetal.Petal = NewPetal(s.Flower.Game, s.PetalConfig, s.Flower)
				livePetal.Petal.Position.X, livePetal.Petal.Position.Y = s.Flower.Position.X, s.Flower.Position.Y
				livePetal.Timer = 0
			}

			continue
		}

		if !livePetal.Petal.Health.IsAlive() {
			livePetal.Petal = nil
			livePetal.Timer = 0
			continue
		}

		var orbitRadius float64
		switch s.Flower.AttackMode {
		case configs.AttackMode_Neutral:
			orbitRadius = s.PetalConfig.Stats.Orbit.Neutral
		case configs.AttackMode_Attack:
			orbitRadius = s.PetalConfig.Stats.Orbit.Attack
		case configs.AttackMode_Defend:
			orbitRadius = s.PetalConfig.Stats.Orbit.Defend
		}

		target.X, target.Y = s.Flower.Position.X, s.Flower.Position.Y

		var wantsHeal bool = s.PetalConfig.Special.Heal != nil && livePetal.Petal.HealTicker >= s.PetalConfig.Special.Heal.Delay && s.Flower.Health.Current < s.Flower.Health.Maximum

		if wantsHeal {
			//
		} else {
			target.DirectionMagnitudeInto(s.Flower.OrbitAngle+(float64(currentClump+subClumpIdx)/float64(numClumps))*TAU, orbitRadius)
			if s.PetalConfig.Special.ClumpSize > 1 {
				target.DirectionMagnitudeInto((float64(i%s.PetalConfig.Special.ClumpSize)/float64(s.PetalConfig.Special.ClumpSize))*TAU, s.PetalConfig.Stats.Radius*1.5)
			}
		}

		var (
			deltaX, deltaY  float64 = target.X - livePetal.Petal.Position.X, target.Y - livePetal.Petal.Position.Y
			distanceSquared float64 = deltaX*deltaX + deltaY*deltaY
		)

		var collisionSize float64 = s.Flower.Size/2 + livePetal.Petal.Size/2
		if wantsHeal && distanceSquared < collisionSize*collisionSize {
			s.Flower.Health.Heal(s.PetalConfig.Special.Heal.Amount)
			livePetal.Petal.Health.Current = -1
		}

		if distanceSquared > 0 {
			var (
				distance      float64 = math.Sqrt(distanceSquared)
				distanceScale float64 = min(2, max(0.01, distanceSquared*petalDistanceFactor))
				acceleration  float64 = deltaMilliseconds * petalSpring * distanceScale
			)

			livePetal.Petal.Velocity.X += deltaX / distance * acceleration
			livePetal.Petal.Velocity.Y += deltaY / distance * acceleration
		}
	}
}
