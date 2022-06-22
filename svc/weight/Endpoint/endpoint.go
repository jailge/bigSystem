package Endpoint

import (
	"bigSystem/svc/common/entity"
	"bigSystem/svc/weight/Service"
	"context"
	"github.com/go-kit/kit/endpoint"
	"go.uber.org/ratelimit"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type EndpointsServer struct {
	GetAllWeightRecordEndpoint           endpoint.Endpoint
	GetWeightRecordPageEndpoint          endpoint.Endpoint
	GetParameterEndpoint                 endpoint.Endpoint
	AddNewRecordEndpoint                 endpoint.Endpoint
	SearchWeightWithMaterialCodeEndpoint endpoint.Endpoint

	GetAllCraftEndpoint                endpoint.Endpoint
	GetAllTextureEndpoint              endpoint.Endpoint
	GetAllProcessEndpoint              endpoint.Endpoint
	GetAllPurchaseStatusEndpoint       endpoint.Endpoint
	AddCraftEndpoint                   endpoint.Endpoint
	AddTextureEndpoint                 endpoint.Endpoint
	AddProcessEndpoint                 endpoint.Endpoint
	AddPurchaseStatusEndpoint          endpoint.Endpoint
	DeleteCraftWithIdEndpoint          endpoint.Endpoint
	DeleteTextureWithIdEndpoint        endpoint.Endpoint
	DeleteProcessWithIdEndpoint        endpoint.Endpoint
	DeletePurchaseStatusWithIdEndpoint endpoint.Endpoint
	UpdateCraftEndpoint                endpoint.Endpoint
	UpdateTextureEndpoint              endpoint.Endpoint
	UpdateProcessEndpoint              endpoint.Endpoint
	UpdatePurchaseStatusEndpoint       endpoint.Endpoint
}

func MakeServerEndpoints(s Service.Service, log *zap.Logger, limit *rate.Limiter, limiter ratelimit.Limiter) EndpointsServer {
	var getAllWeightRecordEndpoint endpoint.Endpoint
	{
		getAllWeightRecordEndpoint = MakeGetAllWeightRecordEndpoint(s)
		getAllWeightRecordEndpoint = LoggingMiddleware(log)(getAllWeightRecordEndpoint)
		// 需要鉴权的方法添加jwt鉴权中间件
		getAllWeightRecordEndpoint = AuthMiddleware(log)(getAllWeightRecordEndpoint)
		// 需要限流的方法添加限流中间件
		getAllWeightRecordEndpoint = NewUberRateMiddleware(limiter)(getAllWeightRecordEndpoint)
	}
	var getWeightRecordPageEndpoint endpoint.Endpoint
	{
		getWeightRecordPageEndpoint = MakeGetWeightRecordPageEndpoint(s)
		getWeightRecordPageEndpoint = LoggingMiddleware(log)(getWeightRecordPageEndpoint)
		// 需要鉴权的方法添加jwt鉴权中间件
		getWeightRecordPageEndpoint = AuthMiddleware(log)(getWeightRecordPageEndpoint)
		// 需要限流的方法添加限流中间件
		getWeightRecordPageEndpoint = NewUberRateMiddleware(limiter)(getWeightRecordPageEndpoint)
	}
	var getParameterEndpoint endpoint.Endpoint
	{
		getParameterEndpoint = MakeGetParameterEndpoint(s)
		getParameterEndpoint = LoggingMiddleware(log)(getParameterEndpoint)
		// 需要鉴权的方法添加jwt鉴权中间件
		//getParameterEndpoint = AuthMiddleware(log)(getParameterEndpoint)
		// 需要限流的方法添加限流中间件
		getParameterEndpoint = NewUberRateMiddleware(limiter)(getParameterEndpoint)
	}
	var addNewRecordEndpoint endpoint.Endpoint
	{
		addNewRecordEndpoint = MakeAddNewRecordEndpoint(s)
		addNewRecordEndpoint = LoggingMiddleware(log)(addNewRecordEndpoint)
		// 需要鉴权的方法添加jwt鉴权中间件
		//addNewRecordEndpoint = AuthMiddleware(log)(addNewRecordEndpoint)
		// 需要限流的方法添加限流中间件JWT
		addNewRecordEndpoint = NewUberRateMiddleware(limiter)(addNewRecordEndpoint)
	}
	var searchWeightWithMaterialCodeEndpoint endpoint.Endpoint
	{
		searchWeightWithMaterialCodeEndpoint = MakeSearchWeightWithMaterialCodeEndpoint(s)
		searchWeightWithMaterialCodeEndpoint = LoggingMiddleware(log)(searchWeightWithMaterialCodeEndpoint)
		// 需要鉴权的方法添加jwt鉴权中间件
		//addNewRecordEndpoint = AuthMiddleware(log)(addNewRecordEndpoint)
		// 需要限流的方法添加限流中间件JWT
		searchWeightWithMaterialCodeEndpoint = NewUberRateMiddleware(limiter)(searchWeightWithMaterialCodeEndpoint)
	}
	var getAllCraftEndpoint endpoint.Endpoint
	{
		getAllCraftEndpoint = MakeGetAllCraftEndpoint(s)
		getAllCraftEndpoint = LoggingMiddleware(log)(getAllCraftEndpoint)
		// 需要鉴权的方法添加jwt鉴权中间件
		//getAllCraftEndpoint = AuthMiddleware(log)(getAllCraftEndpoint)
		// 需要限流的方法添加限流中间件JWT
		getAllCraftEndpoint = NewUberRateMiddleware(limiter)(getAllCraftEndpoint)
	}
	var getAllTextureEndpoint endpoint.Endpoint
	{
		getAllTextureEndpoint = MakeGetAllTextureEndpoint(s)
		getAllTextureEndpoint = LoggingMiddleware(log)(getAllTextureEndpoint)
		// 需要鉴权的方法添加jwt鉴权中间件
		//getAllCraftEndpoint = AuthMiddleware(log)(getAllCraftEndpoint)
		// 需要限流的方法添加限流中间件JWT
		getAllTextureEndpoint = NewUberRateMiddleware(limiter)(getAllTextureEndpoint)
	}
	var getAllProcessEndpoint endpoint.Endpoint
	{
		getAllProcessEndpoint = MakeGetAllProcessEndpoint(s)
		getAllProcessEndpoint = LoggingMiddleware(log)(getAllProcessEndpoint)
		// 需要鉴权的方法添加jwt鉴权中间件
		//getAllCraftEndpoint = AuthMiddleware(log)(getAllCraftEndpoint)
		// 需要限流的方法添加限流中间件JWT
		getAllProcessEndpoint = NewUberRateMiddleware(limiter)(getAllProcessEndpoint)
	}
	var getAllPurchaseStatusEndpoint endpoint.Endpoint
	{
		getAllPurchaseStatusEndpoint = MakeGetAllPurchaseStatusEndpoint(s)
		getAllPurchaseStatusEndpoint = LoggingMiddleware(log)(getAllPurchaseStatusEndpoint)
		// 需要鉴权的方法添加jwt鉴权中间件
		//getAllCraftEndpoint = AuthMiddleware(log)(getAllCraftEndpoint)
		// 需要限流的方法添加限流中间件JWT
		getAllPurchaseStatusEndpoint = NewUberRateMiddleware(limiter)(getAllPurchaseStatusEndpoint)
	}

	var addCraftEndpoint endpoint.Endpoint
	{
		addCraftEndpoint = MakeAddCraftEndpoint(s)
		addCraftEndpoint = LoggingMiddleware(log)(addCraftEndpoint)
		// 需要鉴权的方法添加jwt鉴权中间件
		addCraftEndpoint = AuthMiddleware(log)(addCraftEndpoint)
		// 需要限流的方法添加限流中间件JWT
		addCraftEndpoint = NewUberRateMiddleware(limiter)(addCraftEndpoint)
	}
	var addPurchaseStatusEndpoint endpoint.Endpoint
	{
		addPurchaseStatusEndpoint = MakeAddPurchaseStatusEndpoint(s)
		addPurchaseStatusEndpoint = LoggingMiddleware(log)(addPurchaseStatusEndpoint)
		// 需要鉴权的方法添加jwt鉴权中间件
		addPurchaseStatusEndpoint = AuthMiddleware(log)(addPurchaseStatusEndpoint)
		// 需要限流的方法添加限流中间件JWT
		addPurchaseStatusEndpoint = NewUberRateMiddleware(limiter)(addPurchaseStatusEndpoint)
	}
	var addTextureEndpoint endpoint.Endpoint
	{
		addTextureEndpoint = MakeAddTextureEndpoint(s)
		addTextureEndpoint = LoggingMiddleware(log)(addTextureEndpoint)
		// 需要鉴权的方法添加jwt鉴权中间件
		addTextureEndpoint = AuthMiddleware(log)(addTextureEndpoint)
		// 需要限流的方法添加限流中间件JWT
		addTextureEndpoint = NewUberRateMiddleware(limiter)(addTextureEndpoint)
	}
	var addProcessEndpoint endpoint.Endpoint
	{
		addProcessEndpoint = MakeAddProcessEndpoint(s)
		addProcessEndpoint = LoggingMiddleware(log)(addProcessEndpoint)
		// 需要鉴权的方法添加jwt鉴权中间件
		addProcessEndpoint = AuthMiddleware(log)(addProcessEndpoint)
		// 需要限流的方法添加限流中间件JWT
		addProcessEndpoint = NewUberRateMiddleware(limiter)(addProcessEndpoint)
	}

	var deleteTextureWithIdEndpoint endpoint.Endpoint
	{
		deleteTextureWithIdEndpoint = MakeDeleteTextureWithIdEndpoint(s)
		deleteTextureWithIdEndpoint = LoggingMiddleware(log)(deleteTextureWithIdEndpoint)
		// 需要鉴权的方法添加jwt鉴权中间件
		deleteTextureWithIdEndpoint = AuthMiddleware(log)(deleteTextureWithIdEndpoint)
		// 需要限流的方法添加限流中间件JWT
		deleteTextureWithIdEndpoint = NewUberRateMiddleware(limiter)(deleteTextureWithIdEndpoint)
	}
	var deleteCraftWithIdEndpoint endpoint.Endpoint
	{
		deleteCraftWithIdEndpoint = MakeDeleteCraftWithIdEndpoint(s)
		deleteCraftWithIdEndpoint = LoggingMiddleware(log)(deleteCraftWithIdEndpoint)
		// 需要鉴权的方法添加jwt鉴权中间件
		deleteCraftWithIdEndpoint = AuthMiddleware(log)(deleteCraftWithIdEndpoint)
		// 需要限流的方法添加限流中间件JWT
		deleteCraftWithIdEndpoint = NewUberRateMiddleware(limiter)(deleteCraftWithIdEndpoint)
	}
	var deleteProcessWithIdEndpoint endpoint.Endpoint
	{
		deleteProcessWithIdEndpoint = MakeDeleteProcessWithIdEndpoint(s)
		deleteProcessWithIdEndpoint = LoggingMiddleware(log)(deleteProcessWithIdEndpoint)
		// 需要鉴权的方法添加jwt鉴权中间件
		deleteProcessWithIdEndpoint = AuthMiddleware(log)(deleteProcessWithIdEndpoint)
		// 需要限流的方法添加限流中间件JWT
		deleteProcessWithIdEndpoint = NewUberRateMiddleware(limiter)(deleteProcessWithIdEndpoint)
	}
	var deletePurchaseStatusWithIdEndpoint endpoint.Endpoint
	{
		deletePurchaseStatusWithIdEndpoint = MakeDeletePurchaseStatusWithIdEndpoint(s)
		deletePurchaseStatusWithIdEndpoint = LoggingMiddleware(log)(deletePurchaseStatusWithIdEndpoint)
		// 需要鉴权的方法添加jwt鉴权中间件
		deletePurchaseStatusWithIdEndpoint = AuthMiddleware(log)(deletePurchaseStatusWithIdEndpoint)
		// 需要限流的方法添加限流中间件JWT
		deletePurchaseStatusWithIdEndpoint = NewUberRateMiddleware(limiter)(deletePurchaseStatusWithIdEndpoint)
	}

	var updateCraftEndpoint endpoint.Endpoint
	{
		updateCraftEndpoint = MakeUpdateCraftEndpoint(s)
		updateCraftEndpoint = LoggingMiddleware(log)(updateCraftEndpoint)
		// 需要鉴权的方法添加jwt鉴权中间件
		updateCraftEndpoint = AuthMiddleware(log)(updateCraftEndpoint)
		// 需要限流的方法添加限流中间件JWT
		updateCraftEndpoint = NewUberRateMiddleware(limiter)(updateCraftEndpoint)
	}
	var updateTextureEndpoint endpoint.Endpoint
	{
		updateTextureEndpoint = MakeUpdateTextureEndpoint(s)
		updateTextureEndpoint = LoggingMiddleware(log)(updateTextureEndpoint)
		// 需要鉴权的方法添加jwt鉴权中间件
		updateTextureEndpoint = AuthMiddleware(log)(updateTextureEndpoint)
		// 需要限流的方法添加限流中间件JWT
		updateTextureEndpoint = NewUberRateMiddleware(limiter)(updateTextureEndpoint)
	}
	var updateProcessEndpoint endpoint.Endpoint
	{
		updateProcessEndpoint = MakeUpdateProcessEndpoint(s)
		updateProcessEndpoint = LoggingMiddleware(log)(updateProcessEndpoint)
		// 需要鉴权的方法添加jwt鉴权中间件
		updateProcessEndpoint = AuthMiddleware(log)(updateProcessEndpoint)
		// 需要限流的方法添加限流中间件JWT
		updateProcessEndpoint = NewUberRateMiddleware(limiter)(updateProcessEndpoint)
	}
	var updatePurchaseStatusEndpoint endpoint.Endpoint
	{
		updatePurchaseStatusEndpoint = MakeUpdatePurchaseStatusEndpoint(s)
		updatePurchaseStatusEndpoint = LoggingMiddleware(log)(updatePurchaseStatusEndpoint)
		// 需要鉴权的方法添加jwt鉴权中间件
		updatePurchaseStatusEndpoint = AuthMiddleware(log)(updatePurchaseStatusEndpoint)
		// 需要限流的方法添加限流中间件JWT
		updatePurchaseStatusEndpoint = NewUberRateMiddleware(limiter)(updatePurchaseStatusEndpoint)
	}

	return EndpointsServer{
		GetAllWeightRecordEndpoint:           getAllWeightRecordEndpoint,
		GetWeightRecordPageEndpoint:          getWeightRecordPageEndpoint,
		GetParameterEndpoint:                 getParameterEndpoint,
		AddNewRecordEndpoint:                 addNewRecordEndpoint,
		SearchWeightWithMaterialCodeEndpoint: searchWeightWithMaterialCodeEndpoint,

		GetAllCraftEndpoint:                getAllCraftEndpoint,
		GetAllTextureEndpoint:              getAllTextureEndpoint,
		GetAllProcessEndpoint:              getAllProcessEndpoint,
		GetAllPurchaseStatusEndpoint:       getAllPurchaseStatusEndpoint,
		AddCraftEndpoint:                   addCraftEndpoint,
		AddProcessEndpoint:                 addCraftEndpoint,
		AddTextureEndpoint:                 addTextureEndpoint,
		AddPurchaseStatusEndpoint:          addPurchaseStatusEndpoint,
		DeleteCraftWithIdEndpoint:          deleteCraftWithIdEndpoint,
		DeleteProcessWithIdEndpoint:        deleteProcessWithIdEndpoint,
		DeleteTextureWithIdEndpoint:        deleteTextureWithIdEndpoint,
		DeletePurchaseStatusWithIdEndpoint: deletePurchaseStatusWithIdEndpoint,
		UpdateCraftEndpoint:                updateCraftEndpoint,
		UpdateTextureEndpoint:              updateTextureEndpoint,
		UpdateProcessEndpoint:              updateProcessEndpoint,
		UpdatePurchaseStatusEndpoint:       updatePurchaseStatusEndpoint,
	}
}

//func (s EndpointsServer) Add(ctx context.Context, in Service.Add) Service.AddAck {
//	res, _ := s.AddEndpoint(ctx, in)
//	return res.(Service.AddAck)
//}
//
//
//// MakeAddEndpoint Service中的TestAdd，转换成endpoint.Endpoint
//func  MakeAddEndpoint(s Service.Service) endpoint.Endpoint {
//	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
//		req := request.(Service.Add)
//		res := s.TestAdd(ctx, req)
//		return res, nil
//	}
//}

//func (s EndpointsServer) Login(ctx context.Context, in *pb.Login) (*pb.LoginAck, error) {
//	res, err := s.LoginEndpoint(ctx, in)
//	if err != nil {
//		return nil, err
//	}
//	return res.(*pb.LoginAck), nil
//}

//func (s EndpointsServer) Login(ctx context.Context, in Service.Login) (Service.LoginAck, error) {
//	res, err := s.LoginEndpoint(ctx, in)
//	return res.(Service.LoginAck), err
//}

//func MakeLoginEndPoint(s Service.Service) endpoint.Endpoint {
//	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
//		req := request.(*pb.Login)
//		return s.Login(ctx, req)
//	}
//}

//
//func MakeLoginEndPoint(s Service.Service) endpoint.Endpoint {
//	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
//		req := request.(Service.Login)
//		return s.Login(ctx, req)
//	}
//}

//func (s EndpointsServer) GetAllWeightRecord(ctx context.Context) Service.AllDocumentsAck {
//	res, _ := s.GetAllWeightRecordEndpoint(ctx)
//	return res.(Service.AllDocumentsAck)
//}

// MakeGetAllWeightRecordEndpoint Service中的GetAllWeightRecord，转换成endpoint.Endpoint
func MakeGetAllWeightRecordEndpoint(s Service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.GetAllWeightRecord(ctx)
	}
}

