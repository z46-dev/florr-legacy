package game

import (
	"slices"

	"github.com/z46-dev/florr-legacy/shared/configs"
)

func NewPetal(g *Game, cfg *configs.Petal, parent *Flower) (petal *Petal) {
	petal = &Petal{
		GenericGameEntity: NewEntity(g),
		Config:            cfg,
		Parent:            parent,
	}

	petal.Type = EntityTypePetal
	petal.refPetal = petal
	petal.Size = cfg.Stats.Radius
	petal.Health.Set(cfg.Stats.Health)
	petal.Damage = cfg.Stats.Damage

	if parent != nil {
		petal.TeamID = parent.TeamID
	}

	g.petals.Add(petal)
	return
}

func (petal *Petal) Update(deltaMilliseconds, damping, positionFactor float64) {
	if !petal.Health.IsAlive() {
		petal.Die()
		return
	}

	petal.GenericGameEntity.Update(damping, positionFactor)
	petal.Game.petalsSpatialHash.Insert(petal)

	if petal.Config.Special.Heal != nil {
		petal.HealTicker += deltaMilliseconds
	}

	if petal.Parent != nil && petal.Config.Special.ConstantHealPerSecond > 0 {
		petal.Parent.Health.Heal(petal.Config.Special.ConstantHealPerSecond * (deltaMilliseconds / 1000))
	}
}

func (petal *Petal) Die() {
	if petal.isDead {
		return
	}

	petal.isDead = true
	petal.Game.petals.Remove(petal.ID)
	petal.Game.petalsCache.Remove(petal.ID)
}

func (petal *Petal) Collide() {
	for _, other := range petal.Game.petalsSpatialHash.Retrieve(petal.AABB) {
		if other.ID == petal.ID || slices.Contains(petal.collisionIDs, other.ID) || !circleIntersects(petal.Position, petal.Size, other.Position, other.Size) {
			continue
		}

		CollideEntities(petal.GenericGameEntity, other.GenericGameEntity)
	}

	for _, mob := range petal.Game.mobsSpatialHash.Retrieve(petal.AABB) {
		if slices.Contains(petal.collisionIDs, mob.ID) || !circleIntersects(petal.Position, petal.Size, mob.Position, mob.Size) {
			continue
		}

		CollideEntities(petal.GenericGameEntity, mob.GenericGameEntity)
	}
}
