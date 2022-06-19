package Service

import (
	"context"
	"fmt"
	"go.uber.org/zap"
)

const ContextReqUUid = "req_uuid"

type NewMiddlewareServer func(service Service) Service

type logMiddlewareServer struct {
	logger *zap.Logger
	next   Service
}

func (l logMiddlewareServer) AddNewRecord(ctx context.Context, in NewRecord) (out NewRecordAck, err error) {
	defer func() {
		l.logger.Debug(
			fmt.Sprint(ctx.Value(ContextReqUUid)),
			zap.Any("调用 AddNewRecord logMiddlewareServer", "AddNewRecord"),
			zap.Any("req", in),
			zap.Any("res", out),
			zap.Any("err", err))
	}()
	out, err = l.next.AddNewRecord(ctx, in)
	return
}

func (l logMiddlewareServer) GetAllWeightRecord(ctx context.Context) (out AllDocumentsAck, err error) {
	defer func() {
		l.logger.Debug(
			fmt.Sprint(ctx.Value(ContextReqUUid)),
			zap.Any("调用 GetAllWeightRecord logMiddlewareServer", "GetAllWeightRecord"),
			//zap.Any("req", in),
			zap.Any("res", out),
			zap.Any("err", err))
	}()
	out, err = l.next.GetAllWeightRecord(ctx)
	return
}

func (l logMiddlewareServer) GetWeightRecordPage(ctx context.Context, in Page) (out AllDocumentsPageAck, err error) {
	defer func() {
		l.logger.Debug(
			fmt.Sprint(ctx.Value(ContextReqUUid)),
			zap.Any("调用 GetWeightRecordPage logMiddlewareServer", "GetWeightRecordPage"),
			zap.Any("req", in),
			zap.Any("res", out),
			zap.Any("err", err))
	}()
	out, err = l.next.GetWeightRecordPage(ctx, in)
	return
}

func (l logMiddlewareServer) GetParameter(ctx context.Context) (out AllParameterAck, err error) {
	defer func() {
		l.logger.Debug(
			fmt.Sprint(ctx.Value(ContextReqUUid)),
			zap.Any("调用 GetParameter logMiddlewareServer", "GetParameter"),
			//zap.Any("req", in),
			zap.Any("res", out),
			zap.Any("err", err))
	}()
	out, err = l.next.GetParameter(ctx)
	return
}

func (l logMiddlewareServer) SearchWeightWithMaterialCode(ctx context.Context, in MaterialCode) (out WeightMaterialCodeAck, err error) {
	defer func() {
		l.logger.Debug(
			fmt.Sprint(ctx.Value(ContextReqUUid)),
			zap.Any("调用 SearchWeightWithMaterialCode logMiddlewareServer", "SearchWeightWithMaterialCode"),
			zap.Any("req", in),
			zap.Any("res", out),
			zap.Any("err", err))
	}()
	out, err = l.next.SearchWeightWithMaterialCode(ctx, in)
	return
}

func NewLogMiddlewareServer(log *zap.Logger) NewMiddlewareServer {
	return func(service Service) Service {
		return logMiddlewareServer{
			logger: log,
			next:   service,
		}
	}
}

//func (l logMiddlewareServer) Login(ctx context.Context, in Login) (out LoginAck, err error) {
//	defer func() {
//		l.logger.Debug(
//			fmt.Sprint(ctx.Value(ContextReqUUid)),
//			zap.Any("调用 Login logMiddlewareServer", "Login"),
//			zap.Any("req", in),
//			zap.Any("res", out),
//			zap.Any("err", err))
//	}()
//	out, err = l.next.Login(ctx, in)
//	return
//}
