package Endpoint

import (
	"bigSystem/svc/weight/Service"
	"context"

	"github.com/go-kit/kit/endpoint"
	"go.uber.org/ratelimit"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type EndpointsServer struct {
	GetAllWeightRecordEndpoint  endpoint.Endpoint
	GetWeightRecordPageEndpoint endpoint.Endpoint
	GetParameterEndpoint        endpoint.Endpoint
	AddNewRecordEndpoint        endpoint.Endpoint
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
		getParameterEndpoint = AuthMiddleware(log)(getParameterEndpoint)
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

	return EndpointsServer{
		GetAllWeightRecordEndpoint:  getAllWeightRecordEndpoint,
		GetWeightRecordPageEndpoint: getWeightRecordPageEndpoint,
		GetParameterEndpoint:        getParameterEndpoint,
		AddNewRecordEndpoint:        addNewRecordEndpoint,
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
