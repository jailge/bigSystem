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

func (l logMiddlewareServer) UpdateTexture(ctx context.Context, id string, in Texture) (out NewParameterAck, err error) {
	defer func() {
		l.logger.Debug(
			fmt.Sprint(ctx.Value(ContextReqUUid)),
			zap.Any("调用 UpdateTexture logMiddlewareServer", "UpdateTexture"),
			zap.Any("req1", id),
			zap.Any("req2", in),
			zap.Any("res", out),
			zap.Any("err", err))
	}()
	out, err = l.next.UpdateTexture(ctx, id, in)
	return
}

func (l logMiddlewareServer) UpdateProcess(ctx context.Context, id string, in Process) (out NewParameterAck, err error) {
	defer func() {
		l.logger.Debug(
			fmt.Sprint(ctx.Value(ContextReqUUid)),
			zap.Any("调用 UpdateProcess logMiddlewareServer", "UpdateProcess"),
			zap.Any("req1", id),
			zap.Any("req2", in),
			zap.Any("res", out),
			zap.Any("err", err))
	}()
	out, err = l.next.UpdateProcess(ctx, id, in)
	return
}

func (l logMiddlewareServer) UpdatePurchaseStatus(ctx context.Context, id string, in PurchaseStatus) (out NewParameterAck, err error) {
	defer func() {
		l.logger.Debug(
			fmt.Sprint(ctx.Value(ContextReqUUid)),
			zap.Any("调用 UpdatePurchaseStatus logMiddlewareServer", "UpdatePurchaseStatus"),
			zap.Any("req1", id),
			zap.Any("req2", in),
			zap.Any("res", out),
			zap.Any("err", err))
	}()
	out, err = l.next.UpdatePurchaseStatus(ctx, id, in)
	return
}

func (l logMiddlewareServer) UpdateCraft(ctx context.Context, id string, in Craft) (out NewParameterAck, err error) {
	defer func() {
		l.logger.Debug(
			fmt.Sprint(ctx.Value(ContextReqUUid)),
			zap.Any("调用 UpdateCraft logMiddlewareServer", "UpdateCraft"),
			zap.Any("req1", id),
			zap.Any("req2", in),
			zap.Any("res", out),
			zap.Any("err", err))
	}()
	out, err = l.next.UpdateCraft(ctx, id, in)
	return
}

func (l logMiddlewareServer) GetAllPurchaseStatus(ctx context.Context) (out AllPurchaseStatusAck, err error) {
	defer func() {
		l.logger.Debug(
			fmt.Sprint(ctx.Value(ContextReqUUid)),
			zap.Any("调用 GetAllPurchaseStatus logMiddlewareServer", "GetAllPurchaseStatus"),
			//zap.Any("req", in),
			zap.Any("res", out),
			zap.Any("err", err))
	}()
	out, err = l.next.GetAllPurchaseStatus(ctx)
	return
}

func (l logMiddlewareServer) GetAllTexture(ctx context.Context) (out AllTextureAck, err error) {
	defer func() {
		l.logger.Debug(
			fmt.Sprint(ctx.Value(ContextReqUUid)),
			zap.Any("调用 GetAllTexture logMiddlewareServer", "GetAllTexture"),
			//zap.Any("req", in),
			zap.Any("res", out),
			zap.Any("err", err))
	}()
	out, err = l.next.GetAllTexture(ctx)
	return
}

func (l logMiddlewareServer) GetAllProcess(ctx context.Context) (out AllProcessAck, err error) {
	defer func() {
		l.logger.Debug(
			fmt.Sprint(ctx.Value(ContextReqUUid)),
			zap.Any("调用 GetAllProcess logMiddlewareServer", "GetAllProcess"),
			//zap.Any("req", in),
			zap.Any("res", out),
			zap.Any("err", err))
	}()
	out, err = l.next.GetAllProcess(ctx)
	return
}

func (l logMiddlewareServer) GetAllCraft(ctx context.Context) (out AllCraftAck, err error) {
	defer func() {
		l.logger.Debug(
			fmt.Sprint(ctx.Value(ContextReqUUid)),
			zap.Any("调用 GetAllCraft logMiddlewareServer", "GetAllCraft"),
			//zap.Any("req", in),
			zap.Any("res", out),
			zap.Any("err", err))
	}()
	out, err = l.next.GetAllCraft(ctx)
	return
}

