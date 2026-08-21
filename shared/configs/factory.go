package configs

import "github.com/z46-dev/florr-legacy/client/render"

// ==================================================
// ==================== Rarities ====================
// ==================================================

func NewRarity(id RarityID, name string, color string) (r *Rarity) {
	r = &Rarity{
		ID:    id,
		Name:  name,
		Color: render.HexToColor(color),
	}

	Rarities[id] = r
	return
}

// ==================================================
// ===================== Petals =====================
// ==================================================

func NewPetal(id PetalID, name string, rarity RarityID, reloadMilliseconds uint32) (p *Petal) {
	p = &Petal{
		ID: id,
		Meta: &PetalMeta{
			Name:   name,
			Rarity: rarity,
		},
		Stats: &PetalStats{
			ReloadMilliseconds: reloadMilliseconds,
			Radius:             10,
			Orbit: &PetalOrbit{
				Neutral: 65,
				Attack:  130,
				Defend:  40,
			},
		},
		Special: &PetalSpecial{
			Count:     1,
			ClumpSize: 1,
		},
	}

	Petals[id] = p
	return
}

func (p *Petal) SetDescription(description string) (self *Petal) {
	p.Meta.Description = description
	self = p
	return
}

func (p *Petal) SetHealthAndDamage(health, damage float64) (self *Petal) {
	p.Stats.Health = health
	p.Stats.Damage = damage
	self = p
	return
}

func (p *Petal) SetRadius(radius float64) (self *Petal) {
	p.Stats.Radius = radius
	self = p
	return
}

func (p *Petal) SetMulti(count, clumpSize int) (self *Petal) {
	p.Special.Count = count
	p.Special.ClumpSize = clumpSize
	self = p
	return
}

func (p *Petal) SetOrbit(neutral, attack, defend float64) (self *Petal) {
	p.Stats.Orbit.Neutral = neutral
	p.Stats.Orbit.Attack = attack
	p.Stats.Orbit.Defend = defend
	self = p
	return
}

func (p *Petal) SetHeal(heal, delayMilliseconds float64) (self *Petal) {
	p.Special.Heal = &PetalHeal{
		Amount: heal,
		Delay:  delayMilliseconds,
	}

	self = p
	return
}

func (p *Petal) SetExtraRotation(rotationRadiansPerSecond float64) (self *Petal) {
	p.Special.RotationRadiansPerSecond = rotationRadiansPerSecond
	self = p
	return
}

func (p *Petal) SetConstantHeal(healPerSecond float64) (self *Petal) {
	p.Special.ConstantHealPerSecond = healPerSecond
	self = p
	return
}

// ==================================================
// ====================== Mobs ======================
// ==================================================

func NewMob(id MobID, name string, rarity RarityID, health, damage, radius float64) (m *Mob) {
	m = &Mob{
		ID: id,
		Meta: &MobMeta{
			Name:   name,
			Rarity: rarity,
		},
		Stats: &MobStats{
			Health:      health,
			Damage:      damage,
			Radius:      radius,
			Density:     1,
			Pushability: 1,
		},
	}

	Mobs[id] = m
	return
}

func (m *Mob) SetDescription(description string) (self *Mob) {
	m.Meta.Description = description
	self = m
	return
}

func (m *Mob) AddDrop(petalID PetalID, chance float64) (self *Mob) {
	m.Drops = append(m.Drops, &MobDrop{
		PetalID: petalID,
		Chance:  chance,
	})

	self = m
	return
}

func (m *Mob) SetVariance(variance *MobVariance) (self *Mob) {
	m.Variance = variance
	self = m
	return
}

func (m *Mob) SetPhysics(pushability, density float64) (self *Mob) {
	m.Stats.Pushability = pushability
	m.Stats.Density = density
	self = m
	return
}

func (m *Mob) SetAI(ai *MobAI) (self *Mob) {
	m.AI = ai
	self = m
	return
}
