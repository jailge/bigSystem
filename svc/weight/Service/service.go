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

// GetAllWeightRecord 获取所有称重记录
func (s baseServer) GetAllWeightRecord(ctx context.Context) (ack AllDocumentsAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "GetAllWeightRecord 处理请求"))

	//var results []*entity.WeightRecord
	var results [][]*processNode

	// 查询总数
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
	// 遍历结果
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
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "GetAllWeightRecord 处理请求"), zap.Any("处理返回值", ack))
	return
}

// GetWeightRecordPage 分页显示称重记录
func (s baseServer) GetWeightRecordPage(ctx context.Context, page Page) (ack AllDocumentsPageAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "GetAllWeightRecord 处理请求"))

	//var results []*entity.WeightRecord
	var results [][]*processNode

	// 查询总数
	_, total := mongodb.NewMgo(collectionRecord).Count()
	//fmt.Printf(" documents name: %+v documents size %d \n", name, size)

	// 分页逻辑
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
	// 遍历结果
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
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "GetAllWeightRecordPage 处理请求"), zap.Any("处理返回值", ack))
	return
}

// GetParameter 获取参数信息
func (s baseServer) GetParameter(ctx context.Context) (ack AllParameterAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "GetParameter 处理请求"))

	var results Parameters

	// 查询总数
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
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "GetParameter 处理请求"), zap.Any("处理返回值", ack))
	return
}

// AddNewRecord 新增记录
func (s baseServer) AddNewRecord(ctx context.Context, newRecord NewRecord) (ack NewRecordAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "AddNewRecord 处理请求"))

	//查询是否存在相同的 material_weight记录
	// material_code, supplier, craft, texture, process,  要相同
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

	//判断是否存在记录
	if cur != nil {
		//如果已存在material_code记录
		//新增称重记录
		//oldFlowProcess := existsRecord.FlowProcess
		record := entity.Record{
			CalPerson: newRecord.CalPerson,
			CalWeight: newRecord.CalWeight,
			CalTime:   newRecord.CalTime,
		}
		same := false
		for i, value := range existsRecord.FlowProcess {
			//如果stage相同
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
		//如果不存在就新增
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
		// 如果查询错误
		//fmt.Println(err)
		ack = NewRecordAck{}
	}
	//fmt.Println(existsRecord)
	//return
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "AddNewRecord 处理请求"), zap.Any("处理返回值", ack))
	return
}

func (s baseServer) SearchWeightWithMaterialCode(ctx context.Context, materialCode MaterialCode) (ack WeightMaterialCodeAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "SearchWeightWithMaterialCode 处理请求"))

	//查询是否存在material_code
	var existsRecord entity.WeightRecord
	err = mongodb.NewMgo(collectionRecord).FindOne("material_code", materialCode.MaterialCode).Decode(&existsRecord)

	//判断是否存在记录
	if err == nil {
		for _, value := range existsRecord.FlowProcess {
			if value.WeighStage == "物料首称" {
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
		// 如果查询错误
		//fmt.Println(err)
		ack = WeightMaterialCodeAck{}
	}
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "SearchWeightWithMaterialCode 处理请求"), zap.Any("处理返回值", ack))
	return
}

// ****************************************************
// ****************************************************
// 参数维护
// ****************************************************
// ****************************************************

// GetAllCraft 所有工艺
func (s baseServer) GetAllCraft(ctx context.Context) (ack AllCraftAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "GetAllCraft 处理请求"))

	var results []entity.Craft

	// 查询总数
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
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "GetAllCraft 处理请求"), zap.Any("处理返回值", ack))
	return
}

// AddCraft 新增工艺
func (s baseServer) AddCraft(ctx context.Context, craft Craft) (ack NewParameterAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "AddCraft 处理请求"))

	if craft.Name == "" {
		err = NoParameters
		return NewParameterAck{}, err
	}

	//查询是否存在craft
	var existCraft entity.Craft
	err = mongodb.NewMgo(collectionCraft).FindOne("name", craft.Name).Decode(&existCraft)
	fmt.Println(err)
	//判断是否存在记录
	if err == nil {
		// 已存在
		ack = NewParameterAck{
			Status:  false,
			Res:     fmt.Sprintf("工艺已存在：%s", craft.Name),
			ErrInfo: "",
		}
	} else if err.Error() == NoDocument.Error() {
		// 不存在，可以新增

		insertOneResult := mongodb.NewMgo(collectionCraft).InsertOne(craft)
		err = nil
		//fmt.Println(insertOneResult)
		ack = NewParameterAck{
			Status:  true,
			Res:     fmt.Sprintf("%s", insertOneResult.InsertedID),
			ErrInfo: "",
		}
	} else {
		// 如果查询错误
		//fmt.Println("***********")
		//fmt.Println(err)
		ack = NewParameterAck{}
	}
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "AddCraft 处理请求"), zap.Any("处理返回值", ack))
	return
}

