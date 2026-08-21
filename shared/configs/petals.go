package configs

var PetalBasic *Petal = NewPetal(PetalID_Basic, "Basic", RarityID_Common, 1000).
	SetDescription("A nice petal, not too strong but not too weak").
	SetHealthAndDamage(10, 10)

var PetalFast *Petal = NewPetal(PetalID_Fast, "Fast", RarityID_Common, 500).
	SetDescription("Weaker than most petals, but recharges very quickly").
	SetHealthAndDamage(5, 8).SetRadius(7).SetOrbit(DefaultPetalOrbit.Neutral, 145, 40)

var PetalHeavy *Petal = NewPetal(PetalID_Heavy, "Heavy", RarityID_Common, 5500).
	SetDescription("Very resilient and deals more damage, but recharges very slowly").
	SetHealthAndDamage(20, 20).SetRadius(12).SetOrbit(DefaultPetalOrbit.Neutral, 120, 50)

var PetalRock *Petal = NewPetal(PetalID_Rock, "Rock", RarityID_Common, 10000).
	SetDescription("Extremely durable, but takes a long time to recharge").
	SetHealthAndDamage(90, 10).SetRadius(14).SetOrbit(DefaultPetalOrbit.Neutral, 110, 55)

var PetalSquare *Petal = NewPetal(PetalID_Square, "Square", RarityID_Unique, 1000).
	SetDescription("This shape... It looks familiar...").
	SetHealthAndDamage(10, 10)

var PetalRose *Petal = NewPetal(PetalID_Rose, "Rose", RarityID_Unusual, 3500).
	SetDescription("Its healing properties are amazing. Not so good at combat though").
	SetHealthAndDamage(5, 5).SetHeal(10, 1000).SetOrbit(DefaultPetalOrbit.Neutral, DefaultPetalOrbit.Neutral, DefaultPetalOrbit.Defend)

var PetalStinger *Petal = NewPetal(PetalID_Stinger, "Stinger", RarityID_Unusual, 4000).
	SetDescription("It really hurts, but it's very fragile").
	SetHealthAndDamage(8, 35).SetRadius(8)

var PetalIris *Petal = NewPetal(PetalID_Iris, "Iris", RarityID_Rare, 6000).
	SetDescription("Very poisonous, but takes a little while to do its work").
	SetHealthAndDamage(5, 5).SetRadius(7).SetPoison(60, 6000)

///////

var PetalFaster *Petal = NewPetal(PetalID_Faster, "Faster", RarityID_Rare, 500).
	SetDescription("It's so light it makes your other petals spin faster").
	SetHealthAndDamage(5, 8).SetRadius(8).SetExtraRotation(0.8)

var PetalTwin *Petal = NewPetal(PetalID_Twin, "Twin", RarityID_Unusual, 1000).
	SetDescription("Why stop at one? Why not TWO?!").
	SetHealthAndDamage(5, 8).SetMulti(2, 1).SetRadius(7)

// Test petal!
var PetalTriplet *Petal = NewPetal(PetalID_Triplet, "Triplet", RarityID_Epic, 1000).
	SetDescription("How about THREE?!").
	SetHealthAndDamage(5, 8).SetMulti(3, 1).SetRadius(7)

var PetalLeaf *Petal = NewPetal(PetalID_Leaf, "Leaf", RarityID_Common, 1000).
	SetDescription("A very light petal, but not very strong").
	SetHealthAndDamage(10, 8).SetRadius(11).SetConstantHeal(1)
