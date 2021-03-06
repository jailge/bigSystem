package Service

import (
	"bigSystem/svc/common/db/mongodb"
	"bigSystem/svc/common/db/redisdb"
	"bigSystem/svc/common/entity"
	"context"
	"errors"
	"fmt"
	"github.com/go-kit/kit/metrics"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"

	//"go-mdm-service/utils/async"
	"go.uber.org/zap"
)

const (
	collectionRecord         = "record"
	collectionParameter      = "parameter"
	collectionCraft          = "craft"
	collectionTexture        = "texture"
	collectionProcess        = "process"
	collectionPurchaseStatus = "purchase_status"
)

type Service interface {
	//TestAdd(ctx context.Context, in Add) AddAck
	//Login(ctx context.Context, in *pb.Login) (ack *pb.LoginAck, err error)
	//Login(ctx context.Context, in Login) (ack LoginAck, err error)

	GetAllWeightRecord(ctx context.Context) (ack AllDocumentsAck, err error)
	GetWeightRecordPage(ctx context.Context, page Page) (ack AllDocumentsPageAck, err error)
	GetParameter(ctx context.Context) (ack AllParameterAck, err error)
	AddNewRecord(ctx context.Context, newRecord NewRecord) (ack NewRecordAck, err error)
	SearchWeightWithMaterialCode(ctx context.Context, materialCode MaterialCode) (ack WeightMaterialCodeAck, err error)

	GetAllCraft(ctx context.Context) (ack AllCraftAck, err error)
	GetAllTexture(ctx context.Context) (ack AllTextureAck, err error)
	GetAllProcess(ctx context.Context) (ack AllProcessAck, err error)
	GetAllPurchaseStatus(ctx context.Context) (ack AllPurchaseStatusAck, err error)
	AddCraft(ctx context.Context, craft Craft) (ack NewParameterAck, err error)
	AddTexture(ctx context.Context, texture Texture) (ack NewParameterAck, err error)
	AddProcess(ctx context.Context, process Process) (ack NewParameterAck, err error)
	AddPurchaseStatus(ctx context.Context, ps PurchaseStatus) (ack NewParameterAck, err error)
	DeleteCraftWithId(ctx context.Context, craftId string) (ack NewParameterAck, err error)
	DeleteTextureWithId(ctx context.Context, textureId string) (ack NewParameterAck, err error)
	DeleteProcessWithId(ctx context.Context, processId string) (ack NewParameterAck, err error)
	DeletePurchaseStatusWithId(ctx context.Context, psId string) (ack NewParameterAck, err error)
	UpdateCraft(ctx context.Context, id string, craft Craft) (ack NewParameterAck, err error)
	UpdateTexture(ctx context.Context, id string, texture Texture) (ack NewParameterAck, err error)
	UpdateProcess(ctx context.Context, id string, process Process) (ack NewParameterAck, err error)
	UpdatePurchaseStatus(ctx context.Context, id string, ps PurchaseStatus) (ack NewParameterAck, err error)
}

type baseServer struct {
	//*mongodb.MdbService
	logger *zap.Logger
}

var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("no document")
	NoErr              = errors.New("no errors")
	NoDocument         = errors.New("mongo: no documents in result")
	NoParameters       = errors.New("no parameters")
)

// NewService func NewService(mdb *mongodb.MdbService, log *zap.Logger) Service {
func NewService(log *zap.Logger, counter metrics.Counter, histogram metrics.Histogram) Service {
	var server Service
	//dir, _ := os.Getwd()
	//configFile := dir + "/svc/common/config/conf.db.yml"
	//"./svc/common/config/conf.db.yml"
	//if err := mdb.InitMongo("./svc/common/config/conf.db.yml"); err != nil {
	//	fmt.Println("The PersonService failed to bind with mysql")
	//}
	//server = &baseServer{mdb, log}
	server = &baseServer{log}
	server = NewMetricsMiddlewareServer(counter, histogram)(server)
	server = NewLogMiddlewareServer(log)(server)
	return server
}

