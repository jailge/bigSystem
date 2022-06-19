package service

import (
	//"bigSystem/svc/common/db"
	"bigSystem/svc/common/db/mongodb"
	"bigSystem/svc/common/utils"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	//usergrpc "bigSystem/svc/user/Transport/grpc"
	//"bigSystem/svc/user/pb"
	"bigSystem/svc/weight/Endpoint"
	"bigSystem/svc/weight/Service"
	"bigSystem/svc/weight/Transport"
	//grpctransport "github.com/go-kit/kit/transport/grpc"

	"flag"
	"fmt"
	metricsprometheus "github.com/go-kit/kit/metrics/prometheus"
	//grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/oklog/oklog/pkg/group"
	"go.uber.org/ratelimit"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	//"google.golang.org/grpc"
	"net"
	netHttp "net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	//fs = flag.NewFlagSet("user", flag.ExitOnError)
	httpAddr       = flag.String("http-addr", "0.0.0.0:8177", "HTTP listen address")
	prometheusAddr = flag.String("p", ":10081", "prometheus addr")
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

	count := metricsprometheus.NewCounterFrom(prometheus.CounterOpts{
		Subsystem: "user_agent",
		Name:      "request_count",
		Help:      "Number of requests",
	}, []string{"method"})

	histogram := metricsprometheus.NewHistogramFrom(prometheus.HistogramOpts{
		Subsystem: "user_agent",
		Name:      "request_consume",
		Help:      "Request consumes time",
	}, []string{"method"})

	// 每秒产生10个令牌，令牌桶可以装1个令牌
	golangLimit := rate.NewLimiter(10, 1)
	// 一秒请求一次
	uberLimit := ratelimit.New(1)

	mongodb.Init("./svc/common/config/conf.db.yml")

	var svc Service.Service
	//svc = Service.NewService(&mongodb.MdbService{}, utils.GetLogger())
	svc = Service.NewService(utils.GetLogger(), count, histogram)

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
	initPrometheus(g)
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

// prometheus服务初始化
func initPrometheus(g *group.Group) {
	httpListenerP, err2 := net.Listen("tcp", *prometheusAddr)
	promHandler := NewPromHandler()
	if err2 != nil {
		utils.GetLogger().Error("Server prometheus Listen err: %s", zap.Error(err2))
	}
	g.Add(func() error {
		utils.GetLogger().Info(fmt.Sprintf("Server prometheus start at addr: %s", *prometheusAddr))
		return netHttp.Serve(httpListenerP, promHandler)
	}, func(err error) {
		utils.GetLogger().Error(fmt.Sprintf("Server prometheus Listen close %s", err2.Error()))
	})

}

func NewPromHandler() *mux.Router {
	r := mux.NewRouter()
	r.Handle("/metrics", promhttp.Handler())
	return r
}
