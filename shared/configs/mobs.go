package configs

var MobRock *Mob = NewMob(MobID_Rock, "Rock", RarityID_Common, 14, 10, 25).
	SetDescription("A rock. It doesn't do much.").SetPhysics(0.1, 1.1).
	AddDrop(PetalID_Rock, 0.00522140646353364).
	AddDrop(PetalID_Fast, 0.034809377044439316).
	AddDrop(PetalID_Heavy, 0.10442813485860825).
	SetVariance(&MobVariance{Health: []float64{0.8, 1.2}, Damage: []float64{0.8, 1.2}, Radius: []float64{0.8, 1.2}})

var MobBoulder *Mob = NewMob(MobID_Boulder, "Boulder", RarityID_Unusual, 60, 10, 45).
	SetDescription("A bigger rock. It doesn't do much either.").SetPhysics(0.05, 1.2).
	AddDrop(PetalID_Basic, 0.041402433067560196).
	AddDrop(PetalID_Rock, 0.041402433067560196).
	AddDrop(PetalID_Heavy, 0.8280486464500427).
	SetVariance(&MobVariance{Health: []float64{0.85, 1.5}, Damage: []float64{0.85, 1.3}, Radius: []float64{0.8, 1.2}})

var MobLadybug *Mob = NewMob(MobID_Ladybug, "Ladybug", RarityID_Common, 25, 10, 25).
	SetDescription("Cute and harmless.").
	AddDrop(PetalID_Fast, 0.041402433067560196).
	SetVariance(&MobVariance{Health: []float64{0.9, 1.1}, Damage: []float64{0.9, 1.1}, Radius: []float64{0.9, 1.1}}).
	SetAI(&MobAI{Mode: MobAIMode_Passive, MovementMode: MobMovementMode_Jumps, IdleSpeed: 0.1})

var MobBee *Mob = NewMob(MobID_Bee, "Bee", RarityID_Common, 15, 50, 25).
	SetDescription("It stings. Don't touch it.").
	AddDrop(PetalID_Fast, 0.041402433067560196).
	SetVariance(&MobVariance{Health: []float64{0.9, 1.1}, Damage: []float64{0.9, 1.1}, Radius: []float64{0.9, 1.1}}).
	SetAI(&MobAI{Mode: MobAIMode_Passive, MovementMode: MobMovementMode_Sine, IdleSpeed: 0.15})

var MobBabyAnt *Mob = NewMob(MobID_BabyAnt, "Baby Ant", RarityID_Common, 10, 10, 12).
	SetDescription("Weak and defenseless, but big dreams.").
	SetAI(&MobAI{Mode: MobAIMode_Passive, MovementMode: MobMovementMode_Direct, IdleSpeed: 0.1})

var MobWorkerAnt *Mob = NewMob(MobID_WorkerAnt, "Worker Ant", RarityID_Common, 25, 10, 12).
	SetDescription("It's a bit temperamental, probably from working all the time.").
	SetAI(&MobAI{Mode: MobAIMode_Neutral, MovementMode: MobMovementMode_Direct, IdleSpeed: 0.1, AltSpeed: 0.92})

var MobSoldierAnt *Mob = NewMob(MobID_SoldierAnt, "Soldier Ant", RarityID_Unusual, 40, 10, 12).
	SetDescription("It's got wings and it's ready to use them.").
	SetAI(&MobAI{Mode: MobAIMode_Aggressive, MovementMode: MobMovementMode_Direct, IdleSpeed: 0.1, AltSpeed: 0.9})

var MobQueenAnt *Mob = NewMob(MobID_QueenAnt, "Queen Ant", RarityID_Rare, 200, 10, 24).
	SetDescription("You must have done something really bad if she's chasing you.").
	SetAI(&MobAI{Mode: MobAIMode_Aggressive, MovementMode: MobMovementMode_Direct, IdleSpeed: 0.1, AltSpeed: 0.875})

var MobAntHole *Mob = NewMob(MobID_AntHole, "Ant Hole", RarityID_Rare, 300, 15, 50).
	SetDescription("Ants go in, and come out. Can't explain that.").
	AddDrop(PetalID_Fast, 0.041402433067560196).SetPhysics(0, 2)
