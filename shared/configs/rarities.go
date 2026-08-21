package configs

var (
	RarityCommon    *Rarity = NewRarity(RarityID_Common, "Common", "#87f05f")
	RarityUnusual   *Rarity = NewRarity(RarityID_Unusual, "Unusual", "#fbe647")
	RarityRare      *Rarity = NewRarity(RarityID_Rare, "Rare", "#654de9")
	RarityEpic      *Rarity = NewRarity(RarityID_Epic, "Epic", "#891cde")
	RarityLegendary *Rarity = NewRarity(RarityID_Legendary, "Legendary", "#e01d26")
	RarityUnique    *Rarity = NewRarity(RarityID_Unique, "Unique", "#da2267") // #555555
	// RarityMythic    *Rarity = NewRarity(RarityID_Mythic, "Mythic", "#4ddcdf")
	// RairtyUltra     *Rarity = NewRarity(RarityID_Ultra, "Ultra", "#f62175")
	// RaritySuper     *Rarity = NewRarity(RarityID_Super, "Super", "#55ff9d")
)
