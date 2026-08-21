package game

import (
	"math"

	"github.com/z46-dev/florr-legacy/shared"
	"github.com/z46-dev/florr-legacy/shared/util"
)

const FLOAT_EPSILON = 0.0001

func differs(a, b float64) bool {
	return math.Abs(a-b) > FLOAT_EPSILON
}

func NewCamera(cilent *Client) (camera *Camera) {
	camera = &Camera{
		Client:      cilent,
		Position:    util.Vec2_0(),
		FOV:         6400.0 / 3.5,
		FlowersSeen: make(map[uint64]uint64),
		PetalsSeen:  make(map[uint64]uint64),
		MobsSeen:    make(map[uint64]uint64),
	}

	return
}

// IsInView reports whether an axis-aligned box intersects this camera.
func (c *Camera) IsInView(aabb *util.AABB) (inside bool) {
	var halfFOV float64 = c.FOV * 0.5
	inside = aabb.X1 < c.Position.X+halfFOV &&
		aabb.X2 > c.Position.X-halfFOV &&
		aabb.Y1 < c.Position.Y+halfFOV &&
		aabb.Y2 > c.Position.Y-halfFOV

	return
}

func (c *Camera) See() {
	if c.Client.Body != nil {
		c.Position.X, c.Position.Y = c.Client.Body.Position.X, c.Client.Body.Position.Y
	}

	var (
		w       *shared.Writer = &shared.Writer{}
		halfFOV float64        = c.FOV * 0.5
		view    *util.AABB     = &util.AABB{
			X1: c.Position.X - halfFOV,
			Y1: c.Position.Y - halfFOV,
			X2: c.Position.X + halfFOV,
			Y2: c.Position.Y + halfFOV,
		}

		flowersNow map[uint64]struct{} = make(map[uint64]struct{})
		petalsNow  map[uint64]struct{} = make(map[uint64]struct{})
		mobsNow    map[uint64]struct{} = make(map[uint64]struct{})

		flowersSeenBefore map[uint64]uint64 = cloneSeenEntities(c.FlowersSeen)
		petalsSeenBefore  map[uint64]uint64 = cloneSeenEntities(c.PetalsSeen)
		mobsSeenBefore    map[uint64]uint64 = cloneSeenEntities(c.MobsSeen)
	)

	w.SetU8(shared.GameMsgTypeWorldStateUpdate)

	w.SetU64(uint64(len(c.Client.Game.flowers.Items)))
	w.SetU64(uint64(len(c.Client.Game.petals.Items)))
	w.SetU64(uint64(len(c.Client.Game.mobs.Items)))
	w.SetF64(c.Client.Game.mspt)
	w.SetU64(c.Client.Game.fps)

	w.SetF32(float32(c.Position.X))
	w.SetF32(float32(c.Position.Y))
	w.SetF32(float32(c.FOV))
	if c.Client.Body != nil && !c.Client.Body.isDead {
		w.SetU64(c.Client.Body.ID)
	} else {
		w.SetU64(0)
	}

	for _, flower := range c.Client.Game.flowersSpatialHash.Retrieve(view) {
		if _, exists := flowersNow[flower.ID]; exists {
			continue
		}

		if !flower.isDead && c.IsInView(flower.GetAABB()) {
			flowersNow[flower.ID] = struct{}{}
			c.seeFlower(w, flower)
		}
	}

	w.SetU64(0)

	var stillSeen bool
	for id := range c.FlowersSeen {
		if _, stillSeen = flowersNow[id]; !stillSeen {
			w.SetU64(id)
			delete(c.FlowersSeen, id)
		}
	}

	w.SetU64(0)

	for _, petal := range c.Client.Game.petalsSpatialHash.Retrieve(view) {
		if _, exists := petalsNow[petal.ID]; exists {
			continue
		}

		if !petal.isDead && c.IsInView(petal.GetAABB()) {
			petalsNow[petal.ID] = struct{}{}
			c.seePetal(w, petal)
		}
	}

	w.SetU64(0)

	for id := range c.PetalsSeen {
		if _, stillSeen = petalsNow[id]; !stillSeen {
			w.SetU64(id)
			delete(c.PetalsSeen, id)
		}
	}

	w.SetU64(0)

	for _, mob := range c.Client.Game.mobsSpatialHash.Retrieve(view) {
		if _, exists := mobsNow[mob.ID]; exists {
			continue
		}

		if !mob.isDead && c.IsInView(mob.GetAABB()) {
			mobsNow[mob.ID] = struct{}{}
			c.seeMob(w, mob)
		}
	}

	w.SetU64(0)

	for id := range c.MobsSeen {
		if _, stillSeen = mobsNow[id]; !stillSeen {
			w.SetU64(id)
			delete(c.MobsSeen, id)
		}
	}

	w.SetU64(0)

	if err := c.Client.Peer.Send(w.Bytes); err != nil {
		c.FlowersSeen = flowersSeenBefore
		c.PetalsSeen = petalsSeenBefore
		c.MobsSeen = mobsSeenBefore
		c.Client.Game.L.Warningf("Failed to send world state (%d bytes): %v", len(w.Bytes), err)
	}
}

