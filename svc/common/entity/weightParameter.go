package entity

type WeightParameter struct {
	Craft          []string `bson:"craft"`
	Texture        []string `bson:"texture"`
	Process        []string `bson:"process"`
	PurchaseStatus []string `bson:"purchase_status"`
}

type Craft struct {
	Id   string `bson:"_id"`
	Name string `bson:"name"`
}

type Texture struct {
	Id   string `bson:"_id"`
	Name string `bson:"name"`
}

type Process struct {
	Id   string `bson:"_id"`
	Name string `bson:"name"`
}

type PurchaseStatus struct {
	Id   string `bson:"_id"`
	Name string `bson:"name"`
}
