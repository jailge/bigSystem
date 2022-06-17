package Transport

import (
	"bigSystem/svc/common/utils"
	"bigSystem/svc/weight/Endpoint"
	"bigSystem/svc/weight/Service"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	handlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"net/http"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

func MakeHTTPHandler(endpoint Endpoint.EndpointsServer, log *zap.Logger) http.Handler {
	r := mux.NewRouter()
	//e := Endpoint.MakeServerEndpoints(s, log)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
		httptransport.ServerBefore(func(ctx context.Context, request *http.Request) context.Context {
			//UUID := uuid.NewV5(uuid.Must(uuid.NewV4()), "req_uuid").String()
			UUID := uuid.NewV4().String()
			log.Debug("给请求添加uuid", zap.Any("UUID", UUID))
			ctx = context.WithValue(ctx, Service.ContextReqUUid, UUID)
			ctx = context.WithValue(ctx, utils.JWT_CONTEXT_KEY, request.Header.Get("Authorization"))
			log.Debug("把请求中的token发到Context中", zap.Any("Token", request.Header.Get("Authorization")))
			return ctx
		}),
	}

	r.Methods("GET").Path("/weight/record").Handler(
		handlers.CORS(handlers.AllowedMethods([]string{"GET"}), handlers.AllowedOrigins([]string{"*"}))(
			httptransport.NewServer(
				endpoint.GetAllWeightRecordEndpoint,
				decodeHTTPGetAllWeightRecordRequest,
				encodeHTTPGenericResponse,
				options...,
			),
		))
	r.Methods("POST").Path("/weight/record_page").Handler(
		handlers.CORS(handlers.AllowedMethods([]string{"POST"}), handlers.AllowedOrigins([]string{"*"}))(
			httptransport.NewServer(
				endpoint.GetWeightRecordPageEndpoint,
				decodeHTTPGetWeightRecordPageRequest,
				encodeHTTPGenericResponse,
				options...,
			),
		))
	r.Methods("GET").Path("/weight/all_parameter").Handler(
		handlers.CORS(handlers.AllowedMethods([]string{"GET"}), handlers.AllowedOrigins([]string{"*"}))(
			httptransport.NewServer(
				endpoint.GetParameterEndpoint,
				decodeHTTPGetParameterRequest,
				encodeHTTPGenericResponse,
				options...,
			),
		))
	// 新增一条record
	r.Methods("POST").Path("/weight/records").Handler(
		handlers.CORS(handlers.AllowedMethods([]string{"POST"}), handlers.AllowedOrigins([]string{"*"}))(
			httptransport.NewServer(
				endpoint.AddNewRecordEndpoint,
				decodeHTTPAddNewRecordRequest,
				encodeHTTPGenericResponse,
				options...,
			),
		))

	return r
}

func decodeHTTPGetAllWeightRecordRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	//var in Service.
	//if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
	//	return nil, err
	//}
	utils.GetLogger().Debug(fmt.Sprint(ctx.Value(Service.ContextReqUUid)), zap.Any(" 开始解析请求数据", nil))
	return nil, nil
}

func decodeHTTPGetWeightRecordPageRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var in Service.Page
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return nil, err
	}
	utils.GetLogger().Debug(fmt.Sprint(ctx.Value(Service.ContextReqUUid)), zap.Any(" 开始解析请求数据", in))
	return in, nil
}

func decodeHTTPGetParameterRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	//var in Service.
	//if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
	//	return nil, err
	//}
	utils.GetLogger().Debug(fmt.Sprint(ctx.Value(Service.ContextReqUUid)), zap.Any(" 开始解析请求数据", nil))
	return nil, nil
}

func decodeHTTPAddNewRecordRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var in Service.NewRecord
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return nil, err
	}
	utils.GetLogger().Debug(fmt.Sprint(ctx.Value(Service.ContextReqUUid)), zap.Any(" 开始解析请求数据", in))
	return in, nil
}

func encodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	utils.GetLogger().Debug(fmt.Sprint(ctx.Value(Service.ContextReqUUid)), zap.Any("请求结束封装返回值", response))
	if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
		encodeError(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": false,
		"res":    "",
		"msg":    err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case Service.ErrNotFound:
		return http.StatusNotFound
	case Service.ErrAlreadyExists, Service.ErrInconsistentIDs:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
