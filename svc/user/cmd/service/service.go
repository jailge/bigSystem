package service

import (
	"bigSystem/svc/common/db"
	"bigSystem/svc/common/utils"
	"bigSystem/svc/user/Endpoint"
	"bigSystem/svc/user/Service"
	"bigSystem/svc/user/Transport"
	usergrpc "bigSystem/svc/user/Transport/grpc"
	"bigSystem/svc/user/pb"
	"flag"
	"fmt"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/oklog/oklog/pkg/group"
	"go.uber.org/ratelimit"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"net"
	netHttp "net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	//fs = flag.NewFlagSet("user", flag.ExitOnError)
	httpAddr = flag.String("http-addr", ":8077", "HTTP listen address")
	grpcAddr = flag.String("grpc-addr", ":8078", "gRPC listen address")
)

func Run() {

	//var (
	//	listen = flag.String("listen", ":8077", "HTTP listen address")
	//	//proxy  = flag.String("proxy", "", "Optional comma-separated list of URLs to proxy uppercase requests")
	//	//servicePort = flag.String("service.port", "8077", "service port")
	//)

	flag.Parse()
	//if err := fs.Parse(os.Args[1:]); err != nil {
	//	panic(err)
	//}

	utils.NewLoggerServer()

	// 每秒产生10个令牌，令牌桶可以装1个令牌
	golangLimit := rate.NewLimiter(10, 1)
	// 一秒请求一次
	uberLimit := ratelimit.New(1)

	var svc Service.Service
	svc = Service.NewService(&db.DbService{}, utils.GetLogger())
	//svc = Service.NewService(utils.GetLogger())

	endpoints := Endpoint.MakeServerEndpoints(svc, utils.GetLogger(), golangLimit, uberLimit)

	//var h http.Handler
	//{
	//	h = Transport.MakeHTTPHandler(endpoints, utils.GetLogger())
	//}

	//errs := make(chan error)
	//go func() {
	//	utils.GetLogger().Info(fmt.Sprintf("Server Address: %d", listen))
	//	errs <- http.ListenAndServe(*listen, h)
	//}()

	//utils.GetLogger().Info(fmt.Sprintf("Server Address: %s", *listen))
	//_ = http.ListenAndServe(*listen, h)
	//utils.GetLogger().Error("exit")

	g := &group.Group{}
	initHttpHandler(endpoints, g)
	initGRPCHandler(endpoints, g)
	initCancelInterrupt(g)

	utils.GetLogger().Error("exit", zap.Error(g.Run()))
}

func initCancelInterrupt(g *group.Group) {
	cancelInterrupt := make(chan struct{})
	g.Add(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			return fmt.Errorf("received signal %s", sig)
		case <-cancelInterrupt:
			return nil
		}
	}, func(error) {
		close(cancelInterrupt)
	})
}

func initHttpHandler(endpoints Endpoint.EndpointsServer, g *group.Group) {
	//opts := []kithttp.ServerOption{
	//	kithttp.ServerErrorHandler(Transport.NewZapLogErrorHandler(utils.GetLogger())),
	//	kithttp.ServerErrorEncoder(encode.JsonError),
	//}

	var httpHandler netHttp.Handler
	httpHandler = Transport.MakeHTTPHandler(endpoints, utils.GetLogger())
	httpListener, err := net.Listen("tcp", *httpAddr)
	if err != nil {
		utils.GetLogger().Error("Server HTTP Listen err: %s", zap.Error(err))
	}
	g.Add(func() error {
		utils.GetLogger().Info(fmt.Sprintf("HTTP Server start at addr: %s", *httpAddr))
		return netHttp.Serve(httpListener, httpHandler)
	}, func(error) {
		utils.GetLogger().Error(fmt.Sprintf("Server HTTP Listen close: %s", httpListener.Close()))
	})
}

func initGRPCHandler(endpoints Endpoint.EndpointsServer, g *group.Group) {
	grpcListener, err := net.Listen("tcp", *grpcAddr)
	if err != nil {
		utils.GetLogger().Error("Server gRPC Listen err: ", zap.Error(err))
	}
	g.Add(func() error {
		grpcServer := usergrpc.NewGRPCServer(endpoints, utils.GetLogger())

		baseServer := grpc.NewServer(grpc.UnaryInterceptor(grpctransport.Interceptor))
		pb.RegisterUserServer(baseServer, grpcServer)
		utils.GetLogger().Info(fmt.Sprintf("gRPC Server start at addr: %s", *grpcAddr))
		return baseServer.Serve(grpcListener)
	}, func(error) {
		utils.GetLogger().Warn("Server gRPC err: ", zap.Error(err))
		grpcListener.Close()
	})
}
