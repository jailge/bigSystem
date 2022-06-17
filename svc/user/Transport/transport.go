package Transport

import (
	"bigSystem/svc/common/utils"
	"bigSystem/svc/user/Endpoint"
	"bigSystem/svc/user/Service"
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

	// POST
	r.Methods("POST").Path("/sum").Handler(
		handlers.CORS(handlers.AllowedMethods([]string{"POST"}), handlers.AllowedOrigins([]string{"*"}))(
			httptransport.NewServer(
				//e.AddEndpoint,
				endpoint.AddEndpoint,
				decodeHTTPAddRequest,      //解析请求值
				encodeHTTPGenericResponse, //返回值
				options...,
			)))
	r.Methods("POST").Path("/login").Handler(
		handlers.CORS(handlers.AllowedMethods([]string{"POST"}), handlers.AllowedOrigins([]string{"*"}))(
			httptransport.NewServer(
				endpoint.LoginEndpoint,
				decodeHTTPLoginRequest,
				encodeHTTPGenericResponse,
				options...,
			)))
	r.Methods("POST").Path("/person/sn").Handler(
		handlers.CORS(handlers.AllowedMethods([]string{"POST"}), handlers.AllowedOrigins([]string{"*"}))(
			httptransport.NewServer(
				endpoint.GetPersonInfoBySnEndPoint,
				decodeHTTPGetPersonInfoBySnRequest,
				encodeHTTPGenericResponse,
				options...,
			)))
	r.Methods("POST").Path("/persons/name").Handler(
		handlers.CORS(handlers.AllowedMethods([]string{"POST"}), handlers.AllowedOrigins([]string{"*"}))(
			httptransport.NewServer(
				endpoint.GetPersonsInfoByNameEndPoint,
				decodeHTTPGetPersonsInfoByNameRequest,
				encodeHTTPGenericResponse,
				options...,
			),
		))
	r.Methods("POST").Path("/persons/all").Handler(
		handlers.CORS(handlers.AllowedMethods([]string{"POST"}), handlers.AllowedOrigins([]string{"*"}))(
			httptransport.NewServer(
				endpoint.GetAllPersonsInfoEndPoint,
				decodeHTTPGetAllPersonsInfoRequest,
				encodeHTTPGenericResponse,
				options...,
			),
		),
	)
	r.Methods("POST").Path("/persons/search").Handler(
		handlers.CORS(handlers.AllowedMethods([]string{"POST"}), handlers.AllowedOrigins([]string{"*"}))(
			httptransport.NewServer(
				endpoint.SearchPersonsInfoByNameEndPoint,
				decodeHTTPSearchPersonsInfoByNameRequest,
				encodeHTTPGenericResponse,
				options...,
			),
		),
	)
	r.Methods("POST").Path("/user/register").Handler(
		handlers.CORS(handlers.AllowedMethods([]string{"POST"}), handlers.AllowedOrigins([]string{"*"}))(
			httptransport.NewServer(
				endpoint.RegisterAccountEndPoint,
				decodeHTTPRegisterAccountRequest,
				encodeHTTPGenericResponse,
				options...,
			)))

	return r
}

func decodeHTTPLoginRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var login Service.Login
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		return nil, err
	}
	utils.GetLogger().Debug(fmt.Sprint(ctx.Value(Service.ContextReqUUid)), zap.Any(" 开始解析请求数据", login))
	return login, nil

}

func decodeHTTPAddRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var (
		in  Service.Add
		err error
	)
	//in.A, err = strconv.Atoi(r.FormValue("a"))
	//in.B, err = strconv.Atoi(r.FormValue("b"))

	//var params map[string]int

	err = json.NewDecoder(r.Body).Decode(&in)
	if err != nil {
		return nil, err
	}

	//utils.GetLogger().Debug("", zap.Any("params", params))

	//in.A, err = strconv.Atoi(a)
	//in.B, err = strconv.Atoi(b)

	//in.A = params["a"]
	//in.B = params["b"]
	//if err != nil {
	//	return in, err
	//}

	utils.GetLogger().Debug(fmt.Sprint(ctx.Value(Service.ContextReqUUid)), zap.Any(" 开始解析请求数据", in))
	return in, nil

}

func decodeHTTPRegisterAccountRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var in Service.User
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return nil, err
	}
	utils.GetLogger().Debug(fmt.Sprint(ctx.Value(Service.ContextReqUUid)), zap.Any(" 开始解析请求数据", in))
	return in, nil
}

func decodeHTTPGetPersonInfoBySnRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var in Service.PersonSn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return nil, err
	}
	utils.GetLogger().Debug(fmt.Sprint(ctx.Value(Service.ContextReqUUid)), zap.Any(" 开始解析请求数据", in))
	return in, nil
}

func decodeHTTPGetPersonsInfoByNameRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var in Service.PersonName
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return nil, err
	}
	utils.GetLogger().Debug(fmt.Sprint(ctx.Value(Service.ContextReqUUid)), zap.Any(" 开始解析请求数据", in))
	return in, nil
}

func decodeHTTPGetAllPersonsInfoRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var in Service.AllPerson
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return nil, err
	}
	utils.GetLogger().Debug(fmt.Sprint(ctx.Value(Service.ContextReqUUid)), zap.Any("开始解析请求数据", in))
	return in, nil
}

func decodeHTTPSearchPersonsInfoByNameRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var in Service.SearchPersons
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return nil, err
	}
	utils.GetLogger().Debug(fmt.Sprint(ctx.Value(Service.ContextReqUUid)), zap.Any("开始解析请求数据", in))
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
	//fmt.Println(err)
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
	case Service.ErrUserLogin:
		return http.StatusBadRequest
	case Service.ErrNotFound:
		return http.StatusNotFound
	case Service.ErrAlreadyExists, Service.ErrInconsistentIDs:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
