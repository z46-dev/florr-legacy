package game

import (
	"math"
	"math/rand/v2"
	"slices"

	"github.com/z46-dev/florr-legacy/shared/configs"
	"github.com/z46-dev/florr-legacy/shared/util"
)

func NewMobAI(mob *Mob) (ai *MobAI) {
	ai = &MobAI{
		Mob: mob,
		Cfg: mob.Config.AI,
	}

	return
}

func (ai *MobAI) ValidateTarget(target *GenericGameEntity, fov float64) (valid bool) {
	valid = target != nil && !target.isDead && target.TeamID != ai.Mob.TeamID

	if fov > 0 {
		valid = valid && ai.Mob.Position.DistSquared(target.Position) <= fov*fov
	}

	return
}

func (ai *MobAI) FindTarget(fov float64) {
	var (
		possibleTargets []*GenericGameEntity
		aabb            *util.AABB = &util.AABB{
			X1: ai.Mob.Position.X - fov,
			Y1: ai.Mob.Position.Y - fov,
			X2: ai.Mob.Position.X + fov,
			Y2: ai.Mob.Position.Y + fov,
		}
	)

	for _, flower := range ai.Mob.Game.flowersSpatialHash.Retrieve(aabb) {
		if ai.ValidateTarget(flower.GenericGameEntity, 0) {
			possibleTargets = append(possibleTargets, flower.GenericGameEntity)
		}
	}

	for _, mob := range ai.Mob.Game.mobsSpatialHash.Retrieve(aabb) {
		if ai.ValidateTarget(mob.GenericGameEntity, 0) {
			possibleTargets = append(possibleTargets, mob.GenericGameEntity)
		}
	}

	if len(possibleTargets) == 0 {
		ai.Target = nil
		return
	}

	slices.SortFunc(possibleTargets, func(a, b *GenericGameEntity) (cmp int) {
		var distA, distB float64 = a.Position.DistSquared(ai.Mob.Position), b.Position.DistSquared(ai.Mob.Position)

		switch {
		case distA < distB:
			cmp = -1
		case distA > distB:
			cmp = 1
		default:
			cmp = 0
		}

		return
	})

	ai.Target = possibleTargets[0]
}

func (ai *MobAI) thinkIdle(deltaMilliseconds float64) {
	if ai.Cfg.IdleSpeed <= 0 {
		return
	}

	// Move in facing direction until our "energy" runs out, then wait, pick new direction, repeat
	if ai.Timer1 > 0 {
		ai.Timer1 -= deltaMilliseconds
		if ai.Timer1 <= 0 {
			ai.Timer2 = rand.Float64()*3000 + 1500
		}

		ai.Mob.Velocity.X += math.Cos(ai.Mob.Facing) * ai.Cfg.IdleSpeed * SPEED_SCALE_FACTOR
		ai.Mob.Velocity.Y += math.Sin(ai.Mob.Facing) * ai.Cfg.IdleSpeed * SPEED_SCALE_FACTOR

		if rand.Float64() < 0.0001 {
			ai.Mob.Facing += rand.Float64()*math.Pi*.2 - math.Pi*0.1
		}
	} else if ai.Timer2 > 0 {
		ai.Timer2 -= deltaMilliseconds
		if ai.Timer2 <= 0 {
			ai.Timer1 = rand.Float64()*1500 + 500
			ai.Mob.Facing += rand.Float64()*math.Pi - math.Pi*0.5
		}
	} else {
		ai.Timer1 = rand.Float64()*1500 + 500
	}
}

func (ai *MobAI) thinkActive(_ float64) {
	var angle float64 = ai.Mob.Position.AngleTo(ai.Target.Position)

	ai.Mob.Velocity.X += math.Cos(angle) * ai.Cfg.AltSpeed * SPEED_SCALE_FACTOR
	ai.Mob.Velocity.Y += math.Sin(angle) * ai.Cfg.AltSpeed * SPEED_SCALE_FACTOR
	ai.Mob.Facing = angle
}

func (ai *MobAI) Think(deltaMilliseconds float64) {
	if ai.Cfg.Mode == configs.MobAIMode_Aggressive {
		var fov float64 = ai.Mob.Size*1.25 + 800
		if ai.Target != nil {
			if !ai.ValidateTarget(ai.Target, fov) {
				ai.Target = nil
			}

			if ai.MayRetarget {
				ai.TargetTimer += deltaMilliseconds
			} else {
				ai.TargetTimer = 0
			}
		} else {
			ai.TargetTimer += deltaMilliseconds
		}

		if ai.TargetTimer >= AI_RETARGET_MS {
			ai.TargetTimer = 0
			ai.FindTarget(fov)
		}
	}

	if ai.Target != nil {
		ai.thinkActive(deltaMilliseconds)
	} else {
		ai.thinkIdle(deltaMilliseconds)
	}
}

func (ai *MobAI) OnHostileCollision(other *GenericGameEntity) {
	if ai.Target != nil || ai.Cfg.Mode != configs.MobAIMode_Neutral {
		return
	}

	if other.refPetal != nil {
		other = other.refPetal.Parent.GenericGameEntity
	}

	ai.Target = other
}
