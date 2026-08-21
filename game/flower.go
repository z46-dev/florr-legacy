package game

import (
	"slices"

	"github.com/z46-dev/florr-legacy/shared/configs"
	"github.com/z46-dev/florr-legacy/shared/util"
)

func NewFlower(g *Game) (flower *Flower) {
	flower = &Flower{
		GenericGameEntity: NewEntity(g),
		NumSlots:          0,
		Slots:             make([]*Slot, 0),
		Movement:          util.Vec2_0(),
	}

	flower.Health.Set(baseFlowerHealth)
	flower.Damage = baseFlowerDamage
	flower.Size = baseFlowerSize
	flower.Type = EntityTypeFlower
	flower.refFlower = flower
	flower.TeamID = TeamID_Players

	var petals = []*configs.Petal{
		// configs.PetalTwin, configs.PetalFaster,
		// configs.PetalTwin, configs.PetalFaster,
		// configs.PetalRose,
		// configs.PetalTwin, configs.PetalFaster,
		// configs.PetalTwin, configs.PetalFaster,
		// configs.PetalRose,
		configs.PetalLeaf, configs.PetalLeaf,
		configs.PetalLeaf, configs.PetalLeaf,
		configs.PetalLeaf, configs.PetalLeaf,
		configs.PetalLeaf, configs.PetalLeaf,
		configs.PetalLeaf, configs.PetalLeaf,
	}

	for i, petal := range petals {
		flower.NumSlots++
		flower.Slots = append(flower.Slots, NewSlot(g, flower, i))
		flower.Slots[i].Define(petal)
	}

	g.flowers.Add(flower)
	return
}

func (flower *Flower) Update(deltaMilliseconds, damping, positionFactor float64) {
	if !flower.Health.IsAlive() {
		flower.Die()
		return
	}

	flower.Velocity.X += flower.Movement.X * SPEED_SCALE_FACTOR
	flower.Velocity.Y += flower.Movement.Y * SPEED_SCALE_FACTOR

	flower.OrbitAngle += (baseOrbitRadiansPerSecond + flower.OrbitRadiansPerSecond) / (1000.0 / deltaMilliseconds)

	var numClumps, currentClump int = 0, 0

	for _, slot := range flower.Slots {
		numClumps += slot.NumClumps
	}

	for _, slot := range flower.Slots {
		slot.Update(deltaMilliseconds, numClumps, currentClump)
		currentClump += slot.NumClumps
	}

	flower.GenericGameEntity.Update(damping, positionFactor)

	if flower.Position.X < -flower.Game.Width*0.5+flower.Size {
		flower.Position.X = -flower.Game.Width*0.5 + flower.Size
		flower.Velocity.X = 0
	} else if flower.Position.X > flower.Game.Width*0.5-flower.Size {
		flower.Position.X = flower.Game.Width*0.5 - flower.Size
		flower.Velocity.X = 0
	}

	if flower.Position.Y < -flower.Game.Height*0.5+flower.Size {
		flower.Position.Y = -flower.Game.Height*0.5 + flower.Size
		flower.Velocity.Y = 0
	} else if flower.Position.Y > flower.Game.Height*0.5-flower.Size {
		flower.Position.Y = flower.Game.Height*0.5 - flower.Size
		flower.Velocity.Y = 0
	}

	flower.Game.flowersSpatialHash.Insert(flower)
}

func (flower *Flower) Die() {
	if flower.isDead {
		return
	}

	flower.isDead = true
	flower.Game.flowers.Remove(flower.ID)
	flower.Game.flowersCache.Remove(flower.ID)

	for _, slot := range flower.Slots {
		for _, livePetal := range slot.Petals {
			if livePetal.Petal != nil {
				livePetal.Petal.Health.Current = -1
			}
		}
	}
}

func (flower *Flower) Collide() {
	for _, other := range flower.Game.flowersSpatialHash.Retrieve(flower.AABB) {
		if other.ID == flower.ID || slices.Contains(flower.collisionIDs, other.ID) || !circleIntersects(flower.Position, flower.Size, other.Position, other.Size) {
			continue
		}

		CollideEntities(flower.GenericGameEntity, other.GenericGameEntity)
	}

	for _, petal := range flower.Game.petalsSpatialHash.Retrieve(flower.AABB) {
		if slices.Contains(flower.collisionIDs, petal.ID) || !circleIntersects(flower.Position, flower.Size, petal.Position, petal.Size) {
			continue
		}

		CollideEntities(flower.GenericGameEntity, petal.GenericGameEntity)
	}

	for _, mob := range flower.Game.mobsSpatialHash.Retrieve(flower.AABB) {
		if slices.Contains(flower.collisionIDs, mob.ID) || !circleIntersects(flower.Position, flower.Size, mob.Position, mob.Size) {
			continue
		}

		CollideEntities(flower.GenericGameEntity, mob.GenericGameEntity)
	}
}
