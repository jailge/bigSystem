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
	"github.com/rs/cors"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"net/http"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrorBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
	ErrorBadRequest = errors.New("invalid request parameter")
)

func MakeHTTPHandler(endpoint Endpoint.EndpointsServer, log *zap.Logger) http.Handler {
	r := mux.NewRouter()
	//e := Endpoint.MakeServerEndpoints(s, log)
	//r.Use(mux.CORSMethodMiddleware(r))
	//r.Headers("Content-Type", "application/json",
	//	"X-Requested-With", "XMLHttpRequest", "Access-Control-Allow-Origin")

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
	r.Methods("POST").Path("/weight/material_code").Handler(
		handlers.CORS(handlers.AllowedMethods([]string{"POST"}), handlers.AllowedOrigins([]string{"*"}), handlers.AllowedHeaders([]string{"*"}))(
			httptransport.NewServer(
				endpoint.SearchWeightWithMaterialCodeEndpoint,
				decodeHTTPSearchWeightWithMaterialCodeRequest,
				encodeHTTPGenericResponse,
				options...,
			),
		))

	r.Methods("GET").Path("/weight/craft").Handler(
		httptransport.NewServer(
			endpoint.GetAllCraftEndpoint,
			decodeHTTPNoParameterRequest,
			encodeHTTPGenericResponse,
			options...,
		),
	)
	r.Methods("POST").Path("/weight/craft").Handler(
		httptransport.NewServer(
			endpoint.AddCraftEndpoint,
			decodeHTTPAddCraftRequest,
			encodeHTTPGenericResponse,
			options...,
		),
	)
	r.Methods("DELETE").Path("/weight/craft/{id}").Handler(
		httptransport.NewServer(
			endpoint.DeleteCraftWithIdEndpoint,
			decodeHTTPDeleteParameterWithIdRequest,
			encodeHTTPGenericResponse,
			options...,
		),
	)
	r.Methods("PUT").Path("/weight/craft/{id}").Handler(
		httptransport.NewServer(
			endpoint.UpdateCraftEndpoint,
			decodeHTTPUpdateCraftRequest,
			encodeHTTPGenericResponse,
			options...,
		),
	)

	r.Methods("GET").Path("/weight/texture").Handler(
		httptransport.NewServer(
			endpoint.GetAllTextureEndpoint,
			decodeHTTPNoParameterRequest,
			encodeHTTPGenericResponse,
			options...,
		),
	)
	r.Methods("POST").Path("/weight/texture").Handler(
		httptransport.NewServer(
			endpoint.AddTextureEndpoint,
			decodeHTTPAddTextureRequest,
			encodeHTTPGenericResponse,
			options...,
		),
	)
	r.Methods("DELETE").Path("/weight/texture/{id}").Handler(
		httptransport.NewServer(
			endpoint.DeleteTextureWithIdEndpoint,
			decodeHTTPDeleteParameterWithIdRequest,
			encodeHTTPGenericResponse,
			options...,
		),
	)
	r.Methods("PUT").Path("/weight/texture/{id}").Handler(
		httptransport.NewServer(
			endpoint.UpdateTextureEndpoint,
			decodeHTTPUpdateTextureRequest,
			encodeHTTPGenericResponse,
			options...,
		),
	)

	r.Methods("GET").Path("/weight/process").Handler(
		httptransport.NewServer(
			endpoint.GetAllProcessEndpoint,
			decodeHTTPNoParameterRequest,
			encodeHTTPGenericResponse,
			options...,
		),
	)
	r.Methods("POST").Path("/weight/process").Handler(
		httptransport.NewServer(
			endpoint.AddProcessEndpoint,
			decodeHTTPAddProcessRequest,
			encodeHTTPGenericResponse,
			options...,
		),
	)
	r.Methods("DELETE").Path("/weight/process/{id}").Handler(
		httptransport.NewServer(
			endpoint.DeleteProcessWithIdEndpoint,
			decodeHTTPDeleteParameterWithIdRequest,
			encodeHTTPGenericResponse,
			options...,
		),
	)
	r.Methods("PUT").Path("/weight/process/{id}").Handler(
		httptransport.NewServer(
			endpoint.UpdateProcessEndpoint,
			decodeHTTPUpdateProcessRequest,
			encodeHTTPGenericResponse,
			options...,
		),
	)

	r.Methods("GET").Path("/weight/purchase_status").Handler(
		httptransport.NewServer(
			endpoint.GetAllPurchaseStatusEndpoint,
			decodeHTTPNoParameterRequest,
			encodeHTTPGenericResponse,
			options...,
		),
	)
	r.Methods("POST").Path("/weight/purchase_status").Handler(
		httptransport.NewServer(
			endpoint.AddPurchaseStatusEndpoint,
			decodeHTTPAddPurchaseStatusRequest,
			encodeHTTPGenericResponse,
			options...,
		),
	)
	r.Methods("DELETE").Path("/weight/purchase_status/{id}").Handler(
		httptransport.NewServer(
			endpoint.DeletePurchaseStatusWithIdEndpoint,
			decodeHTTPDeleteParameterWithIdRequest,
			encodeHTTPGenericResponse,
			options...,
		),
	)
	r.Methods("PUT").Path("/weight/purchase_status/{id}").Handler(
		httptransport.NewServer(
			endpoint.UpdatePurchaseStatusEndpoint,
			decodeHTTPUpdatePurchaseStatusRequest,
			encodeHTTPGenericResponse,
			options...,
		),
	)

	r.Use(mux.CORSMethodMiddleware(r))
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://10.10.181.60:5500", "http://gm.he.link", "http://10.10.181.111", "http://10.10.6.51"},
		//AllowedOrigins:         []string{"*"},
		AllowOriginFunc:        nil,
		AllowOriginRequestFunc: nil,
		AllowedMethods:         []string{"POST", "GET", "PUT", "DELETE"},
		AllowedHeaders:         []string{"Content-Type", "application/json"},
		ExposedHeaders:         nil,
		MaxAge:                 0,
		AllowCredentials:       false,
		OptionsPassthrough:     false,
		OptionsSuccessStatus:   0,
		Debug:                  false,
	})
	cors.Default()
	h := c.Handler(r)

	return h
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
	fmt.Printf("%s\n", r.Header)
	fmt.Println("**********")
	fmt.Printf("%s\n", r.Body)
	var in Service.NewRecord
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return nil, err
	}
	utils.GetLogger().Debug(fmt.Sprint(ctx.Value(Service.ContextReqUUid)), zap.Any(" 开始解析请求数据", in))
	return in, nil
}