// GetAllWeightRecord ????????????????????????
func (s baseServer) GetAllWeightRecord(ctx context.Context) (ack AllDocumentsAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "GetAllWeightRecord ????????????"))

	//var results []*entity.WeightRecord
	var results [][]*processNode

	// ????????????
	name, size := mongodb.NewMgo(collectionRecord).Count()
	fmt.Printf(" documents name: %+v documents size %d \n", name, size)
	cur := mongodb.NewMgo(collectionRecord).FindAll(0, size, 1)

	defer cur.Close(context.TODO())
	if cur != nil {
		fmt.Println("FindAll err:", cur)
	}
	for cur.Next(context.TODO()) {

		var elem entity.WeightRecord
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		//tmpList := make([]*processNode, 4, 6)
		var tmpList []*processNode
		for _, value := range elem.FlowProcess {
			pn := processNode{
				Id:                 elem.Id,
				MaterialCode:       elem.MaterialCode,
				MaterialType:       elem.MaterialType,
				MaterialName:       elem.MaterialName,
				Specifications:     elem.Specifications,
				Supplier:           elem.Supplier,
				Craft:              elem.Craft,
				Texture:            elem.Texture,
				Process:            elem.Process,
				PurchaseStatus:     elem.PurchaseStatus,
				ReceivingWarehouse: elem.ReceivingWarehouse,
				WeighStage:         value.WeighStage,
				RecordLog:          value.RecordLog,
			}
			tmpList = append(tmpList, &pn)
		}
		results = append(results, tmpList)
	}
	fmt.Println(results)
	if err := cur.Err(); err != nil {
		//log.Fatal(err)
		return AllDocumentsAck{
			Status:  false,
			Res:     results,
			ErrInfo: err.Error(),
		}, err
	}
	// ????????????
	for k, v := range results {
		fmt.Printf("Found  documents  %d  %v \n", k, v)
	}

	//document, err := s.MdbService.FindAllDocument(s.MdbService.Cli())
	//if err != nil {
	//	return AllDocumentsAck{
	//		Status: false,
	//		Res: document,
	//		ErrInfo: err.Error(),
	//	}, err
	//}
	ack = AllDocumentsAck{
		Status:  true,
		Res:     results,
		ErrInfo: "",
	}
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "GetAllWeightRecord ????????????"), zap.Any("???????????????", ack))
	return
}

// GetWeightRecordPage ????????????????????????
func (s baseServer) GetWeightRecordPage(ctx context.Context, page Page) (ack AllDocumentsPageAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "GetAllWeightRecord ????????????"))

	//var results []*entity.WeightRecord
	var results [][]*processNode

	// ????????????
	_, total := mongodb.NewMgo(collectionRecord).Count()
	//fmt.Printf(" documents name: %+v documents size %d \n", name, size)

	// ????????????
	//var pages int64 = int64(size / pageSize)
	var skip int64 = int64(page.PageSize * (page.PageNum - 1))

	//mongodb.NewMgo().
	cur := mongodb.NewMgo(collectionRecord).FindAll(skip, page.PageSize, 1)

	defer cur.Close(context.TODO())
	if cur != nil {
		fmt.Println("FindAll err:", cur)
	}
	for cur.Next(context.TODO()) {

		var elem entity.WeightRecord
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		//tmpList := make([]*processNode, 4, 6)
		var tmpList []*processNode
		for _, value := range elem.FlowProcess {
			pn := processNode{
				Id:                 elem.Id,
				MaterialCode:       elem.MaterialCode,
				MaterialType:       elem.MaterialType,
				MaterialName:       elem.MaterialName,
				Specifications:     elem.Specifications,
				Supplier:           elem.Supplier,
				Craft:              elem.Craft,
				Texture:            elem.Texture,
				Process:            elem.Process,
				PurchaseStatus:     elem.PurchaseStatus,
				ReceivingWarehouse: elem.ReceivingWarehouse,
				WeighStage:         value.WeighStage,
				RecordLog:          value.RecordLog,
			}
			tmpList = append(tmpList, &pn)
		}
		results = append(results, tmpList)
	}
	fmt.Println(results)

	rData := resultData{
		Data:  results,
		Total: total,
	}

	if err := cur.Err(); err != nil {
		//log.Fatal(err)
		return AllDocumentsPageAck{
			Status: false,
			Res:    rData,
			//ErrInfo: err.Error(),
		}, err
	}
	// ????????????
	//for k, v := range results {
	//	fmt.Printf("Found  documents  %d  %v \n", k, v)
	//}

	//document, err := s.MdbService.FindAllDocument(s.MdbService.Cli())
	//if err != nil {
	//	return AllDocumentsAck{
	//		Status: false,
	//		Res: document,
	//		ErrInfo: err.Error(),
	//	}, err
	//}

	ack = AllDocumentsPageAck{
		Status: true,
		Res:    rData,
		//ErrInfo: "",
	}
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "GetAllWeightRecordPage ????????????"), zap.Any("???????????????", ack))
	return
}

// GetParameter ??????????????????
func (s baseServer) GetParameter(ctx context.Context) (ack AllParameterAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "GetParameter ????????????"))

	var results Parameters

	// ????????????
	name, size := mongodb.NewMgo(collectionParameter).Count()
	fmt.Printf(" documents name: %+v documents size %d \n", name, size)
	cur := mongodb.NewMgo(collectionParameter).FindAll(0, size, 1)

	defer cur.Close(context.TODO())
	if cur != nil {
		fmt.Println("FindAll err:", cur)
	}
	for cur.Next(context.TODO()) {

		var elem entity.WeightParameter
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results.Craft = elem.Craft
		results.Texture = elem.Texture
		results.Process = elem.Process
		results.PurchaseStatus = elem.PurchaseStatus
	}

	//fmt.Println(results)
	if err := cur.Err(); err != nil {
		//log.Fatal(err)
		return AllParameterAck{
			Status:  false,
			Res:     results,
			ErrInfo: err.Error(),
		}, err
	}

	ack = AllParameterAck{
		Status:  true,
		Res:     results,
		ErrInfo: "",
	}
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "GetParameter ????????????"), zap.Any("???????????????", ack))
	return
}

