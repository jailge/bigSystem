package entity

type WeightRecord struct {
	Id                 string             `bson:"_id"`
	MaterialCode       string             `bson:"material_code"`
	MaterialType       string             `bson:"material_type"`
	MaterialName       string             `bson:"material_name"`
	Specifications     string             `bson:"specifications"`
	Supplier           string             `bson:"supplier"`
	Craft              string             `bson:"craft"`
	Texture            string             `bson:"texture"`
	Process            string             `bson:"process"`
	PurchaseStatus     string             `bson:"purchase_status"`
	ReceivingWarehouse string             `bson:"receiving_warehouse"`
	FlowProcess        []FlowProcessStage `bson:"flow_process"`

	//Bias string `bson:"bias"`
	//AccordingTo string `bson:"according_to"`
	//UnitPrice float64 `bson:"unit_price"`
	//TotalPrice float64 `bson:"total_price"`
}

type FlowProcessStage struct {
	WeighStage string   `bson:"weigh_stage"`
	RecordLog  []Record `bson:"record_log"`
}

type Record struct {
	CalPerson string  `bson:"cal_person"`
	CalWeight float64 `bson:"cal_weight"`
	CalTime   string  `bson:"cal_time"`
}

type NewWeightRecord struct {
	MaterialCode       string             `bson:"material_code"`
	MaterialType       string             `bson:"material_type"`
	MaterialName       string             `bson:"material_name"`
	Specifications     string             `bson:"specifications"`
	Supplier           string             `bson:"supplier"`
	Craft              string             `bson:"craft"`
	Texture            string             `bson:"texture"`
	Process            string             `bson:"process"`
	PurchaseStatus     string             `bson:"purchase_status"`
	ReceivingWarehouse string             `bson:"receiving_warehouse"`
	FlowProcess        []FlowProcessStage `bson:"flow_process"`
}