func MakeGetWeightRecordPageEndpoint(s Service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Service.Page)
		return s.GetWeightRecordPage(ctx, req)
	}
}

func MakeGetParameterEndpoint(s Service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.GetParameter(ctx)
	}
}

func MakeAddNewRecordEndpoint(s Service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Service.NewRecord)
		return s.AddNewRecord(ctx, req)
	}
}

func MakeSearchWeightWithMaterialCodeEndpoint(s Service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Service.MaterialCode)
		return s.SearchWeightWithMaterialCode(ctx, req)
	}
}

// **************************

// **************************

func MakeGetAllCraftEndpoint(s Service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.GetAllCraft(ctx)
	}
}

func MakeGetAllTextureEndpoint(s Service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.GetAllTexture(ctx)
	}
}

func MakeGetAllProcessEndpoint(s Service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.GetAllProcess(ctx)
	}
}

func MakeGetAllPurchaseStatusEndpoint(s Service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.GetAllPurchaseStatus(ctx)
	}
}

func MakeAddCraftEndpoint(s Service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Service.Craft)
		return s.AddCraft(ctx, req)
	}
}

func MakeAddTextureEndpoint(s Service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Service.Texture)
		return s.AddTexture(ctx, req)
	}
}

func MakeAddProcessEndpoint(s Service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Service.Process)
		return s.AddProcess(ctx, req)
	}
}