// AddNewRecord ????????????
func (s baseServer) AddNewRecord(ctx context.Context, newRecord NewRecord) (ack NewRecordAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "AddNewRecord ????????????"))

	//??????????????????????????? material_weight??????
	// material_code, supplier, craft, texture, process,  ?????????
	var existsRecord entity.WeightRecord
	filter := bson.M{
		"material_code":       newRecord.MaterialCode,
		"craft":               newRecord.Craft,
		"texture":             newRecord.Texture,
		"process":             newRecord.Process,
		"purchase_status":     newRecord.PurchaseStatus,
		"receiving_warehouse": newRecord.ReceivingWarehouse,
		"supplier":            newRecord.Supplier,
	}
	cur := mongodb.NewMgo(collectionRecord).Find(filter)
	defer cur.Close(context.TODO())
	//if cur != nil {
	//	fmt.Println("FindAll err:", cur)
	//}
	for cur.Next(context.TODO()) {
		err := cur.Decode(&existsRecord)
		if err != nil {
			log.Fatal(err)
			//utils.GetLogger().Debug(fmt.Sprintf("Get existRecord err: %s", err.Error()))
		}
	}
	//err = mongodb.NewMgo(collectionRecord).FindOne("material_code", newRecord.MaterialCode).Decode(&existsRecord)

	//????????????????????????
	if cur != nil {
		//???????????????material_code??????
		//??????????????????
		//oldFlowProcess := existsRecord.FlowProcess
		record := entity.Record{
			CalPerson: newRecord.CalPerson,
			CalWeight: newRecord.CalWeight,
			CalTime:   newRecord.CalTime,
		}
		same := false
		for i, value := range existsRecord.FlowProcess {
			//??????stage??????
			if value.WeighStage == newRecord.WeighStage {
				//fmt.Println("value")
				//fmt.Println(value)
				value.RecordLog = append(value.RecordLog, entity.Record{
					CalPerson: record.CalPerson,
					CalWeight: record.CalWeight,
					CalTime:   record.CalTime,
				})
				//fmt.Println(value.RecordLog)
				//fmt.Println(existsRecord.FlowProcess)
				existsRecord.FlowProcess[i] = value
				//fmt.Println(existsRecord.FlowProcess)
				same = true
				break
			}
		}
		//fmt.Println(same)
		if same == false {
			existsRecord.FlowProcess = append(existsRecord.FlowProcess, entity.FlowProcessStage{
				WeighStage: newRecord.WeighStage,
				RecordLog:  []entity.Record{record},
			})
		}
		//fmt.Println(existsRecord.FlowProcess)
		//fmt.Println(&existsRecord.FlowProcess)

		//update := bson.D{
		//	{
		//		"$set", bson.D{
		//			{"flow_process", existsRecord.FlowProcess},
		//		}},
		//}
		m := bson.M{
			"$set": UpdateFlowRecord{FlowProcess: existsRecord.FlowProcess},
		}

		//fmt.Println(update)
		//fmt.Println(m)

		updateResult := mongodb.NewMgo(collectionRecord).UpdateMany("_id", existsRecord.Id, m)
		//fmt.Println("########", updateResult)

		ack = NewRecordAck{
			Status:  true,
			Res:     fmt.Sprintf("MatchedCount:%d,  UpsertedCount:%d, ModifiedCount:%d", updateResult.MatchedCount, updateResult.UpsertedCount, updateResult.ModifiedCount),
			ErrInfo: "",
		}

	} else if cur == nil {
		//????????????????????????
		record := entity.Record{
			CalPerson: newRecord.CalPerson,
			CalWeight: newRecord.CalWeight,
			CalTime:   newRecord.CalTime,
		}
		processStage := entity.FlowProcessStage{
			WeighStage: newRecord.WeighStage,
			RecordLog:  []entity.Record{record},
		}
		weightRecord := entity.NewWeightRecord{
			MaterialCode:       newRecord.MaterialCode,
			MaterialType:       newRecord.MaterialName,
			MaterialName:       newRecord.MaterialName,
			Specifications:     newRecord.Specifications,
			Supplier:           newRecord.Supplier,
			Craft:              newRecord.Craft,
			Texture:            newRecord.Texture,
			Process:            newRecord.Process,
			PurchaseStatus:     newRecord.PurchaseStatus,
			ReceivingWarehouse: newRecord.ReceivingWarehouse,
			FlowProcess:        []entity.FlowProcessStage{processStage},
		}
		insertOneResult := mongodb.NewMgo(collectionRecord).InsertOne(weightRecord)
		//fmt.Println(insertOneResult)
		ack = NewRecordAck{
			Status:  true,
			Res:     fmt.Sprintf("%s", insertOneResult.InsertedID),
			ErrInfo: "",
		}
	} else {
		// ??????????????????
		//fmt.Println(err)
		ack = NewRecordAck{}
	}
	//fmt.Println(existsRecord)
	//return
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "AddNewRecord ????????????"), zap.Any("???????????????", ack))
	return
}

