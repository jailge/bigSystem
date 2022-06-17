package Service

import (
	"bigSystem/svc/common/db/mongodb"
	"bigSystem/svc/common/entity"
	"context"
	"errors"
	"fmt"
	"github.com/go-kit/kit/metrics"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	//"go-mdm-service/utils/async"
	"go.uber.org/zap"
)

const (
	collectionRecord    = "record"
	collectionParameter = "parameter"
)

type Service interface {
	//TestAdd(ctx context.Context, in Add) AddAck
	//Login(ctx context.Context, in *pb.Login) (ack *pb.LoginAck, err error)
	//Login(ctx context.Context, in Login) (ack LoginAck, err error)

	GetAllWeightRecord(ctx context.Context) (ack AllDocumentsAck, err error)
	GetWeightRecordPage(ctx context.Context, page Page) (ack AllDocumentsPageAck, err error)
	GetParameter(ctx context.Context) (ack AllParameterAck, err error)
	AddNewRecord(ctx context.Context, newRecord NewRecord) (ack NewRecordAck, err error)
}

type baseServer struct {
	//*mongodb.MdbService
	logger *zap.Logger
}

var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
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
			Status:  false,
			Res:     rData,
			ErrInfo: err.Error(),
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
		Status:  true,
		Res:     rData,
		ErrInfo: "",
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

	//查询是否存在material_code
	var existsRecord entity.WeightRecord
	err = mongodb.NewMgo(collectionRecord).FindOne("material_code", newRecord.MaterialCode).Decode(&existsRecord)

	//判断是否存在记录
	if err == nil {
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
		fmt.Println(same)
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

		updateResult := mongodb.NewMgo(collectionRecord).UpdateMany("material_code", existsRecord.MaterialCode, m)
		//fmt.Println("########", updateResult)

		return NewRecordAck{
			Status:  true,
			Res:     fmt.Sprintf("MatchedCount:%d,  UpsertedCount:%d, ModifiedCount:%d", updateResult.MatchedCount, updateResult.UpsertedCount, updateResult.ModifiedCount),
			ErrInfo: "",
		}, nil

	} else if err.Error() == "mongo: no documents in result" {
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
		return NewRecordAck{
			Status:  true,
			Res:     fmt.Sprintf("%s", insertOneResult.InsertedID),
			ErrInfo: "",
		}, nil
	} else {
		// 如果查询错误
		fmt.Println(err)
		return NewRecordAck{}, err
	}
	//fmt.Println(existsRecord)
	//return

}