// DeleteCraftWithId 删除工艺
func (s baseServer) DeleteCraftWithId(ctx context.Context, craftId string) (ack NewParameterAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "DeleteCraftWithId 处理请求"))

	objId, _ := primitive.ObjectIDFromHex(craftId)
	//查询是否存在craft
	var existCraft entity.Craft
	err = mongodb.NewMgo(collectionCraft).FindOne("_id", objId).Decode(&existCraft)
	//fmt.Println(craftId)
	//fmt.Println(err)
	//判断是否存在记录
	if err == nil {
		// 已存在，可删除
		deleteResult := mongodb.NewMgo(collectionCraft).Delete("_id", objId)
		err = nil
		ack = NewParameterAck{
			Status:  true,
			Res:     fmt.Sprintf("工艺已删除：%d", deleteResult),
			ErrInfo: "",
		}
	} else {
		// 如果查询错误
		//fmt.Println(err)
		ack = NewParameterAck{}
	}
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "DeleteCraftWithId 处理请求"), zap.Any("处理返回值", ack))
	return
}

// UpdateCraft 更新工艺
func (s baseServer) UpdateCraft(ctx context.Context, id string, craft Craft) (ack NewParameterAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "UpdateCraft 处理请求"))

	objId, _ := primitive.ObjectIDFromHex(id)

	//查询是否存在craft
	var existCraft entity.Craft
	err = mongodb.NewMgo(collectionCraft).FindOne("_id", objId).Decode(&existCraft)

	//判断是否存在记录
	if err == nil {
		// 已存在,可更新
		// redis 分布式锁
		//res, err := redisdb.RdsClient.SetNX(clientId, existCraft.Id, time.Minute).Result()
		if r, _ := redisdb.GetLock(existCraft.Id, craft.ClientId); r {
			//fmt.Println("***********************")
			//fmt.Println(r)
			// 设置成功
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
				Res:     fmt.Sprintf("%d ms之后解锁", ms),
				ErrInfo: "",
			}
		}

	} else if err.Error() == NoDocument.Error() {
		// 不存在
		//err = errors.New("no errors")
		err = errors.New("no document")
		ack = NewParameterAck{}
	} else {
		// 如果查询错误
		//fmt.Println(err)
		ack = NewParameterAck{}
	}
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "UpdateCraft 处理请求"), zap.Any("处理返回值", ack))
	return
}

// GetAllTexture 所有材质
func (s baseServer) GetAllTexture(ctx context.Context) (ack AllTextureAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "GetAllTexture 处理请求"))

	var results []entity.Texture

	// 查询总数
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
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "GetAllTexture 处理请求"), zap.Any("处理返回值", ack))
	return
}

// AddTexture 新增材质
func (s baseServer) AddTexture(ctx context.Context, texture Texture) (ack NewParameterAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "AddTexture 处理请求"))

	if texture.Name == "" {
		err = NoParameters
		return NewParameterAck{}, err
	}

	//查询是否存在texture
	var existTexture entity.Texture
	err = mongodb.NewMgo(collectionTexture).FindOne("name", texture.Name).Decode(&existTexture)
	//fmt.Println(err)
	//判断是否存在记录
	if err == nil {
		// 已存在
		ack = NewParameterAck{
			Status:  false,
			Res:     fmt.Sprintf("材质已存在：%s", texture.Name),
			ErrInfo: "",
		}
	} else if err.Error() == NoDocument.Error() {
		// 不存在，可以新增
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
		// 如果查询错误
		//fmt.Println(err)
		ack = NewParameterAck{}
	}
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "AddTexture 处理请求"), zap.Any("处理返回值", ack))
	return
}