func (s baseServer) SearchWeightWithMaterialCode(ctx context.Context, materialCode MaterialCode) (ack WeightMaterialCodeAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "SearchWeightWithMaterialCode ????????????"))

	//??????????????????material_code
	var existsRecord entity.WeightRecord
	err = mongodb.NewMgo(collectionRecord).FindOne("material_code", materialCode.MaterialCode).Decode(&existsRecord)

	//????????????????????????
	if err == nil {
		for _, value := range existsRecord.FlowProcess {
			if value.WeighStage == "????????????" {
				return WeightMaterialCodeAck{
					Status:  true,
					Res:     value.RecordLog[len(value.RecordLog)-1].CalWeight,
					ErrInfo: "",
				}, err
			}
		}
		ack = WeightMaterialCodeAck{
			Status:  false,
			Res:     0,
			ErrInfo: "",
		}

	} else if err.Error() == NoDocument.Error() {
		err = errors.New("no document")
		ack = WeightMaterialCodeAck{
			Status:  false,
			Res:     0,
			ErrInfo: "",
		}
	} else {
		// ??????????????????
		//fmt.Println(err)
		ack = WeightMaterialCodeAck{}
	}
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "SearchWeightWithMaterialCode ????????????"), zap.Any("???????????????", ack))
	return
}

// ****************************************************
// ****************************************************
// ????????????
// ****************************************************
// ****************************************************

// GetAllCraft ????????????
func (s baseServer) GetAllCraft(ctx context.Context) (ack AllCraftAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "GetAllCraft ????????????"))

	var results []entity.Craft

	// ????????????
	_, size := mongodb.NewMgo(collectionCraft).Count()
	//fmt.Printf(" documents name: %+v documents size %d \n", name, size)
	cur := mongodb.NewMgo(collectionCraft).FindAll(0, size, 1)

	defer cur.Close(context.TODO())
	if cur != nil {
		fmt.Println("FindAll :", cur)
		//err = errors.New("FindAll err")
	}
	for cur.Next(context.TODO()) {
		var elem entity.Craft
		err := cur.Decode(&elem)
		if err != nil {
			//err = errors.New("FindAll err")
			log.Fatal(err)
		}
		results = append(results, elem)
	}
	if err := cur.Err(); err != nil {
		//log.Fatal(err)
		return AllCraftAck{
			Status:  false,
			Res:     results,
			ErrInfo: err.Error(),
		}, err
	}
	ack = AllCraftAck{
		Status:  true,
		Res:     results,
		ErrInfo: "",
	}
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "GetAllCraft ????????????"), zap.Any("???????????????", ack))
	return
}

// AddCraft ????????????
func (s baseServer) AddCraft(ctx context.Context, craft Craft) (ack NewParameterAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "AddCraft ????????????"))

	if craft.Name == "" {
		err = NoParameters
		return NewParameterAck{}, err
	}

	//??????????????????craft
	var existCraft entity.Craft
	err = mongodb.NewMgo(collectionCraft).FindOne("name", craft.Name).Decode(&existCraft)
	fmt.Println(err)
	//????????????????????????
	if err == nil {
		// ?????????
		ack = NewParameterAck{
			Status:  false,
			Res:     fmt.Sprintf("??????????????????%s", craft.Name),
			ErrInfo: "",
		}
	} else if err.Error() == NoDocument.Error() {
		// ????????????????????????

		insertOneResult := mongodb.NewMgo(collectionCraft).InsertOne(craft)
		err = nil
		//fmt.Println(insertOneResult)
		ack = NewParameterAck{
			Status:  true,
			Res:     fmt.Sprintf("%s", insertOneResult.InsertedID),
			ErrInfo: "",
		}
	} else {
		// ??????????????????
		//fmt.Println("***********")
		//fmt.Println(err)
		ack = NewParameterAck{}
	}
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "AddCraft ????????????"), zap.Any("???????????????", ack))
	return
}

