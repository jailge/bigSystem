package Service

import (
	"bigSystem/svc/user/pb"
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

func NewLogMiddlewareServer(log *zap.Logger) NewMiddlewareServer {
	return func(service Service) Service {
		return logMiddlewareServer{
			logger: log,
			next:   service,
		}
	}
}

func (l logMiddlewareServer) Login(ctx context.Context, in *pb.Login) (out *pb.LoginAck, err error) {
	defer func() {
		l.logger.Debug(
			fmt.Sprint(ctx.Value(ContextReqUUid)),
			zap.Any("调用 Login logMiddlewareServer", "Login"),
			zap.Any("req", in),
			zap.Any("res", out),
			zap.Any("err", err))
	}()
	out, err = l.next.Login(ctx, in)
	return
}

func (l logMiddlewareServer) LoginHTTP(ctx context.Context, in Login) (out LoginAck, err error) {
	defer func() {
		l.logger.Debug(
			fmt.Sprint(ctx.Value(ContextReqUUid)),
			zap.Any("调用 LoginHTTP logMiddlewareServer", "LoginHTTP"),
			zap.Any("req", in),
			zap.Any("res", out),
			zap.Any("err", err))
	}()
	out, err = l.next.LoginHTTP(ctx, in)
	return
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

func (l logMiddlewareServer) TestAdd(ctx context.Context, in Add) (out AddAck) {
	defer func() {
		l.logger.Debug(
			fmt.Sprint(ctx.Value(ContextReqUUid)),
			zap.Any("调用 Service logMiddlewareServer", "TestAdd"),
			zap.Any("req", in),
			zap.Any("res", out),
		)
	}()
	out = l.next.TestAdd(ctx, in)
	return out
}

func (l logMiddlewareServer) GetPersonInfoBySn(ctx context.Context, in PersonSn) (out PersonSnAck, err error) {
	defer func() {
		l.logger.Debug(
			fmt.Sprint(ctx.Value(ContextReqUUid)),
			zap.Any("调用 Service logMiddlewareServer", "GetPersonInfoBySn"),
			zap.Any("req", in),
			zap.Any("res", out),
		)
	}()
	out, err = l.next.GetPersonInfoBySn(ctx, in)
	return
}

func (l logMiddlewareServer) GetPersonsInfoByName(ctx context.Context, in PersonName) (out PersonsAck, err error) {
	defer func() {
		l.logger.Debug(
			fmt.Sprint(ctx.Value(ContextReqUUid)),
			zap.Any("调用 Service logMiddlewareServer", "GetPersonsInfoByName"),
			zap.Any("req", in),
			zap.Any("res", out),
		)
	}()
	out, err = l.next.GetPersonsInfoByName(ctx, in)
	return
}

func (l logMiddlewareServer) GetAllPersonsInfo(ctx context.Context, in AllPerson) (out PersonsAllAck, err error) {
	defer func() {
		l.logger.Debug(
			fmt.Sprint(ctx.Value(ContextReqUUid)),
			zap.Any("调用 Service logMiddlewareServer", "GetAllPersonsInfo"),
			zap.Any("req", in),
			zap.Any("res", out),
		)
	}()
	out, err = l.next.GetAllPersonsInfo(ctx, in)
	return
}

func (l logMiddlewareServer) SearchPersonsInfoByName(ctx context.Context, in SearchPersons) (out SearchPersonsAck, err error) {
	defer func() {
		l.logger.Debug(
			fmt.Sprint(ctx.Value(ContextReqUUid)),
			zap.Any("调用 Service logMiddlewareServer", "SearchPersonsInfoByName"),
			zap.Any("req", in),
			zap.Any("res", out),
		)
	}()
	out, err = l.next.SearchPersonsInfoByName(ctx, in)
	return
}

func (l logMiddlewareServer) RegisterAccount(ctx context.Context, in User) (out UserAck, err error) {
	defer func() {
		l.logger.Debug(
			fmt.Sprint(ctx.Value(ContextReqUUid)),
			zap.Any("调用 Service logMiddlewareServer", "RegisterAccount"),
			zap.Any("req", in),
			zap.Any("res", out),
		)
	}()
	out, err = l.next.RegisterAccount(ctx, in)
	return
}