// cloneSeenEntities snapshots camera visibility so failed world updates can be retried.
func cloneSeenEntities(seen map[uint64]uint64) (clone map[uint64]uint64) {
	clone = make(map[uint64]uint64, len(seen))
	for id, asOf := range seen {
		clone[id] = asOf
	}

	return
}

func (c *Camera) seeFlower(w *shared.Writer, flower *Flower) {
	var cache *FlowerCache

	if cache = c.Client.Game.flowersCache.Items[flower.ID]; cache == nil {
		cache = &FlowerCache{}
		cache.ID = flower.ID
		cache.Size = flower.Size
		c.Client.Game.flowersCache.Add(cache)
	}

	if cache.AsOf != c.Client.Game.Time {
		cache.PosChanged = false
		cache.SizeChanged = false
		cache.RotChanged = false
		cache.HealthChanged = false
		cache.AttackModeChanged = false
		cache.StatusFlagsChanged = false
		var newFlags uint8 = 0

		if differs(cache.X, flower.Position.X) || differs(cache.Y, flower.Position.Y) {
			cache.PosChanged = true
			cache.X = flower.Position.X
			cache.Y = flower.Position.Y
		}

		if differs(cache.Size, flower.Size) {
			cache.SizeChanged = true
			cache.Size = flower.Size
		}

		if differs(cache.Rotation, flower.Facing) {
			cache.RotChanged = true
			cache.Rotation = flower.Facing
		}

		var newHealth = flower.Health.Ratio()
		if newHealth != cache.HealthRatio {
			cache.HealthChanged = true

			if newHealth < cache.HealthRatio {
				newFlags |= 1 << 0
			}

			cache.HealthRatio = newHealth
		}

		if cache.AttackMode != flower.AttackMode {
			cache.AttackModeChanged = true
			cache.AttackMode = flower.AttackMode
		}

		if cache.StatusFlags != newFlags {
			cache.StatusFlagsChanged = true
			cache.StatusFlags = newFlags
		}

		if cache.PosChanged || cache.SizeChanged || cache.RotChanged || cache.HealthChanged || cache.AttackModeChanged || cache.StatusFlagsChanged {
			cache.AsOf = c.Client.Game.Time

			if cache.New != nil {
				cache.New.Reset()
			}

			if cache.Old != nil {
				cache.Old.Reset()
			}
		}
	}

	var (
		seenAt uint64
		seen   bool
	)

	if seenAt, seen = c.FlowersSeen[flower.ID]; !seen {
		if cache.New == nil || cache.New.GetLength() == 0 {
			cache.New = &shared.Writer{}
			cache.New.SetU64(cache.ID)
			cache.New.SetU8(0)
			cache.New.SetF32(float32(cache.X))
			cache.New.SetF32(float32(cache.Y))
			cache.New.SetF32(float32(cache.Size))
			cache.New.SetF32(float32(cache.Rotation))
			cache.New.SetStringUTF8(flower.Name)
			cache.New.SetU8(uint8(cache.HealthRatio * 255.0))
			cache.New.SetU8(uint8(cache.AttackMode))
			cache.New.SetU8(cache.StatusFlags)
		}

		c.FlowersSeen[flower.ID] = cache.AsOf
		w.Append(cache.New)
		return
	}

	if seenAt == cache.AsOf || (!cache.PosChanged && !cache.SizeChanged && !cache.RotChanged && !cache.HealthChanged && !cache.AttackModeChanged && !cache.StatusFlagsChanged) {
		return
	}

	if cache.Old == nil || cache.Old.GetLength() == 0 {
		cache.Old = &shared.Writer{}
		cache.Old.SetU64(cache.ID)
		cache.Old.SetU8(1)

		var flags uint8 = 0

		if cache.PosChanged {
			flags |= 1 << 0
		}

		if cache.SizeChanged {
			flags |= 1 << 1
		}

		if cache.RotChanged {
			flags |= 1 << 2
		}

		if cache.HealthChanged {
			flags |= 1 << 3
		}

		if cache.AttackModeChanged {
			flags |= 1 << 4
		}

		if cache.StatusFlagsChanged {
			flags |= 1 << 5
		}

		cache.Old.SetU8(flags)

		if cache.PosChanged {
			cache.Old.SetF32(float32(cache.X))
			cache.Old.SetF32(float32(cache.Y))
		}

		if cache.SizeChanged {
			cache.Old.SetF32(float32(cache.Size))
		}

		if cache.RotChanged {
			cache.Old.SetF32(float32(cache.Rotation))
		}

		if cache.HealthChanged {
			cache.Old.SetU8(uint8(cache.HealthRatio * 255.0))
		}

		if cache.AttackModeChanged {
			cache.Old.SetU8(uint8(cache.AttackMode))
		}

		if cache.StatusFlagsChanged {
			cache.Old.SetU8(cache.StatusFlags)
		}
	}

	c.FlowersSeen[flower.ID] = cache.AsOf
	w.Append(cache.Old)
}