// DeleteCraftWithId ????????????
func (s baseServer) DeleteCraftWithId(ctx context.Context, craftId string) (ack NewParameterAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "DeleteCraftWithId ????????????"))

	objId, _ := primitive.ObjectIDFromHex(craftId)
	//??????????????????craft
	var existCraft entity.Craft
	err = mongodb.NewMgo(collectionCraft).FindOne("_id", objId).Decode(&existCraft)
	//fmt.Println(craftId)
	//fmt.Println(err)
	//????????????????????????
	if err == nil {
		// ?????????????????????
		deleteResult := mongodb.NewMgo(collectionCraft).Delete("_id", objId)
		err = nil
		ack = NewParameterAck{
			Status:  true,
			Res:     fmt.Sprintf("??????????????????%d", deleteResult),
			ErrInfo: "",
		}
	} else {
		// ??????????????????
		//fmt.Println(err)
		ack = NewParameterAck{}
	}
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "DeleteCraftWithId ????????????"), zap.Any("???????????????", ack))
	return
}

// UpdateCraft ????????????
func (s baseServer) UpdateCraft(ctx context.Context, id string, craft Craft) (ack NewParameterAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "UpdateCraft ????????????"))

	objId, _ := primitive.ObjectIDFromHex(id)

	//??????????????????craft
	var existCraft entity.Craft
	err = mongodb.NewMgo(collectionCraft).FindOne("_id", objId).Decode(&existCraft)

	//????????????????????????
	if err == nil {
		// ?????????,?????????
		// redis ????????????
		//res, err := redisdb.RdsClient.SetNX(clientId, existCraft.Id, time.Minute).Result()
		if r, _ := redisdb.GetLock(existCraft.Id, craft.ClientId); r {
			//fmt.Println("***********************")
			//fmt.Println(r)
			// ????????????
			update := bson.D{
				{
					"$set", bson.D{
						{"name", craft.Name},
					}},
			}
			time.Sleep(30 * time.Second)
			updateResult := mongodb.NewMgo(collectionCraft).UpdateOne("_id", objId, update)

			err = nil
			_, _ = redisdb.RunEvalDel(existCraft.Id, craft.ClientId)
			ack = NewParameterAck{
				Status:  true,
				Res:     fmt.Sprintf("Matched %d documents and updated %d documents.", updateResult.MatchedCount, updateResult.ModifiedCount),
				ErrInfo: "",
			}
		} else {
			err = nil
			ms, _ := redisdb.Pttl(existCraft.Id)
			ack = NewParameterAck{
				Status:  false,
				Res:     fmt.Sprintf("%d ms????????????", ms),
				ErrInfo: "",
			}
		}

	} else if err.Error() == NoDocument.Error() {
		// ?????????
		//err = errors.New("no errors")
		err = errors.New("no document")
		ack = NewParameterAck{}
	} else {
		// ??????????????????
		//fmt.Println(err)
		ack = NewParameterAck{}
	}
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "UpdateCraft ????????????"), zap.Any("???????????????", ack))
	return
}

// GetAllTexture ????????????
func (s baseServer) GetAllTexture(ctx context.Context) (ack AllTextureAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "GetAllTexture ????????????"))

	var results []entity.Texture

	// ????????????
	_, size := mongodb.NewMgo(collectionTexture).Count()
	//fmt.Printf(" documents name: %+v documents size %d \n", name, size)
	cur := mongodb.NewMgo(collectionTexture).FindAll(0, size, 1)

	defer cur.Close(context.TODO())
	if cur != nil {
		fmt.Println("FindAll :", cur)
		//err = errors.New("FindAll err")
	}
	for cur.Next(context.TODO()) {
		var elem entity.Texture
		err := cur.Decode(&elem)
		if err != nil {
			//err = errors.New("FindAll err")
			log.Fatal(err)
		}
		results = append(results, elem)
	}
	if err := cur.Err(); err != nil {
		//log.Fatal(err)
		return AllTextureAck{
			Status:  false,
			Res:     results,
			ErrInfo: err.Error(),
		}, err
	}
	ack = AllTextureAck{
		Status:  true,
		Res:     results,
		ErrInfo: "",
	}
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "GetAllTexture ????????????"), zap.Any("???????????????", ack))
	return
}

// AddTexture ????????????
func (s baseServer) AddTexture(ctx context.Context, texture Texture) (ack NewParameterAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "AddTexture ????????????"))

	if texture.Name == "" {
		err = NoParameters
		return NewParameterAck{}, err
	}

	//??????????????????texture
	var existTexture entity.Texture
	err = mongodb.NewMgo(collectionTexture).FindOne("name", texture.Name).Decode(&existTexture)
	//fmt.Println(err)
	//????????????????????????
	if err == nil {
		// ?????????
		ack = NewParameterAck{
			Status:  false,
			Res:     fmt.Sprintf("??????????????????%s", texture.Name),
			ErrInfo: "",
		}
	} else if err.Error() == NoDocument.Error() {
		// ????????????????????????
		//err = errors.New("no errors")
		insertOneResult := mongodb.NewMgo(collectionTexture).InsertOne(texture)
		err = nil
		//fmt.Println(insertOneResult)
		ack = NewParameterAck{
			Status:  true,
			Res:     fmt.Sprintf("%s", insertOneResult.InsertedID),
			ErrInfo: "",
		}
	} else {
		// ??????????????????
		//fmt.Println(err)
		ack = NewParameterAck{}
	}
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "AddTexture ????????????"), zap.Any("???????????????", ack))
	return
}

