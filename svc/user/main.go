package main

import (
	"bigSystem/svc/common/db"
	"bigSystem/svc/common/utils"
	"bigSystem/svc/user/Endpoint"
	"bigSystem/svc/user/Service"
	"bigSystem/svc/user/Transport"
	"flag"
	"fmt"
	"go.uber.org/ratelimit"
	"golang.org/x/time/rate"
	"net/http"
)

func main() {

	var (
		listen = flag.String("listen", ":8077", "HTTP listen address")
		//proxy  = flag.String("proxy", "", "Optional comma-separated list of URLs to proxy uppercase requests")
		//servicePort = flag.String("service.port", "8077", "service port")
	)

	flag.Parse()

	utils.NewLoggerServer()

	// 每秒产生10个令牌，令牌桶可以装1个令牌
	golangLimit := rate.NewLimiter(10, 1)
	// 一秒请求一次
	uberLimit := ratelimit.New(1)

	var svc Service.Service
	svc = Service.NewService(&db.DbService{}, utils.GetLogger())

	endpoints := Endpoint.MakeServerEndpoints(svc, utils.GetLogger(), golangLimit, uberLimit)

	var h http.Handler
	{
		h = Transport.MakeHTTPHandler(endpoints, utils.GetLogger())
	}

	//errs := make(chan error)
	//go func() {
	//	utils.GetLogger().Info(fmt.Sprintf("Server Address: %d", listen))
	//	errs <- http.ListenAndServe(*listen, h)
	//}()

	utils.GetLogger().Info(fmt.Sprintf("Server Address: %s", *listen))
	_ = http.ListenAndServe(*listen, h)
	//utils.GetLogger().Error("exit")

}
