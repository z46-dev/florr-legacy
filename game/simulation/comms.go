package simulation

import (
	"time"

	"github.com/z46-dev/florr-legacy/shared"
	"github.com/z46-dev/florr-legacy/shared/configs"
)

func findFlower(collection *shared.Collection[*Flower], id uint64) (flower *Flower) {
	if flower = collection.Items[id]; flower != nil {
		return
	}

	for _, pending := range collection.ToAppend {
		if pending != nil && pending.ID == id {
			flower = pending
			return
		}
	}

	return
}

func findPetal(collection *shared.Collection[*Petal], id uint64) (petal *Petal) {
	if petal = collection.Items[id]; petal != nil {
		return
	}

	for _, pending := range collection.ToAppend {
		if pending != nil && pending.ID == id {
			petal = pending
			return
		}
	}

	return
}

func findMob(collection *shared.Collection[*Mob], id uint64) (mob *Mob) {
	if mob = collection.Items[id]; mob != nil {
		return
	}

	for _, pending := range collection.ToAppend {
		if pending != nil && pending.ID == id {
			mob = pending
			return
		}
	}

	return
}

func (s *Simulation) Spawn(name string) (err error) {
	var writer *shared.Writer = &shared.Writer{}
	writer.SetU8(shared.GameMsgTypeSpawn)
	writer.SetStringUTF8(name)
	err = s.Client.Send(writer.Bytes)
	return
}