// DeleteTextureWithId ????????????
func (s baseServer) DeleteTextureWithId(ctx context.Context, textureId string) (ack NewParameterAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "DeleteTextureWithId ????????????"))

	objId, _ := primitive.ObjectIDFromHex(textureId)
	//??????????????????texture
	var existTexture entity.Texture
	err = mongodb.NewMgo(collectionTexture).FindOne("_id", objId).Decode(&existTexture)
	//fmt.Println(craftId)
	fmt.Println(err)
	//????????????????????????
	if err == nil {
		// ?????????????????????
		deleteResult := mongodb.NewMgo(collectionTexture).Delete("_id", objId)
		err = nil
		ack = NewParameterAck{
			Status:  true,
			Res:     fmt.Sprintf("??????????????????%d", deleteResult),
			ErrInfo: "",
		}
	} else {
		// ??????????????????
		//fmt.Println(err)
		ack = NewParameterAck{}
	}
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "DeleteTextureWithId ????????????"), zap.Any("???????????????", ack))
	return
}

// UpdateTexture ????????????
func (s baseServer) UpdateTexture(ctx context.Context, id string, texture Texture) (ack NewParameterAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "UpdateTexture ????????????"))

	objId, _ := primitive.ObjectIDFromHex(id)

	//??????????????????texture
	var exist entity.Texture
	err = mongodb.NewMgo(collectionTexture).FindOne("_id", objId).Decode(&exist)
	//????????????????????????
	if err == nil {
		// ?????????,?????????
		// redis ????????????
		//res, err := redisdb.RdsClient.SetNX(clientId, existCraft.Id, time.Minute).Result()
		if r, _ := redisdb.GetLock(exist.Id, texture.ClientId); r {
			update := bson.D{
				{
					"$set", bson.D{
						{"name", texture.Name},
					}},
			}
			updateResult := mongodb.NewMgo(collectionTexture).UpdateOne("_id", objId, update)
			err = nil
			_, _ = redisdb.RunEvalDel(exist.Id, texture.ClientId)
			ack = NewParameterAck{
				Status:  true,
				Res:     fmt.Sprintf("Matched %d documents and updated %d documents.", updateResult.MatchedCount, updateResult.ModifiedCount),
				ErrInfo: "",
			}
		} else {
			err = nil
			ms, _ := redisdb.Pttl(exist.Id)
			ack = NewParameterAck{
				Status:  false,
				Res:     fmt.Sprintf("%d ms????????????", ms),
				ErrInfo: "",
			}
		}

	} else if err.Error() == NoDocument.Error() {
		// ????????????????????????
		//err = errors.New("no errors")
		err = errors.New("no document")
		ack = NewParameterAck{}
	} else {
		// ??????????????????
		//fmt.Println(err)
		ack = NewParameterAck{}
	}
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "UpdateTexture ????????????"), zap.Any("???????????????", ack))
	return
}

// GetAllProcess ????????????
func (s baseServer) GetAllProcess(ctx context.Context) (ack AllProcessAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "GetAllProcess ????????????"))

	var results []entity.Process

	// ????????????
	_, size := mongodb.NewMgo(collectionProcess).Count()
	//fmt.Printf(" documents name: %+v documents size %d \n", name, size)
	cur := mongodb.NewMgo(collectionProcess).FindAll(0, size, 1)

	defer cur.Close(context.TODO())
	if cur != nil {
		fmt.Println("FindAll :", cur)
		//err = errors.New("FindAll err")
	}
	for cur.Next(context.TODO()) {
		var elem entity.Process
		err := cur.Decode(&elem)
		if err != nil {
			//err = errors.New("FindAll err")
			log.Fatal(err)
		}
		results = append(results, elem)
	}
	if err := cur.Err(); err != nil {
		//log.Fatal(err)
		return AllProcessAck{
			Status:  false,
			Res:     results,
			ErrInfo: err.Error(),
		}, err
	}
	ack = AllProcessAck{
		Status:  true,
		Res:     results,
		ErrInfo: "",
	}
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "GetAllProcess ????????????"), zap.Any("???????????????", ack))
	return
}