// DeleteTextureWithId 删除材质
func (s baseServer) DeleteTextureWithId(ctx context.Context, textureId string) (ack NewParameterAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "DeleteTextureWithId 处理请求"))

	objId, _ := primitive.ObjectIDFromHex(textureId)
	//查询是否存在texture
	var existTexture entity.Texture
	err = mongodb.NewMgo(collectionTexture).FindOne("_id", objId).Decode(&existTexture)
	//fmt.Println(craftId)
	fmt.Println(err)
	//判断是否存在记录
	if err == nil {
		// 已存在，可删除
		deleteResult := mongodb.NewMgo(collectionTexture).Delete("_id", objId)
		err = nil
		ack = NewParameterAck{
			Status:  true,
			Res:     fmt.Sprintf("材质已删除：%d", deleteResult),
			ErrInfo: "",
		}
	} else {
		// 如果查询错误
		//fmt.Println(err)
		ack = NewParameterAck{}
	}
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "DeleteTextureWithId 处理请求"), zap.Any("处理返回值", ack))
	return
}

// UpdateTexture 更新材质
func (s baseServer) UpdateTexture(ctx context.Context, id string, texture Texture) (ack NewParameterAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "UpdateTexture 处理请求"))

	objId, _ := primitive.ObjectIDFromHex(id)

	//查询是否存在texture
	var exist entity.Texture
	err = mongodb.NewMgo(collectionTexture).FindOne("_id", objId).Decode(&exist)
	//判断是否存在记录
	if err == nil {
		// 已存在,可更新
		// redis 分布式锁
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
				Res:     fmt.Sprintf("%d ms之后解锁", ms),
				ErrInfo: "",
			}
		}

	} else if err.Error() == NoDocument.Error() {
		// 不存在，可以新增
		//err = errors.New("no errors")
		err = errors.New("no document")
		ack = NewParameterAck{}
	} else {
		// 如果查询错误
		//fmt.Println(err)
		ack = NewParameterAck{}
	}
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "UpdateTexture 处理请求"), zap.Any("处理返回值", ack))
	return
}

// GetAllProcess 所有工序
func (s baseServer) GetAllProcess(ctx context.Context) (ack AllProcessAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "GetAllProcess 处理请求"))

	var results []entity.Process

	// 查询总数
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
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "GetAllProcess 处理请求"), zap.Any("处理返回值", ack))
	return
}

// AddProcess 新增工序
func (s baseServer) AddProcess(ctx context.Context, process Process) (ack NewParameterAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "AddProcess 处理请求"))

	if process.Name == "" {
		err = NoParameters
		return NewParameterAck{}, err
	}

	//查询是否存在process
	var existProcess entity.Process
	err = mongodb.NewMgo(collectionProcess).FindOne("name", process.Name).Decode(&existProcess)
	//fmt.Println(err)
	//判断是否存在记录
	if err == nil {
		// 已存在
		ack = NewParameterAck{
			Status:  false,
			Res:     fmt.Sprintf("工序已存在：%s", process.Name),
			ErrInfo: "",
		}
	} else if err.Error() == NoDocument.Error() {
		// 不存在，可以新增
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
		// 如果查询错误
		//fmt.Println(err)
		ack = NewParameterAck{}
	}
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "AddProcess 处理请求"), zap.Any("处理返回值", ack))
	return
}

// DeleteProcessWithId 删除工序
func (s baseServer) DeleteProcessWithId(ctx context.Context, processId string) (ack NewParameterAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "DeleteProcessWithId 处理请求"))

	objId, _ := primitive.ObjectIDFromHex(processId)
	//查询是否存在Process
	var existProcess entity.Process
	err = mongodb.NewMgo(collectionProcess).FindOne("_id", objId).Decode(&existProcess)
	//fmt.Println(craftId)
	//判断是否存在记录
	if err == nil {
		// 已存在，可删除
		deleteResult := mongodb.NewMgo(collectionProcess).Delete("_id", objId)
		err = nil
		ack = NewParameterAck{
			Status:  true,
			Res:     fmt.Sprintf("工序已删除：%d", deleteResult),
			ErrInfo: "",
		}
	} else {
		// 如果查询错误
		//fmt.Println(err)
		ack = NewParameterAck{}
	}
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "DeleteProcessWithId 处理请求"), zap.Any("处理返回值", ack))
	return
}

