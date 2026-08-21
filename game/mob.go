package game

import (
	"math/rand/v2"
	"slices"

	"github.com/z46-dev/florr-legacy/shared/configs"
)

func NewMob(g *Game, cfg *configs.Mob) (mob *Mob) {
	mob = &Mob{
		GenericGameEntity: NewEntity(g),
		Config:            cfg,
	}

	mob.Type = EntityTypeMob
	mob.refMob = mob
	mob.Size = cfg.Stats.Radius
	mob.Health.Set(cfg.Stats.Health)
	mob.Damage = cfg.Stats.Damage
	mob.Density = cfg.Stats.Density
	mob.Pushability = cfg.Stats.Pushability

	if cfg.Variance != nil {
		if len(cfg.Variance.Health) == 2 {
			mob.Health.Set(cfg.Stats.Health * (cfg.Variance.Health[0] + rand.Float64()*(cfg.Variance.Health[1]-cfg.Variance.Health[0])))
		}

		if len(cfg.Variance.Damage) == 2 {
			mob.Damage = cfg.Stats.Damage * (cfg.Variance.Damage[0] + rand.Float64()*(cfg.Variance.Damage[1]-cfg.Variance.Damage[0]))
		}

		if len(cfg.Variance.Radius) == 2 {
			mob.Size = cfg.Stats.Radius * (cfg.Variance.Radius[0] + rand.Float64()*(cfg.Variance.Radius[1]-cfg.Variance.Radius[0]))
		}
	}

	if cfg.AI != nil {
		mob.AI = NewMobAI(mob)
	}

	mob.Facing = rand.Float64() * TAU
	g.mobs.Add(mob)
	return
}

func (mob *Mob) Update(deltaMilliseconds, damping, positionFactor float64) {
	if !mob.Health.IsAlive() {
		mob.Die()
		return
	}

	mob.GenericGameEntity.Update(damping, positionFactor)
	// if mob.Velocity.X != 0 || mob.Velocity.Y != 0 {
	// 	mob.Facing = math.Atan2(mob.Velocity.Y, mob.Velocity.X)
	// }

	if mob.Position.X < -mob.Game.Width*0.5+mob.Size {
		mob.Position.X = -mob.Game.Width*0.5 + mob.Size
		mob.Velocity.X = 0
	} else if mob.Position.X > mob.Game.Width*0.5-mob.Size {
		mob.Position.X = mob.Game.Width*0.5 - mob.Size
		mob.Velocity.X = 0
	}

	if mob.Position.Y < -mob.Game.Height*0.5+mob.Size {
		mob.Position.Y = -mob.Game.Height*0.5 + mob.Size
		mob.Velocity.Y = 0
	} else if mob.Position.Y > mob.Game.Height*0.5-mob.Size {
		mob.Position.Y = mob.Game.Height*0.5 - mob.Size
		mob.Velocity.Y = 0
	}

	mob.Game.mobsSpatialHash.Insert(mob)
}

func (mob *Mob) Die() {
	if mob.isDead {
		return
	}

	mob.isDead = true
	mob.Game.mobs.Remove(mob.ID)
	mob.Game.mobsCache.Remove(mob.ID)
}

func (mob *Mob) Collide() {
	for _, other := range mob.Game.mobsSpatialHash.Retrieve(mob.AABB) {
		if other.ID == mob.ID || slices.Contains(mob.collisionIDs, other.ID) || !circleIntersects(mob.Position, mob.Size, other.Position, other.Size) {
			continue
		}

		CollideEntities(mob.GenericGameEntity, other.GenericGameEntity)
	}
}