// AddProcess ????????????
func (s baseServer) AddProcess(ctx context.Context, process Process) (ack NewParameterAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "AddProcess ????????????"))

	if process.Name == "" {
		err = NoParameters
		return NewParameterAck{}, err
	}

	//??????????????????process
	var existProcess entity.Process
	err = mongodb.NewMgo(collectionProcess).FindOne("name", process.Name).Decode(&existProcess)
	//fmt.Println(err)
	//????????????????????????
	if err == nil {
		// ?????????
		ack = NewParameterAck{
			Status:  false,
			Res:     fmt.Sprintf("??????????????????%s", process.Name),
			ErrInfo: "",
		}
	} else if err.Error() == NoDocument.Error() {
		// ????????????????????????
		//err = errors.New("no errors")
		insertOneResult := mongodb.NewMgo(collectionProcess).InsertOne(process)
		err = nil
		//fmt.Println(insertOneResult)
		ack = NewParameterAck{
			Status:  true,
			Res:     fmt.Sprintf("%s", insertOneResult.InsertedID),
			ErrInfo: "",
		}
	} else {
		// ??????????????????
		//fmt.Println(err)
		ack = NewParameterAck{}
	}
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "AddProcess ????????????"), zap.Any("???????????????", ack))
	return
}

// DeleteProcessWithId ????????????
func (s baseServer) DeleteProcessWithId(ctx context.Context, processId string) (ack NewParameterAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "DeleteProcessWithId ????????????"))

	objId, _ := primitive.ObjectIDFromHex(processId)
	//??????????????????Process
	var existProcess entity.Process
	err = mongodb.NewMgo(collectionProcess).FindOne("_id", objId).Decode(&existProcess)
	//fmt.Println(craftId)
	//????????????????????????
	if err == nil {
		// ?????????????????????
		deleteResult := mongodb.NewMgo(collectionProcess).Delete("_id", objId)
		err = nil
		ack = NewParameterAck{
			Status:  true,
			Res:     fmt.Sprintf("??????????????????%d", deleteResult),
			ErrInfo: "",
		}
	} else {
		// ??????????????????
		//fmt.Println(err)
		ack = NewParameterAck{}
	}
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "DeleteProcessWithId ????????????"), zap.Any("???????????????", ack))
	return
}

// UpdateProcess ????????????
func (s baseServer) UpdateProcess(ctx context.Context, id string, process Process) (ack NewParameterAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "UpdateProcess ????????????"))

	objId, _ := primitive.ObjectIDFromHex(id)

	//??????????????????texture
	var exist entity.Process
	err = mongodb.NewMgo(collectionProcess).FindOne("_id", objId).Decode(&exist)
	//????????????????????????
	if err == nil {
		// ?????????,?????????
		// redis ????????????
		//res, err := redisdb.RdsClient.SetNX(clientId, existCraft.Id, time.Minute).Result()
		if r, _ := redisdb.GetLock(exist.Id, process.ClientId); r {
			update := bson.D{
				{
					"$set", bson.D{
						{"name", process.Name},
					}},
			}
			updateResult := mongodb.NewMgo(collectionProcess).UpdateOne("_id", objId, update)
			err = nil
			_, _ = redisdb.RunEvalDel(exist.Id, process.ClientId)
			ack = NewParameterAck{
				Status:  true,
				Res:     fmt.Sprintf("Matched %d documents and updated %d documents.", updateResult.MatchedCount, updateResult.ModifiedCount),
				ErrInfo: "",
			}
		} else {
			err = nil
			ms, _ := redisdb.Pttl(exist.Id)
			ack = NewParameterAck{
				Status:  false,
				Res:     fmt.Sprintf("%d ms????????????", ms),
				ErrInfo: "",
			}
		}

	} else if err.Error() == NoDocument.Error() {
		// ????????????????????????
		//err = errors.New("no errors")
		err = errors.New("no document")
		ack = NewParameterAck{}
	} else {
		// ??????????????????
		//fmt.Println(err)
		ack = NewParameterAck{}
	}
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "UpdateProcess ????????????"), zap.Any("???????????????", ack))
	return
}

// GetAllPurchaseStatus ??????????????????
func (s baseServer) GetAllPurchaseStatus(ctx context.Context) (ack AllPurchaseStatusAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "GetAllPurchaseStatus ????????????"))

	var results []entity.PurchaseStatus

	// ????????????
	_, size := mongodb.NewMgo(collectionPurchaseStatus).Count()
	//fmt.Printf(" documents name: %+v documents size %d \n", name, size)
	cur := mongodb.NewMgo(collectionPurchaseStatus).FindAll(0, size, 1)

	defer cur.Close(context.TODO())
	if cur != nil {
		fmt.Println("FindAll :", cur)
		//err = errors.New("FindAll err")
	}
	for cur.Next(context.TODO()) {
		var elem entity.PurchaseStatus
		err := cur.Decode(&elem)
		if err != nil {
			//err = errors.New("FindAll err")
			log.Fatal(err)
		}
		results = append(results, elem)
	}
	if err := cur.Err(); err != nil {
		//log.Fatal(err)
		ack = AllPurchaseStatusAck{
			Status:  false,
			Res:     results,
			ErrInfo: err.Error(),
		}
	} else {
		ack = AllPurchaseStatusAck{
			Status:  true,
			Res:     results,
			ErrInfo: "",
		}
	}

	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "GetAllPurchaseStatus ????????????"), zap.Any("???????????????", ack))
	return
}