func (l logMiddlewareServer) DeleteTextureWithId(ctx context.Context, in string) (out NewParameterAck, err error) {
	defer func() {
		l.logger.Debug(
			fmt.Sprint(ctx.Value(ContextReqUUid)),
			zap.Any("调用 DeleteTextureWithId logMiddlewareServer", "DeleteTextureWithId"),
			zap.Any("req", in),
			zap.Any("res", out),
			zap.Any("err", err))
	}()
	out, err = l.next.DeleteTextureWithId(ctx, in)
	return
}

func (l logMiddlewareServer) DeleteProcessWithId(ctx context.Context, in string) (out NewParameterAck, err error) {
	defer func() {
		l.logger.Debug(
			fmt.Sprint(ctx.Value(ContextReqUUid)),
			zap.Any("调用 DeleteProcessWithId logMiddlewareServer", "DeleteProcessWithId"),
			zap.Any("req", in),
			zap.Any("res", out),
			zap.Any("err", err))
	}()
	out, err = l.next.DeleteProcessWithId(ctx, in)
	return
}

func (l logMiddlewareServer) DeleteCraftWithId(ctx context.Context, in string) (out NewParameterAck, err error) {
	defer func() {
		l.logger.Debug(
			fmt.Sprint(ctx.Value(ContextReqUUid)),
			zap.Any("调用 DeleteCraftWithId logMiddlewareServer", "DeleteCraftWithId"),
			zap.Any("req", in),
			zap.Any("res", out),
			zap.Any("err", err))
	}()
	out, err = l.next.DeleteCraftWithId(ctx, in)
	return
}

func (l logMiddlewareServer) DeletePurchaseStatusWithId(ctx context.Context, in string) (out NewParameterAck, err error) {
	defer func() {
		l.logger.Debug(
			fmt.Sprint(ctx.Value(ContextReqUUid)),
			zap.Any("调用 DeletePurchaseStatusWithId logMiddlewareServer", "DeletePurchaseStatusWithId"),
			zap.Any("req", in),
			zap.Any("res", out),
			zap.Any("err", err))
	}()
	out, err = l.next.DeletePurchaseStatusWithId(ctx, in)
	return
}

func (l logMiddlewareServer) AddProcess(ctx context.Context, in Process) (out NewParameterAck, err error) {
	defer func() {
		l.logger.Debug(
			fmt.Sprint(ctx.Value(ContextReqUUid)),
			zap.Any("调用 AddProcess logMiddlewareServer", "AddProcess"),
			zap.Any("req", in),
			zap.Any("res", out),
			zap.Any("err", err))
	}()
	out, err = l.next.AddProcess(ctx, in)
	return
}

func (l logMiddlewareServer) AddTexture(ctx context.Context, in Texture) (out NewParameterAck, err error) {
	defer func() {
		l.logger.Debug(
			fmt.Sprint(ctx.Value(ContextReqUUid)),
			zap.Any("调用 AddTexture logMiddlewareServer", "AddTexture"),
			zap.Any("req", in),
			zap.Any("res", out),
			zap.Any("err", err))
	}()
	out, err = l.next.AddTexture(ctx, in)
	return
}

func (l logMiddlewareServer) AddCraft(ctx context.Context, in Craft) (out NewParameterAck, err error) {
	defer func() {
		l.logger.Debug(
			fmt.Sprint(ctx.Value(ContextReqUUid)),
			zap.Any("调用 AddCraft logMiddlewareServer", "AddCraft"),
			zap.Any("req", in),
			zap.Any("res", out),
			zap.Any("err", err))
	}()
	out, err = l.next.AddCraft(ctx, in)
	return
}

func (l logMiddlewareServer) AddPurchaseStatus(ctx context.Context, in PurchaseStatus) (out NewParameterAck, err error) {
	defer func() {
		l.logger.Debug(
			fmt.Sprint(ctx.Value(ContextReqUUid)),
			zap.Any("调用 AddPurchaseStatus logMiddlewareServer", "AddPurchaseStatus"),
			zap.Any("req", in),
			zap.Any("res", out),
			zap.Any("err", err))
	}()
	out, err = l.next.AddPurchaseStatus(ctx, in)
	return
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
