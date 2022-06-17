package Endpoint

import (
	"bigSystem/svc/user/Service"
	"bigSystem/svc/user/pb"
	"context"
	"github.com/go-kit/kit/endpoint"
	"go.uber.org/ratelimit"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type EndpointsServer struct {
	AddEndpoint                     endpoint.Endpoint
	LoginEndpoint                   endpoint.Endpoint
	GetPersonInfoBySnEndPoint       endpoint.Endpoint
	GetPersonsInfoByNameEndPoint    endpoint.Endpoint
	GetAllPersonsInfoEndPoint       endpoint.Endpoint
	SearchPersonsInfoByNameEndPoint endpoint.Endpoint

	RegisterAccountEndPoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service.Service, log *zap.Logger, limit *rate.Limiter, limiter ratelimit.Limiter) EndpointsServer {
	var addEndPoint endpoint.Endpoint
	{
		addEndPoint = MakeAddEndpoint(s)
		addEndPoint = LoggingMiddleware(log)(addEndPoint)
		// 需要鉴权的方法添加jwt鉴权中间件
		//addEndPoint = AuthMiddleware(log)(addEndPoint)
		// 需要限流的方法添加限流中间件
		addEndPoint = NewUberRateMiddleware(limiter)(addEndPoint)
	}
	// 登录endpoint
	var loginEndPoint endpoint.Endpoint
	{
		loginEndPoint = MakeLoginHTTPEndPoint(s)
		loginEndPoint = LoggingMiddleware(log)(loginEndPoint)
		// 需要限流的方法添加限流中间件
		loginEndPoint = NewGolangRateAllowMiddleware(limit)(loginEndPoint)
	}
	var registerAccountEndPoint endpoint.Endpoint
	{
		registerAccountEndPoint = MakeRegisterAccountEndPoint(s)
		// 需要鉴权的方法添加jwt鉴权中间件
		registerAccountEndPoint = AuthMiddleware(log)(registerAccountEndPoint)
		registerAccountEndPoint = LoggingMiddleware(log)(registerAccountEndPoint)
		// 需要限流的方法添加限流中间件
		registerAccountEndPoint = NewUberRateMiddleware(limiter)(registerAccountEndPoint)
	}
	// GetPersonInfoBySn endpoint
	var getPersonInfoBySnEndPoint endpoint.Endpoint
	{
		getPersonInfoBySnEndPoint = MakeGetPersonInfoBySnEndPoint(s)
		getPersonInfoBySnEndPoint = LoggingMiddleware(log)(getPersonInfoBySnEndPoint)
		// 需要限流的方法添加限流中间件
		getPersonInfoBySnEndPoint = NewUberRateMiddleware(limiter)(getPersonInfoBySnEndPoint)
	}

	var getPersonsInfoByNameEndPoint endpoint.Endpoint
	{
		getPersonsInfoByNameEndPoint = MakeGetPersonsInfoByNameEndPoint(s)
		getPersonsInfoByNameEndPoint = LoggingMiddleware(log)(getPersonsInfoByNameEndPoint)
		// 需要限流的方法添加限流中间件
		getPersonsInfoByNameEndPoint = NewUberRateMiddleware(limiter)(getPersonsInfoByNameEndPoint)
	}

	var getAllPersonsInfoEndPoint endpoint.Endpoint
	{
		getAllPersonsInfoEndPoint = MakeGetAllPersonsInfoEndPoint(s)
		getAllPersonsInfoEndPoint = LoggingMiddleware(log)(getAllPersonsInfoEndPoint)
		// 需要限流的方法添加限流中间件
		getAllPersonsInfoEndPoint = NewUberRateMiddleware(limiter)(getAllPersonsInfoEndPoint)
	}

	var searchPersonsInfoByNameEndPoint endpoint.Endpoint
	{
		searchPersonsInfoByNameEndPoint = MakeSearchPersonsInfoByNameEndPoint(s)
		searchPersonsInfoByNameEndPoint = LoggingMiddleware(log)(searchPersonsInfoByNameEndPoint)
		// 需要限流的方法添加限流中间件
		searchPersonsInfoByNameEndPoint = NewUberRateMiddleware(limiter)(searchPersonsInfoByNameEndPoint)
	}
	return EndpointsServer{
		AddEndpoint:                     addEndPoint,
		LoginEndpoint:                   loginEndPoint,
		GetPersonInfoBySnEndPoint:       getPersonInfoBySnEndPoint,
		GetPersonsInfoByNameEndPoint:    getPersonsInfoByNameEndPoint,
		GetAllPersonsInfoEndPoint:       getAllPersonsInfoEndPoint,
		SearchPersonsInfoByNameEndPoint: searchPersonsInfoByNameEndPoint,
		RegisterAccountEndPoint:         registerAccountEndPoint,
	}
}

func (s EndpointsServer) Add(ctx context.Context, in Service.Add) Service.AddAck {
	res, _ := s.AddEndpoint(ctx, in)
	return res.(Service.AddAck)
}

// MakeAddEndpoint Service中的TestAdd，转换成endpoint.Endpoint
func MakeAddEndpoint(s Service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Service.Add)
		res := s.TestAdd(ctx, req)
		return res, nil
	}
}

func (s EndpointsServer) Login(ctx context.Context, in *pb.Login) (*pb.LoginAck, error) {
	res, err := s.LoginEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.LoginAck), nil
}

func MakeLoginEndPoint(s Service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.Login)
		return s.Login(ctx, req)
	}
}

func (s EndpointsServer) LoginHTTP(ctx context.Context, in Service.Login) (Service.LoginAck, error) {
	res, err := s.LoginEndpoint(ctx, in)
	return res.(Service.LoginAck), err
}

func MakeLoginHTTPEndPoint(s Service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Service.Login)
		return s.LoginHTTP(ctx, req)
	}
}

func (s EndpointsServer) RegisterAccount(ctx context.Context, in Service.User) Service.UserAck {
	res, _ := s.RegisterAccountEndPoint(ctx, in)
	return res.(Service.UserAck)
}

func MakeRegisterAccountEndPoint(s Service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Service.User)
		return s.RegisterAccount(ctx, req)
	}
}

func MakeGetPersonInfoBySnEndPoint(s Service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Service.PersonSn)
		return s.GetPersonInfoBySn(ctx, req)
	}
}

func MakeGetPersonsInfoByNameEndPoint(s Service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Service.PersonName)
		return s.GetPersonsInfoByName(ctx, req)
	}
}

func MakeGetAllPersonsInfoEndPoint(s Service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Service.AllPerson)
		return s.GetAllPersonsInfo(ctx, req)
	}
}

func MakeSearchPersonsInfoByNameEndPoint(s Service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(Service.SearchPersons)
		return s.SearchPersonsInfoByName(ctx, req)
	}
}
