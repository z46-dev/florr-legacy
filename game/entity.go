package game

import "github.com/z46-dev/florr-legacy/shared/util"

func NewEntity(g *Game) (entity *GenericGameEntity) {
	entity = &GenericGameEntity{
		Game:        g,
		ID:          g.nextEntityID,
		TeamID:      g.nextEntityID,
		Position:    util.NewVec2(0, 0),
		Velocity:    util.NewVec2(0, 0),
		Size:        baseFlowerSize,
		Facing:      0,
		Density:     1,
		Pushability: 1,
		Health:      NewHealthComponent(1),
		AABB:        &util.AABB{},
	}

	g.nextEntityID++
	return
}

func (entity *GenericGameEntity) GetAABB() *util.AABB {
	return entity.AABB
}

// Universal physics update. If an entity which embeds GenericGameEntity has its own
// Update() method, it should call this method at the start of its own Update() method.
func (entity *GenericGameEntity) Update(damping, positionFactor float64) {
	entity.Position.X += entity.Velocity.X * positionFactor
	entity.Position.Y += entity.Velocity.Y * positionFactor
	entity.Velocity.X *= damping
	entity.Velocity.Y *= damping
	syncEntityAABB(entity)
	entity.collisionIDs = entity.collisionIDs[:0]
}

func syncEntityAABB(entity *GenericGameEntity) {
	entity.AABB.X1 = entity.Position.X - entity.Size
	entity.AABB.Y1 = entity.Position.Y - entity.Size
	entity.AABB.X2 = entity.Position.X + entity.Size
	entity.AABB.Y2 = entity.Position.Y + entity.Size
}

func NewHealthComponent(health float64) (hc *HealthComponent) {
	hc = &HealthComponent{
		Current:      health,
		Maximum:      health,
		Regeneration: 0,
	}

	return
}

func (hc *HealthComponent) IsAlive() bool {
	return hc.Current > 0
}

func (hc *HealthComponent) Damage(amount float64) {
	hc.Current -= amount
}

func (hc *HealthComponent) Heal(amount float64) {
	hc.Current = min(hc.Current+amount, hc.Maximum)
}

func (hc *HealthComponent) Set(health float64) {
	var diff float64 = health - hc.Maximum
	hc.Maximum = health
	hc.Current += diff
}

func (hc *HealthComponent) Ratio() float64 {
	if hc.Maximum <= 0 || hc.Current <= 0 {
		return 0
	}

	return min(hc.Current/hc.Maximum, 1)
}