func MakeAddPurchaseStatusEndpoint(s Service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Service.PurchaseStatus)
		return s.AddPurchaseStatus(ctx, req)
	}
}

func MakeDeleteCraftWithIdEndpoint(s Service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		//fmt.Println("request:  ")
		//fmt.Println(request)
		req := request.(string)
		return s.DeleteCraftWithId(ctx, req)
	}
}

func MakeDeleteTextureWithIdEndpoint(s Service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		//fmt.Println("request:  ")
		//fmt.Println(request)
		req := request.(string)
		return s.DeleteTextureWithId(ctx, req)
	}
}

func MakeDeleteProcessWithIdEndpoint(s Service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		//fmt.Println("request:  ")
		//fmt.Println(request)
		req := request.(string)
		return s.DeleteProcessWithId(ctx, req)
	}
}

func MakeDeletePurchaseStatusWithIdEndpoint(s Service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)
		return s.DeletePurchaseStatusWithId(ctx, req)
	}
}

func MakeUpdateCraftEndpoint(s Service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(entity.Craft)
		craft := Service.Craft{Name: req.Name}
		return s.UpdateCraft(ctx, req.Id, craft)
	}
}

func MakeUpdateTextureEndpoint(s Service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(entity.Texture)
		texture := Service.Texture{Name: req.Name}
		return s.UpdateTexture(ctx, req.Id, texture)
	}
}

func MakeUpdateProcessEndpoint(s Service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(entity.Process)
		process := Service.Process{Name: req.Name}
		return s.UpdateProcess(ctx, req.Id, process)
	}
}

func MakeUpdatePurchaseStatusEndpoint(s Service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(entity.PurchaseStatus)
		purchaseStatus := Service.PurchaseStatus{Name: req.Name}
		return s.UpdatePurchaseStatus(ctx, req.Id, purchaseStatus)
	}
}
