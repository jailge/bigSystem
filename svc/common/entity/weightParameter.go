package entity

type WeightParameter struct {
	Craft          []string `bson:"craft"`
	Texture        []string `bson:"texture"`
	Process        []string `bson:"process"`
	PurchaseStatus []string `bson:"purchase_status"`
}