func (c *Camera) seePetal(w *shared.Writer, petal *Petal) {
	var cache *PetalCache

	if cache = c.Client.Game.petalsCache.Items[petal.ID]; cache == nil {
		cache = &PetalCache{}
		cache.ID = petal.ID
		cache.Size = petal.Size
		c.Client.Game.petalsCache.Add(cache)
	}

	if cache.AsOf != c.Client.Game.Time {
		cache.PosChanged = false
		cache.SizeChanged = false
		cache.RotChanged = false
		cache.StatusFlagsChanged = false

		var newFlags uint8 = 0

		if differs(cache.X, petal.Position.X) || differs(cache.Y, petal.Position.Y) {
			cache.PosChanged = true
			cache.X = petal.Position.X
			cache.Y = petal.Position.Y
		}

		if differs(cache.Size, petal.Size) {
			cache.SizeChanged = true
			cache.Size = petal.Size
		}

		if differs(cache.Rotation, petal.Facing) {
			cache.RotChanged = true
			cache.Rotation = petal.Facing
		}

		var newHealth = petal.Health.Ratio()
		if newHealth != cache.HealthRatio {
			if newHealth < cache.HealthRatio {
				newFlags |= 1 << 0
			}

			cache.HealthRatio = newHealth
		}

		if cache.StatusFlags != newFlags {
			cache.StatusFlagsChanged = true
			cache.StatusFlags = newFlags
		}

		if cache.PosChanged || cache.SizeChanged || cache.RotChanged || cache.StatusFlagsChanged {
			cache.AsOf = c.Client.Game.Time

			if cache.New != nil {
				cache.New.Reset()
			}

			if cache.Old != nil {
				cache.Old.Reset()
			}
		}
	}

	var (
		seenAt uint64
		seen   bool
	)

	if seenAt, seen = c.PetalsSeen[petal.ID]; !seen {
		if cache.New == nil || cache.New.GetLength() == 0 {
			cache.New = &shared.Writer{}
			cache.New.SetU64(cache.ID)
			cache.New.SetU8(0)
			cache.New.SetF32(float32(cache.X))
			cache.New.SetF32(float32(cache.Y))
			cache.New.SetF32(float32(cache.Size))
			cache.New.SetF32(float32(cache.Rotation))
			cache.New.SetU8(uint8(petal.Config.ID))
			cache.New.SetU8(uint8(petal.Config.Meta.Rarity))
			cache.New.SetU8(cache.StatusFlags)
		}

		c.PetalsSeen[petal.ID] = cache.AsOf
		w.Append(cache.New)
		return
	}

	if seenAt == cache.AsOf || (!cache.PosChanged && !cache.SizeChanged && !cache.RotChanged && !cache.StatusFlagsChanged) {
		return
	}

	if cache.Old == nil || cache.Old.GetLength() == 0 {
		cache.Old = &shared.Writer{}
		cache.Old.SetU64(cache.ID)
		cache.Old.SetU8(1)

		var flags uint8 = 0

		if cache.PosChanged {
			flags |= 1 << 0
		}

		if cache.SizeChanged {
			flags |= 1 << 1
		}

		if cache.RotChanged {
			flags |= 1 << 2
		}

		if cache.StatusFlagsChanged {
			flags |= 1 << 3
		}

		cache.Old.SetU8(flags)

		if cache.PosChanged {
			cache.Old.SetF32(float32(cache.X))
			cache.Old.SetF32(float32(cache.Y))
		}

		if cache.SizeChanged {
			cache.Old.SetF32(float32(cache.Size))
		}

		if cache.RotChanged {
			cache.Old.SetF32(float32(cache.Rotation))
		}

		if cache.StatusFlagsChanged {
			cache.Old.SetU8(cache.StatusFlags)
		}
	}

	c.PetalsSeen[petal.ID] = cache.AsOf
	w.Append(cache.Old)
}