// AddPurchaseStatus ??????????????????
func (s baseServer) AddPurchaseStatus(ctx context.Context, ps PurchaseStatus) (ack NewParameterAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "AddPurchaseStatus ????????????"))

	if ps.Name == "" {
		err = NoParameters
		return NewParameterAck{}, err
	}

	//??????????????????purchase_status
	var existPurchaseStatus entity.PurchaseStatus
	err = mongodb.NewMgo(collectionPurchaseStatus).FindOne("name", ps.Name).Decode(&existPurchaseStatus)
	//fmt.Println(err)
	//????????????????????????
	if err == nil {
		// ?????????
		ack = NewParameterAck{
			Status:  false,
			Res:     fmt.Sprintf("????????????????????????%s", ps.Name),
			ErrInfo: "",
		}
	} else if err.Error() == NoDocument.Error() {
		// ????????????????????????
		//err = errors.New("no errors")
		insertOneResult := mongodb.NewMgo(collectionPurchaseStatus).InsertOne(ps)
		err = nil
		//fmt.Println(insertOneResult)
		ack = NewParameterAck{
			Status:  true,
			Res:     fmt.Sprintf("%s", insertOneResult.InsertedID),
			ErrInfo: "",
		}
	} else {
		// ??????????????????
		//fmt.Println(err)
		ack = NewParameterAck{}
	}
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "AddPurchaseStatus ????????????"), zap.Any("???????????????", ack))
	return
}

// DeletePurchaseStatusWithId ??????????????????
func (s baseServer) DeletePurchaseStatusWithId(ctx context.Context, psId string) (ack NewParameterAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "DeletePurchaseStatusWithId ????????????"))

	objId, _ := primitive.ObjectIDFromHex(psId)
	//??????????????????purchase_status
	var existPurchaseStatus entity.PurchaseStatus
	err = mongodb.NewMgo(collectionPurchaseStatus).FindOne("_id", objId).Decode(&existPurchaseStatus)

	//????????????????????????
	if err == nil {
		// ?????????????????????
		deleteResult := mongodb.NewMgo(collectionPurchaseStatus).Delete("_id", objId)
		err = nil
		ack = NewParameterAck{
			Status:  true,
			Res:     fmt.Sprintf("????????????????????????%d", deleteResult),
			ErrInfo: "",
		}
	} else {
		// ??????????????????
		//fmt.Println(err)
		ack = NewParameterAck{}
	}

	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "DeletePurchaseStatusWithId ????????????"), zap.Any("???????????????", ack))
	return
}

// UpdatePurchaseStatus ??????????????????
func (s baseServer) UpdatePurchaseStatus(ctx context.Context, id string, ps PurchaseStatus) (ack NewParameterAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "UpdatePurchaseStatus ????????????"))

	objId, _ := primitive.ObjectIDFromHex(id)

	//??????????????????texture
	var exist entity.PurchaseStatus
	err = mongodb.NewMgo(collectionPurchaseStatus).FindOne("_id", objId).Decode(&exist)
	//????????????????????????
	if err == nil {
		// ?????????,?????????
		// redis ????????????
		//res, err := redisdb.RdsClient.SetNX(clientId, existCraft.Id, time.Minute).Result()
		if r, _ := redisdb.GetLock(exist.Id, ps.ClientId); r {
			update := bson.D{
				{
					"$set", bson.D{
						{"name", ps.Name},
					}},
			}
			updateResult := mongodb.NewMgo(collectionPurchaseStatus).UpdateOne("_id", objId, update)
			err = nil
			_, _ = redisdb.RunEvalDel(exist.Id, ps.ClientId)
			ack = NewParameterAck{
				Status:  true,
				Res:     fmt.Sprintf("Matched %d documents and updated %d documents.", updateResult.MatchedCount, updateResult.ModifiedCount),
				ErrInfo: "",
			}
		} else {
			err = nil
			ms, _ := redisdb.Pttl(exist.Id)
			ack = NewParameterAck{
				Status:  false,
				Res:     fmt.Sprintf("%d ms????????????", ms),
				ErrInfo: "",
			}
		}

	} else if err.Error() == NoDocument.Error() {
		// ????????????
		//err = errors.New("no errors")
		err = errors.New("no document")
		ack = NewParameterAck{}
	} else {
		// ??????????????????
		//fmt.Println(err)
		ack = NewParameterAck{}
	}
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("?????? Service", "UpdatePurchaseStatus ????????????"), zap.Any("???????????????", ack))
	return
}
