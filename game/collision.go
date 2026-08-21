package game

import (
	"math"

	"github.com/z46-dev/florr-legacy/shared/util"
)

type collisionResponse struct {
	enabled                  bool
	bounceA, bounceB         float64
	separationA, separationB float64
}

func (g *Game) Collide(deltaMilliseconds float64) {
	for _, flower := range g.flowers.Items {
		flower.Collide()
	}

	for _, petal := range g.petals.Items {
		petal.Collide()
	}

	for _, mob := range g.mobs.Items {
		mob.Collide()

		if mob.AI != nil {
			mob.AI.Think(deltaMilliseconds)
		}
	}
}

// CollideEntities resolves one entity-entity circle collision using the archived normal/reflection solver.
func CollideEntities(entityA, entityB *GenericGameEntity) (collided bool) {
	var (
		deltaX, deltaY             float64 = entityA.Position.X - entityB.Position.X, entityA.Position.Y - entityB.Position.Y
		distanceSquared, radiusSum float64 = deltaX*deltaX + deltaY*deltaY, entityA.Size + entityB.Size
		response                   collisionResponse
	)

	if distanceSquared <= 0 || distanceSquared >= radiusSum*radiusSum {
		return
	}

	response = collisionResponseFor(entityA, entityB)
	if !response.enabled {
		return
	}

	var (
		relativeVelocityX, relativeVelocityY float64 = entityA.Velocity.X - entityB.Velocity.X, entityA.Velocity.Y - entityB.Velocity.Y
		reflection                           float64 = (relativeVelocityX*deltaX + relativeVelocityY*deltaY) * -2 / distanceSquared
	)

	entityA.Velocity.X += deltaX * reflection * response.bounceA
	entityA.Velocity.Y += deltaY * reflection * response.bounceA
	entityB.Velocity.X -= deltaX * reflection * response.bounceB
	entityB.Velocity.Y -= deltaY * reflection * response.bounceB

	var (
		overlap, distDiv float64 = radiusSum - math.Sqrt(distanceSquared), 1 / math.Sqrt(distanceSquared)
	)

	entityA.Position.X += deltaX * distDiv * overlap * response.separationA
	entityA.Position.Y += deltaY * distDiv * overlap * response.separationA
	entityB.Position.X -= deltaX * distDiv * overlap * response.separationB
	entityB.Position.Y -= deltaY * distDiv * overlap * response.separationB
	entityA.collisionIDs, entityB.collisionIDs = append(entityA.collisionIDs, entityB.ID), append(entityB.collisionIDs, entityA.ID)
	syncEntityAABB(entityA)
	syncEntityAABB(entityB)

	if entityA.TeamID != entityB.TeamID {
		entityA.Health.Damage(entityB.Damage)
		entityB.Health.Damage(entityA.Damage)
		entityA.Health.LastDamageTick, entityB.Health.LastDamageTick = entityA.Game.Time, entityB.Game.Time

		if entityA.refMob != nil && entityA.refMob.AI != nil {
			entityA.refMob.AI.OnHostileCollision(entityB)
		}

		if entityB.refMob != nil && entityB.refMob.AI != nil {
			entityB.refMob.AI.OnHostileCollision(entityA)
		}
	}

	collided = true
	return
}

func collisionResponseFor(entityA, entityB *GenericGameEntity) (response collisionResponse) {
	response = collisionResponse{
		enabled: true,
		bounceA: 1,
		bounceB: 1,
	}

	switch {
	case entityA.Type == EntityTypeFlower && entityB.Type == EntityTypeFlower:
		if entityA.TeamID == entityB.TeamID {
			response.bounceA, response.bounceB = 0.35, 0.35
		}
	case entityA.Type == EntityTypePetal && entityB.Type == EntityTypePetal:
		if entityA.TeamID == entityB.TeamID {
			response.enabled = false
			return
		}

		response.bounceA, response.bounceB = 0.2, 0.2
	case entityA.Type == EntityTypeMob && entityB.Type == EntityTypeMob:
	case entityA.Type == EntityTypePetal && entityB.Type == EntityTypeMob:
		response.bounceA, response.bounceB = 0.2, 0
	case entityA.Type == EntityTypeMob && entityB.Type == EntityTypePetal:
		response.bounceA, response.bounceB = 0, 0.2
	default:
		if entityA.TeamID == entityB.TeamID {
			response.enabled = false
			return
		}
	}

	response.separationA, response.separationB = splitCollisionInfluence(
		collisionInfluence(entityB, entityA),
		collisionInfluence(entityA, entityB),
	)
	response.bounceA *= response.separationA
	response.bounceB *= response.separationB

	return
}

func collisionInfluence(pusher, receiver *GenericGameEntity) (influence float64) {
	var (
		pusherSize   float64 = max(pusher.Size, 0.001)
		receiverSize float64 = max(receiver.Size, 0.001)
		sizeRatio    float64 = pusherSize / receiverSize
	)

	influence = max(0, pusher.Density) * max(0, receiver.Pushability) * sizeRatio * sizeRatio
	return
}

func splitCollisionInfluence(influenceA, influenceB float64) (shareA, shareB float64) {
	var total float64 = influenceA + influenceB
	if total <= 0 {
		return
	}

	shareA, shareB = influenceA/total, influenceB/total
	return
}

func circleIntersects(posA *util.Vec2, radiusA float64, posB *util.Vec2, radiusB float64) (intersects bool) {
	intersects = posA.DistSquared(posB) < (radiusA+radiusB)*(radiusA+radiusB)
	return
}
