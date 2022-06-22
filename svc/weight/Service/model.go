package Service

import "bigSystem/svc/common/entity"

// GetAllPersonsInfo

//type AllDo struct {
//	PageNum int `json:"page_num"`
//	PageSize int `json:"page_size"`
//}

//type AllDocumentsAck struct {
//	Status bool `json:"status"`
//	Res []*entity.WeightRecord `json:"res"`
//	ErrInfo string `json:"err_info"`
//}

type Page struct {
	PageSize int64 `json:"page_size"`
	PageNum  int64 `json:"page_num"`
}

type AllDocumentsAck struct {
	Status  bool             `json:"status"`
	Res     [][]*processNode `json:"res"`
	ErrInfo string           `json:"err_info"`
}

type AllDocumentsPageAck struct {
	Status  bool       `json:"status"`
	Res     resultData `json:"res"`
	ErrInfo string     `json:"err_info"`
}

type resultData struct {
	Data  [][]*processNode `json:"data"`
	Total int64            `json:"total"`
}

type processNode struct {
	Id                 string          `bson:"_id"`
	MaterialCode       string          `bson:"material_code"`
	MaterialType       string          `bson:"material_type"`
	MaterialName       string          `bson:"material_name"`
	Specifications     string          `bson:"specifications"`
	Supplier           string          `bson:"supplier"`
	Craft              string          `bson:"craft"`
	Texture            string          `bson:"texture"`
	Process            string          `bson:"process"`
	PurchaseStatus     string          `bson:"purchase_status"`
	ReceivingWarehouse string          `bson:"receiving_warehouse"`
	WeighStage         string          `bson:"weigh_stage"`
	RecordLog          []entity.Record `bson:"record_log"`
}

type AllParameterAck struct {
	Status  bool       `json:"status"`
	Res     Parameters `json:"res"`
	ErrInfo string     `json:"err_info"`
}

type Parameters struct {
	Craft          []string `json:"crafts"`
	Texture        []string `json:"texture"`
	Process        []string `json:"process"`
	PurchaseStatus []string `json:"purchase_status"`
}

type NewRecord struct {
	MaterialCode       string  `json:"material_code"`
	MaterialType       string  `json:"material_type"`
	MaterialName       string  `json:"material_name"`
	Specifications     string  `json:"specifications"`
	Supplier           string  `json:"supplier"`
	Craft              string  `json:"craft"`
	Texture            string  `json:"texture"`
	Process            string  `json:"process"`
	PurchaseStatus     string  `json:"purchase_status"`
	ReceivingWarehouse string  `json:"receiving_warehouse"`
	WeighStage         string  `json:"weigh_stage"`
	CalPerson          string  `json:"cal_person"`
	CalWeight          float64 `json:"cal_weight"`
	CalTime            string  `json:"cal_time"`
}

type NewRecordAck struct {
	Status  bool   `json:"status"`
	Res     string `json:"res"`
	ErrInfo string `json:"err_info"`
}

type UpdateFlowRecord struct {
	FlowProcess []entity.FlowProcessStage `bson:"flow_process"`
}

type MaterialCode struct {
	MaterialCode string `json:"material_code" form:"material_code" binding:"required"`
}

type WeightMaterialCodeAck struct {
	Status  bool    `json:"status"`
	Res     float64 `json:"res"`
	ErrInfo string  `json:"err_info"`
}

type Craft struct {
	Name string `json:"name"`
}

type Texture struct {
	Name string `json:"name"`
}

type Process struct {
	Name string `json:"name"`
}

type PurchaseStatus struct {
	Name string `json:"name"`
}

type CraftId struct {
	Id string `json:"id"`
}

type TextureId struct {
	Id string `json:"id"`
}

type ProcessId struct {
	Id string `json:"id"`
}

//type PurchaseStatus struct {
//	Name string `bson:"name"`
//}

type NewParameterAck struct {
	Status  bool   `json:"status"`
	Res     string `json:"res"`
	ErrInfo string `json:"err_info"`
}

type AllCraftAck struct {
	Status  bool           `json:"status"`
	Res     []entity.Craft `json:"res"`
	ErrInfo string         `json:"err_info"`
}

type AllTextureAck struct {
	Status  bool             `json:"status"`
	Res     []entity.Texture `json:"res"`
	ErrInfo string           `json:"err_info"`
}

type AllProcessAck struct {
	Status  bool             `json:"status"`
	Res     []entity.Process `json:"res"`
	ErrInfo string           `json:"err_info"`
}

type AllPurchaseStatusAck struct {
	Status  bool                    `json:"status"`
	Res     []entity.PurchaseStatus `json:"res"`
	ErrInfo string                  `json:"err_info"`
}
