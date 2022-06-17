package Endpoint

import (
	"bigSystem/svc/common/utils"
	"bigSystem/svc/weight/Service"
	"context"
	"errors"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"go.uber.org/ratelimit"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"time"
)

// LoggingMiddleware 日志中间件
func LoggingMiddleware(logger *zap.Logger) endpoint.Middleware {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				logger.Debug(
					fmt.Sprint(ctx.Value(Service.ContextReqUUid)),
					zap.Any("调用 endpoint LoggingMiddleware", "处理完请求"),
					zap.Any("耗时毫秒", time.Since(begin).Milliseconds()),
				)
			}(time.Now())
			return e(ctx, request)
		}
	}
}

// AuthMiddleware jwt鉴权中间件
func AuthMiddleware(logger *zap.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			token := fmt.Sprint(ctx.Value(utils.JWT_CONTEXT_KEY))
			if token == "" {
				err = errors.New("请登录")
				logger.Debug(fmt.Sprint(ctx.Value(Service.ContextReqUUid)), zap.Any("[AuthMiddleware]", "token == empty"), zap.Error(err))
				return "", err
			}
			jwtInfo, err := utils.ParseToken(token)
			if err != nil {
				logger.Debug(fmt.Sprint(ctx.Value(Service.ContextReqUUid)), zap.Any("[AuthMiddleware]", "ParseToken"), zap.Error(err))
				return "", err
			}
			if v, ok := jwtInfo["Name"]; ok {
				ctx = context.WithValue(ctx, "name", v)
			}
			return next(ctx, request)
		}
	}
}

// NewGolangRateWaitMiddleware 基于golang.org/x/time/rate的限流中间件
func NewGolangRateWaitMiddleware(limit *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if err = limit.Wait(ctx); err != nil {
				return "", errors.New("limit req  Wait")
			}
			return next(ctx, request)
		}
	}
}

// NewGolangRateAllowMiddleware 基于go.uber.org/ratelimit的限流中间件
func NewGolangRateAllowMiddleware(limit *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !limit.Allow() {
				return "", errors.New("limit req  Allow")
			}
			return next(ctx, request)
		}
	}
}

// NewUberRateMiddleware 基于go.uber.org/ratelimit的限流中间件
func NewUberRateMiddleware(limit ratelimit.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			limit.Take()
			return next(ctx, request)
		}
	}
}