func (s *Simulation) OnMessage(r *shared.Reader) {
	var msgType uint8 = r.GetU8()

	switch msgType {
	case shared.GameMsgTypeWorldStateUpdate:
		s.Flowers.Lock.Lock()
		s.Petals.Lock.Lock()
		s.Mobs.Lock.Lock()
		defer s.Mobs.Lock.Unlock()
		defer s.Petals.Lock.Unlock()
		defer s.Flowers.Lock.Unlock()

		for _, flower := range s.Flowers.Items {
			flower.PrevX, flower.PrevY, flower.PrevSize, flower.PrevRotation = flower.RealX, flower.RealY, flower.RealSize, flower.RealRotation
		}

		for _, petal := range s.Petals.Items {
			petal.PrevX, petal.PrevY, petal.PrevSize, petal.PrevRotation = petal.RealX, petal.RealY, petal.RealSize, petal.RealRotation
		}

		for _, mob := range s.Mobs.Items {
			mob.PrevX, mob.PrevY, mob.PrevSize, mob.PrevRotation = mob.RealX, mob.RealY, mob.RealSize, mob.RealRotation
		}

		s.WorldLerpAlpha = 0
		s.LastWorldStateAt = time.Now()

		s.serverStats.numFlowers = int(r.GetU64())
		s.serverStats.numPetals = int(r.GetU64())
		s.serverStats.numMobs = int(r.GetU64())
		s.serverStats.mspt = r.GetF64()
		s.serverStats.fps = r.GetU64()

		var (
			cameraX   float32 = r.GetF32()
			cameraY   float32 = r.GetF32()
			cameraFOV float32 = r.GetF32()
		)

		s.Camera.PrevX, s.Camera.PrevY = s.Camera.RealX, s.Camera.RealY
		s.Camera.PrevFieldOfView = s.Camera.RealFieldOfView
		s.Camera.RealX, s.Camera.RealY, s.Camera.RealFieldOfView = cameraX, cameraY, cameraFOV
		s.PlayerID = r.GetU64()

		var id uint64

		// Flowers
		for id = r.GetU64(); id != 0; id = r.GetU64() {
			if r.GetU8() == 0 {
				var flower *Flower = &Flower{MouthCurve: 5}
				flower.ID = id
				flower.RealX, flower.RealY = r.GetF32(), r.GetF32()
				flower.RealSize, flower.RealRotation = r.GetF32(), r.GetF32()
				flower.Name = r.GetStringUTF8()
				flower.RealHealthRatio = float32(r.GetU8()) / 255.0
				flower.AttackMode = configs.AttackMode(r.GetU8())

				var statusFlags uint8 = r.GetU8()
				flower.Damaged = statusFlags&0x01 != 0

				flower.PrevX, flower.PrevY = flower.RealX, flower.RealY
				flower.PrevSize, flower.PrevRotation = flower.RealSize, flower.RealRotation
				flower.X, flower.Y, flower.Size, flower.Rotation = flower.RealX, flower.RealY, flower.RealSize, flower.RealRotation
				flower.HealthRatio, flower.SubHealthRatio, flower.LastDamagedAt = flower.RealHealthRatio, flower.RealHealthRatio, 0
				flower.Alpha = 1.0

				s.Flowers.Items[flower.ID] = flower
				// s.L.Infof("New flower: ID=%d, Name=%s, X=%.2f, Y=%.2f, Size=%.2f, Rotation=%.2f, HealthRatio=%.2f\n", flower.ID, flower.Name, flower.RealX, flower.RealY, flower.RealSize, flower.RealRotation, flower.HealthRatio)
			} else {
				var flags uint8 = r.GetU8()
				var flower *Flower = findFlower(s.Flowers, id)
				if flower == nil {
					s.L.Warningf("Discard flower update for unknown ID=%d\n", id)
				}

				if flags&0x01 != 0 {
					var x, y float32 = r.GetF32(), r.GetF32()
					if flower != nil {
						flower.PrevX, flower.PrevY = flower.RealX, flower.RealY
						flower.RealX, flower.RealY = x, y
					}
				}

				if flags&0x02 != 0 {
					var size float32 = r.GetF32()
					if flower != nil {
						flower.PrevSize = flower.RealSize
						flower.RealSize = size
					}
				}

				if flags&0x04 != 0 {
					var rotation float32 = r.GetF32()
					if flower != nil {
						flower.PrevRotation = flower.RealRotation
						flower.RealRotation = rotation
					}
				}

				if flags&0x08 != 0 {
					var healthRatio float32 = float32(r.GetU8()) / 255.0
					if flower != nil {
						if flower.RealHealthRatio = healthRatio; flower.RealHealthRatio < flower.HealthRatio {
							flower.LastDamagedAt = HEALTH_DELTA_DURATION
						}
					}
				}

				if flags&0x10 != 0 {
					var attackMode configs.AttackMode = configs.AttackMode(r.GetU8())
					if flower != nil {
						flower.AttackMode = attackMode
					}
				}

				if flags&0x20 != 0 {
					var statusFlags uint8 = r.GetU8()
					if flower != nil {
						flower.Damaged = statusFlags&0x01 != 0
					}
				}
			}
		}

		s.Player = findFlower(s.Flowers, s.PlayerID)

		for id = r.GetU64(); id != 0; id = r.GetU64() {
			var flower *Flower = findFlower(s.Flowers, id)
			if flower != nil {
				// s.L.Infof("Remove flower: ID=%d, Name=%s\n", flower.ID, flower.Name)
				flower.IsDead = true
			} else {
				s.L.Warningf("Discard flower removal for unknown ID=%d\n", id)
			}
		}

		// Petals

		for id = r.GetU64(); id != 0; id = r.GetU64() {
			if r.GetU8() == 0 {
				var petal *Petal = &Petal{}
				petal.ID = id
				petal.RealX, petal.RealY = r.GetF32(), r.GetF32()
				petal.RealSize, petal.RealRotation = r.GetF32(), r.GetF32()
				petal.PetalID = configs.PetalID(r.GetU8())
				petal.RarityID = configs.RarityID(r.GetU8())

				var statusFlags uint8 = r.GetU8()
				petal.Damaged = statusFlags&0x01 != 0

				petal.PrevX, petal.PrevY = petal.RealX, petal.RealY
				petal.PrevSize, petal.PrevRotation = petal.RealSize, petal.RealRotation
				petal.X, petal.Y, petal.Size, petal.Rotation = petal.RealX, petal.RealY, petal.RealSize, petal.RealRotation
				petal.Alpha = 1.0

				s.Petals.Items[petal.ID] = petal
				// s.L.Infof("New petal: ID=%d, PetalID=%d, RarityID=%d, X=%.2f, Y=%.2f, Size=%.2f, Rotation=%.2f\n", petal.ID, petal.PetalID, petal.RarityID, petal.RealX, petal.RealY, petal.RealSize, petal.RealRotation)
			} else {
				var flags uint8 = r.GetU8()
				var petal *Petal = findPetal(s.Petals, id)
				if petal == nil {
					s.L.Warningf("Discard petal update for unknown ID=%d\n", id)
				}

				if flags&0x01 != 0 {
					var x, y float32 = r.GetF32(), r.GetF32()
					if petal != nil {
						petal.PrevX, petal.PrevY = petal.RealX, petal.RealY
						petal.RealX, petal.RealY = x, y
					}
				}

				if flags&0x02 != 0 {
					var size float32 = r.GetF32()
					if petal != nil {
						petal.PrevSize = petal.RealSize
						petal.RealSize = size
					}
				}

				if flags&0x04 != 0 {
					var rotation float32 = r.GetF32()
					if petal != nil {
						petal.PrevRotation = petal.RealRotation
						petal.RealRotation = rotation
					}
				}

				if flags&0x08 != 0 {
					var statusFlags uint8 = r.GetU8()
					if petal != nil {
						petal.Damaged = statusFlags&0x01 != 0
					}
				}
			}
		}

		for id = r.GetU64(); id != 0; id = r.GetU64() {
			// s.Petals.Remove(id)

			var petal *Petal = findPetal(s.Petals, id)
			if petal != nil {
				// s.L.Infof("Remove petal: ID=%d, PetalID=%d, RarityID=%d\n", petal.ID, petal.PetalID, petal.RarityID)
				petal.IsDead = true
			} else {
				s.L.Warningf("Discard petal removal for unknown ID=%d\n", id)
			}
		}

		// Mobs
		for id = r.GetU64(); id != 0; id = r.GetU64() {
			if r.GetU8() == 0 {
				var mob *Mob = findMob(s.Mobs, id)
				if mob == nil {
					mob = &Mob{}
				}

				mob.ID = id
				mob.RealX, mob.RealY = r.GetF32(), r.GetF32()
				mob.RealSize, mob.RealRotation = r.GetF32(), r.GetF32()
				mob.RealHealthRatio = float32(r.GetU8()) / 255.0
				mob.MobID = configs.MobID(r.GetU8())
				mob.RarityID = configs.RarityID(r.GetU8())

				var statusFlags uint8 = r.GetU8()
				mob.Damaged = statusFlags&0x01 != 0
				mob.Attacking = statusFlags&0x02 != 0

				mob.PrevX, mob.PrevY = mob.RealX, mob.RealY
				mob.PrevSize, mob.PrevRotation = mob.RealSize, mob.RealRotation
				mob.X, mob.Y, mob.Size, mob.Rotation = mob.RealX, mob.RealY, mob.RealSize, mob.RealRotation
				mob.HealthRatio, mob.SubHealthRatio, mob.LastDamagedAt = mob.RealHealthRatio, mob.RealHealthRatio, 0
				mob.Alpha = 1.0

				s.Mobs.Items[mob.ID] = mob
				// s.L.Infof("New mob: ID=%d, MobID=%d, RarityID=%d, X=%.2f, Y=%.2f, Size=%.2f, Rotation=%.2f, HealthRatio=%.2f\n", mob.ID, mob.MobID, mob.RarityID, mob.RealX, mob.RealY, mob.RealSize, mob.RealRotation, mob.HealthRatio)
			} else {
				var flags uint8 = r.GetU8()
				var mob *Mob = findMob(s.Mobs, id)
				if mob == nil {
					s.L.Warningf("Discard mob update for unknown ID=%d\n", id)
				}

				if flags&0x01 != 0 {
					var x, y float32 = r.GetF32(), r.GetF32()
					if mob != nil {
						mob.PrevX, mob.PrevY = mob.RealX, mob.RealY
						mob.RealX, mob.RealY = x, y
					}
				}

				if flags&0x02 != 0 {
					var size float32 = r.GetF32()
					if mob != nil {
						mob.PrevSize = mob.RealSize
						mob.RealSize = size
					}
				}

				if flags&0x04 != 0 {
					var rotation float32 = r.GetF32()
					if mob != nil {
						mob.PrevRotation = mob.RealRotation
						mob.RealRotation = rotation
					}
				}

				if flags&0x08 != 0 {
					var healthRatio float32 = float32(r.GetU8()) / 255.0
					if mob != nil {
						if mob.RealHealthRatio = healthRatio; mob.RealHealthRatio < mob.HealthRatio {
							mob.LastDamagedAt = HEALTH_DELTA_DURATION
						}
					}
				}

				if flags&0x10 != 0 {
					var statusFlags uint8 = r.GetU8()
					if mob != nil {
						mob.Damaged = statusFlags&0x01 != 0
						mob.Attacking = statusFlags&0x02 != 0
					}
				}
			}
		}

		for id = r.GetU64(); id != 0; id = r.GetU64() {
			// s.Mobs.Remove(id)

			var mob *Mob = findMob(s.Mobs, id)
			if mob != nil {
				// s.L.Infof("Remove mob: ID=%d, MobID=%d, RarityID=%d\n", mob.ID, mob.MobID, mob.RarityID)
				mob.IsDead = true
			} else {
				s.L.Warningf("Discard mob removal for unknown ID=%d\n", id)
			}
		}
	}
}