// UpdateProcess 更新工序
func (s baseServer) UpdateProcess(ctx context.Context, id string, process Process) (ack NewParameterAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "UpdateProcess 处理请求"))

	objId, _ := primitive.ObjectIDFromHex(id)

	//查询是否存在texture
	var exist entity.Process
	err = mongodb.NewMgo(collectionProcess).FindOne("_id", objId).Decode(&exist)
	//判断是否存在记录
	if err == nil {
		// 已存在,可更新
		// redis 分布式锁
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
				Res:     fmt.Sprintf("%d ms之后解锁", ms),
				ErrInfo: "",
			}
		}

	} else if err.Error() == NoDocument.Error() {
		// 不存在，可以新增
		//err = errors.New("no errors")
		err = errors.New("no document")
		ack = NewParameterAck{}
	} else {
		// 如果查询错误
		//fmt.Println(err)
		ack = NewParameterAck{}
	}
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "UpdateProcess 处理请求"), zap.Any("处理返回值", ack))
	return
}

// GetAllPurchaseStatus 所有采购状态
func (s baseServer) GetAllPurchaseStatus(ctx context.Context) (ack AllPurchaseStatusAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "GetAllPurchaseStatus 处理请求"))

	var results []entity.PurchaseStatus

	// 查询总数
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

	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "GetAllPurchaseStatus 处理请求"), zap.Any("处理返回值", ack))
	return
}

// AddPurchaseStatus 新增采购状态
func (s baseServer) AddPurchaseStatus(ctx context.Context, ps PurchaseStatus) (ack NewParameterAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "AddPurchaseStatus 处理请求"))

	if ps.Name == "" {
		err = NoParameters
		return NewParameterAck{}, err
	}

	//查询是否存在purchase_status
	var existPurchaseStatus entity.PurchaseStatus
	err = mongodb.NewMgo(collectionPurchaseStatus).FindOne("name", ps.Name).Decode(&existPurchaseStatus)
	//fmt.Println(err)
	//判断是否存在记录
	if err == nil {
		// 已存在
		ack = NewParameterAck{
			Status:  false,
			Res:     fmt.Sprintf("采购状态已存在：%s", ps.Name),
			ErrInfo: "",
		}
	} else if err.Error() == NoDocument.Error() {
		// 不存在，可以新增
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
		// 如果查询错误
		//fmt.Println(err)
		ack = NewParameterAck{}
	}
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "AddPurchaseStatus 处理请求"), zap.Any("处理返回值", ack))
	return
}

// DeletePurchaseStatusWithId 删除采购状态
func (s baseServer) DeletePurchaseStatusWithId(ctx context.Context, psId string) (ack NewParameterAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "DeletePurchaseStatusWithId 处理请求"))

	objId, _ := primitive.ObjectIDFromHex(psId)
	//查询是否存在purchase_status
	var existPurchaseStatus entity.PurchaseStatus
	err = mongodb.NewMgo(collectionPurchaseStatus).FindOne("_id", objId).Decode(&existPurchaseStatus)

	//判断是否存在记录
	if err == nil {
		// 已存在，可删除
		deleteResult := mongodb.NewMgo(collectionPurchaseStatus).Delete("_id", objId)
		err = nil
		ack = NewParameterAck{
			Status:  true,
			Res:     fmt.Sprintf("采购状态已删除：%d", deleteResult),
			ErrInfo: "",
		}
	} else {
		// 如果查询错误
		//fmt.Println(err)
		ack = NewParameterAck{}
	}

	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "DeletePurchaseStatusWithId 处理请求"), zap.Any("处理返回值", ack))
	return
}

// UpdatePurchaseStatus 更新采购状态
func (s baseServer) UpdatePurchaseStatus(ctx context.Context, id string, ps PurchaseStatus) (ack NewParameterAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "UpdatePurchaseStatus 处理请求"))

	objId, _ := primitive.ObjectIDFromHex(id)

	//查询是否存在texture
	var exist entity.PurchaseStatus
	err = mongodb.NewMgo(collectionPurchaseStatus).FindOne("_id", objId).Decode(&exist)
	//判断是否存在记录
	if err == nil {
		// 已存在,可更新
		// redis 分布式锁
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
				Res:     fmt.Sprintf("%d ms之后解锁", ms),
				ErrInfo: "",
			}
		}

	} else if err.Error() == NoDocument.Error() {
		// 不存在，
		//err = errors.New("no errors")
		err = errors.New("no document")
		ack = NewParameterAck{}
	} else {
		// 如果查询错误
		//fmt.Println(err)
		ack = NewParameterAck{}
	}
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "UpdatePurchaseStatus 处理请求"), zap.Any("处理返回值", ack))
	return
}