func (c *Camera) seeMob(w *shared.Writer, mob *Mob) {
	var cache *MobCache

	if cache = c.Client.Game.mobsCache.Items[mob.ID]; cache == nil {
		cache = &MobCache{}
		cache.ID = mob.ID
		cache.X = mob.Position.X
		cache.Y = mob.Position.Y
		cache.Size = mob.Size
		cache.Rotation = mob.Facing
		cache.HealthRatio = mob.Health.Ratio()
		if mob.AI != nil && mob.AI.Target != nil {
			cache.StatusFlags = 1 << 1
		}

		c.Client.Game.mobsCache.Add(cache)
	}

	if cache.AsOf != c.Client.Game.Time {
		cache.PosChanged = false
		cache.SizeChanged = false
		cache.RotChanged = false
		cache.HealthChanged = false
		cache.StatusFlagsChanged = false

		var newFlags uint8 = 0

		if differs(cache.X, mob.Position.X) || differs(cache.Y, mob.Position.Y) {
			cache.PosChanged = true
			cache.X = mob.Position.X
			cache.Y = mob.Position.Y
		}

		if differs(cache.Size, mob.Size) {
			cache.SizeChanged = true
			cache.Size = mob.Size
		}

		if differs(cache.Rotation, mob.Facing) {
			cache.RotChanged = true
			cache.Rotation = mob.Facing
		}

		// var healthRatio uint8 = uint8(mob.Health.Ratio() * 255.0)
		// if healthRatio != cache.HealthRatio {
		// 	cache.HealthChanged = true
		// 	cache.HealthRatio = healthRatio
		// }

		var newHealth = mob.Health.Ratio()
		if newHealth != cache.HealthRatio {
			cache.HealthChanged = true

			if newHealth < cache.HealthRatio {
				newFlags |= 1 << 0
			}

			cache.HealthRatio = newHealth
		}

		if mob.AI != nil && mob.AI.Target != nil {
			newFlags |= 1 << 1
		}

		if cache.StatusFlags != newFlags {
			cache.StatusFlagsChanged = true
			cache.StatusFlags = newFlags
		}

		if cache.PosChanged || cache.SizeChanged || cache.RotChanged || cache.HealthChanged || cache.StatusFlagsChanged {
			cache.AsOf = c.Client.Game.Time

			if cache.New != nil {
				cache.New.Reset()
			}

			if cache.Old != nil {
				cache.Old.Reset()
			}
		}
	}

	var (
		seenAt uint64
		seen   bool
	)

	if seenAt, seen = c.MobsSeen[mob.ID]; !seen || c.Client.Game.Time%TicksPerSecond == 0 {
		if cache.New == nil || cache.New.GetLength() == 0 {
			cache.New = &shared.Writer{}
			cache.New.SetU64(cache.ID)
			cache.New.SetU8(0)
			cache.New.SetF32(float32(mob.Position.X))
			cache.New.SetF32(float32(mob.Position.Y))
			cache.New.SetF32(float32(mob.Size))
			cache.New.SetF32(float32(mob.Facing))
			cache.New.SetU8(uint8(mob.Health.Ratio() * 255.0))
			cache.New.SetU8(uint8(mob.Config.ID))
			cache.New.SetU8(uint8(mob.Config.Meta.Rarity))
			cache.New.SetU8(cache.StatusFlags)
		}

		c.MobsSeen[mob.ID] = c.Client.Game.Time
		w.Append(cache.New)
		return
	}

	if seenAt == cache.AsOf || (!cache.PosChanged && !cache.SizeChanged && !cache.RotChanged && !cache.HealthChanged && !cache.StatusFlagsChanged) {
		return
	}

	if cache.Old == nil || cache.Old.GetLength() == 0 {
		cache.Old = &shared.Writer{}
		cache.Old.SetU64(cache.ID)
		cache.Old.SetU8(1)

		var flags uint8 = 0

		if cache.PosChanged {
			flags |= 1 << 0
		}

		if cache.SizeChanged {
			flags |= 1 << 1
		}

		if cache.RotChanged {
			flags |= 1 << 2
		}

		if cache.HealthChanged {
			flags |= 1 << 3
		}

		if cache.StatusFlagsChanged {
			flags |= 1 << 4
		}

		cache.Old.SetU8(flags)

		if cache.PosChanged {
			cache.Old.SetF32(float32(cache.X))
			cache.Old.SetF32(float32(cache.Y))
		}

		if cache.SizeChanged {
			cache.Old.SetF32(float32(cache.Size))
		}

		if cache.RotChanged {
			cache.Old.SetF32(float32(cache.Rotation))
		}

		if cache.HealthChanged {
			cache.Old.SetU8(uint8(cache.HealthRatio * 255.0))
		}

		if cache.StatusFlagsChanged {
			cache.Old.SetU8(cache.StatusFlags)
		}
	}

	c.MobsSeen[mob.ID] = cache.AsOf
	w.Append(cache.Old)
}