// ********************************
// 参数
// ********************************

func decodeHTTPNoParameterRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	utils.GetLogger().Debug(fmt.Sprint(ctx.Value(Service.ContextReqUUid)), zap.Any(" 开始解析请求数据", nil))
	return nil, nil
}

func decodeHTTPAddCraftRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var in Service.Craft
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return nil, err
	}
	utils.GetLogger().Debug(fmt.Sprint(ctx.Value(Service.ContextReqUUid)), zap.Any(" 开始解析请求数据", in))
	return in, nil
}

func decodeHTTPAddTextureRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var in Service.Texture
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return nil, err
	}
	utils.GetLogger().Debug(fmt.Sprint(ctx.Value(Service.ContextReqUUid)), zap.Any(" 开始解析请求数据", in))
	return in, nil
}

func decodeHTTPAddProcessRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var in Service.Process
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return nil, err
	}
	utils.GetLogger().Debug(fmt.Sprint(ctx.Value(Service.ContextReqUUid)), zap.Any(" 开始解析请求数据", in))
	return in, nil
}

func decodeHTTPAddPurchaseStatusRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var in Service.PurchaseStatus
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return nil, err
	}
	utils.GetLogger().Debug(fmt.Sprint(ctx.Value(Service.ContextReqUUid)), zap.Any(" 开始解析请求数据", in))
	return in, nil
}

func decodeHTTPUpdateCraftRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	var t Service.Craft
	id := vars["id"]
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return nil, err
	}
	in := Service.Craft{
		Id:       id,
		Name:     t.Name,
		ClientId: t.ClientId,
	}
	utils.GetLogger().Debug(fmt.Sprint(ctx.Value(Service.ContextReqUUid)), zap.Any(" 开始解析请求数据", in))
	return in, nil
}

func decodeHTTPUpdateTextureRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	var t Service.Texture
	id := vars["id"]
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return nil, err
	}
	in := Service.Texture{
		Id:       id,
		Name:     t.Name,
		ClientId: t.ClientId,
	}
	utils.GetLogger().Debug(fmt.Sprint(ctx.Value(Service.ContextReqUUid)), zap.Any(" 开始解析请求数据", in))
	return in, nil
}

func decodeHTTPUpdateProcessRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	var t Service.Process
	id := vars["id"]
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return nil, err
	}
	in := Service.Process{
		Id:       id,
		Name:     t.Name,
		ClientId: t.ClientId,
	}
	utils.GetLogger().Debug(fmt.Sprint(ctx.Value(Service.ContextReqUUid)), zap.Any(" 开始解析请求数据", in))
	return in, nil
}

func decodeHTTPUpdatePurchaseStatusRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	var t Service.PurchaseStatus
	id := vars["id"]
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return nil, err
	}
	in := Service.PurchaseStatus{
		Id:       id,
		Name:     t.Name,
		ClientId: t.ClientId,
	}
	utils.GetLogger().Debug(fmt.Sprint(ctx.Value(Service.ContextReqUUid)), zap.Any(" 开始解析请求数据", in))
	return in, nil
}

//func decodeHTTPDeleteCraftWithIdRequest(ctx context.Context, r *http.Request) (interface{}, error) {
//	//var in Service.CraftId
//	vars := mux.Vars(r)
//	//if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
//	//	return nil, err
//	//}
//	id := vars["id"]
//	utils.GetLogger().Debug(fmt.Sprint(ctx.Value(Service.ContextReqUUid)), zap.Any(" 开始解析请求数据", id))
//	return id, nil
//}
//
//func decodeHTTPDeleteTextureWithIdRequest(ctx context.Context, r *http.Request) (interface{}, error) {
//	//var in Service.CraftId
//	vars := mux.Vars(r)
//	id := vars["id"]
//	utils.GetLogger().Debug(fmt.Sprint(ctx.Value(Service.ContextReqUUid)), zap.Any(" 开始解析请求数据", id))
//	return id, nil
//}
//
//func decodeHTTPDeleteProcessWithIdRequest(ctx context.Context, r *http.Request) (interface{}, error) {
//	//var in Service.CraftId
//	vars := mux.Vars(r)
//	id := vars["id"]
//	utils.GetLogger().Debug(fmt.Sprint(ctx.Value(Service.ContextReqUUid)), zap.Any(" 开始解析请求数据", id))
//	return id, nil
//}

func decodeHTTPDeleteParameterWithIdRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	utils.GetLogger().Debug(fmt.Sprint(ctx.Value(Service.ContextReqUUid)), zap.Any(" 开始解析请求数据", id))
	return id, nil
}

func decodeHTTPSearchWeightWithMaterialCodeRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var in Service.MaterialCode
	//vars := mux.Vars(r)
	//in, err := vars["material_code"]
	//fmt.Printf(in)
	//if err {
	//	return nil, ErrorBadRequest
	//}
	fmt.Printf("%s\n", r.Header)
	fmt.Println("**********")
	//fmt.Printf("%s\n", r.Body)
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return nil, err
	}
	fmt.Println(in)
	utils.GetLogger().Debug(fmt.Sprint(ctx.Value(Service.ContextReqUUid)), zap.Any(" 开始解析请求数据", in))
	return in, nil
}

func encodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	utils.GetLogger().Debug(fmt.Sprint(ctx.Value(Service.ContextReqUUid)), zap.Any("请求结束封装返回值", response))
	if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
		encodeError(ctx, f.Failed(), w)
		return nil
	}
	//f, _ := response.(endpoint.Failer)
	//fmt.Sprintf("f.failed() %x ", msg.Failed())
	//w.WriteHeader(codeFrom(f.Failed()))
	w.Header().Set("Content-Type", "application/json")
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   false,
		"res":      "",
		"err_info": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err.Error() {
	case Service.NoErr.Error():
		return http.StatusOK
	case Service.ErrNotFound.Error():
		return http.StatusOK
	case Service.ErrAlreadyExists.Error(), Service.ErrInconsistentIDs.Error(), Service.NoParameters.Error():
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
