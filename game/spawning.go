package game

import (
	"math/rand/v2"

	"github.com/z46-dev/florr-legacy/shared/configs"
)

func (g *Game) randomPosition() (x, y float64) {
	x, y = rand.Float64()*g.Width-g.Width*0.5, rand.Float64()*g.Height-g.Height*0.5
	return
}

func (g *Game) SpawnMobs() {
	if len(g.mobs.Items) >= 16384 {
		return
	}

	var mobChoices []*configs.Mob = []*configs.Mob{
		configs.MobRock,
		configs.MobBoulder,
		configs.MobLadybug,
		configs.MobBee,
		configs.MobBabyAnt,
		configs.MobWorkerAnt,
		configs.MobSoldierAnt,
		configs.MobQueenAnt,
		configs.MobAntHole,
	}

	var mob *Mob = NewMob(g, mobChoices[rand.IntN(len(mobChoices))])
	mob.Position.X, mob.Position.Y = g.randomPosition()
	mob.TeamID = TeamID_Mobs
}
