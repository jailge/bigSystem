package Service

import (
	"bigSystem/svc/common/db"
	"bigSystem/svc/common/entity"
	"bigSystem/svc/common/utils"
	"bigSystem/svc/user/pb"
	"context"
	"errors"
	"fmt"
	//"go-mdm-service/utils/async"
	"go.uber.org/zap"
	"strconv"
)

type Service interface {
	TestAdd(ctx context.Context, in Add) AddAck
	Login(ctx context.Context, in *pb.Login) (ack *pb.LoginAck, err error)
	//Login(ctx context.Context, in Login) (ack LoginAck, err error)
	LoginHTTP(ctx context.Context, in Login) (ack LoginAck, err error)

	GetPersonInfoBySn(ctx context.Context, in PersonSn) (ack PersonSnAck, err error)
	GetPersonsInfoByName(ctx context.Context, in PersonName) (ack PersonsAck, err error)
	GetAllPersonsInfo(ctx context.Context, in AllPerson) (ack PersonsAllAck, err error)
	SearchPersonsInfoByName(ctx context.Context, in SearchPersons) (ack SearchPersonsAck, err error)

	RegisterAccount(ctx context.Context, in User) (ack UserAck, err error)
}

// Repository new add
//type Repository interface {
//
//
//}

type baseServer struct {
	*db.DbService
	logger *zap.Logger
}

var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
	ErrUserLogin       = errors.New("用户信息错误")
)

func NewService(db *db.DbService, log *zap.Logger) Service {
	//func NewService(log *zap.Logger) Service {
	var server Service
	//dir, _ := os.Getwd()
	//configFile := dir + "/svc/common/config/conf.db.yml"
	//"./svc/common/config/conf.db.yml"
	if err := db.Bind("./svc/common/config/conf.db.yml"); err != nil {
		fmt.Println("The PersonService failed to bind with mysql")
	}
	server = &baseServer{db, log}
	//server = &baseServer{log}

	server = NewLogMiddlewareServer(log)(server)
	return server
}

func (s baseServer) TestAdd(ctx context.Context, in Add) AddAck {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "TestAdd 处理请求"))
	ack := AddAck{
		Status: true,
		Res:    in.A + in.B,
	}
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "TestAdd 处理请求"), zap.Any("处理返回值", ack))
	return ack
}

func (s baseServer) RegisterAccount(ctx context.Context, in User) (ack UserAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "RegisterAccount 处理请求"))
	exists, err := db.CheckUserExists(s.Engine(), in.UserName)
	if err != nil {
		return UserAck{}, err
	}
	if exists == true {
		return UserAck{
			Status:  false,
			Res:     nil,
			ErrInfo: "已存在该账号",
		}, err
	}
	password := utils.Md5Password(in.Password) // 密码md5加密
	err = db.InsertUser(s.Engine(), in.UserName, password, in.Email)
	if err != nil {
		return UserAck{
			Status:  false,
			Res:     nil,
			ErrInfo: err.Error(),
		}, err
	}
	ack = UserAck{
		Status: true,
		Res: &entity.User{
			UserName: in.UserName,
			Password: in.Password,
			Email:    in.Email,
		},
		ErrInfo: "",
	}
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "RegisterAccount 处理请求"), zap.Any("处理返回值", ack))
	return
}

func (s baseServer) GetPersonInfoBySn(ctx context.Context, in PersonSn) (ack PersonSnAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "GetPersonInfoBySn 处理请求"))
	data := &entity.Person{
		Sn: in.PersonSn,
	}
	var errInfo string
	status, err := s.Engine().Get(data)
	if status == false {
		data = nil
	}
	if err != nil {
		errInfo = err.Error()
	}

	//person, err := db.SelectPersonBySn(in.PersonSn)
	//fmt.Println(person)
	//if err != nil {
	//	return
	//}
	ack = PersonSnAck{
		Status:  status,
		Res:     data,
		ErrInfo: errInfo,
	}
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "GetPersonInfoBySn 处理请求"), zap.Any("处理返回值", ack))
	return
}

func (s baseServer) GetPersonsInfoByName(ctx context.Context, in PersonName) (ack PersonsAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "GetPersonsInfoByName 处理请求"))
	status := true
	res := make([]entity.Person, 0)
	var errInfo string
	err = s.Engine().Where("person_name=?", in.PersonName).Find(&res)
	if err != nil {
		status = false
		errInfo = err.Error()
	}
	//for _, u := range res {
	//	fmt.Println("person", u)
	//}
	ack = PersonsAck{
		Status:  status,
		Res:     res,
		ErrInfo: errInfo,
	}
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "GetPersonsInfoByName 处理请求"), zap.Any("处理返回值", ack))
	return
}

func (s baseServer) GetAllPersonsInfo(ctx context.Context, in AllPerson) (ack PersonsAllAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "GetAllPersonsInfo 处理请求"))
	status := true
	persons := make([]entity.Person, 0)
	cp := CountAndPerson{}
	// 分页
	var pageNum, pageSize = in.PageNum, in.PageSize
	var errInfo string
	if pageNum-1 <= 0 {
		pageNum = 1
	}
	err = s.Engine().Limit(pageSize*(pageNum), (pageNum-1)*pageSize).Find(&persons)
	if err != nil {
		status = false
		errInfo = err.Error()
	}
	cp.Persons = persons
	// 查询所有记录
	sql := "select count(*) count from person"
	countRes, err := s.Engine().Query(sql)

	if err != nil {
		status = false
		errInfo = err.Error()
	}

	countTmp := "0"
	for _, v := range countRes {
		countTmp = string(v["count"])
	}

	cp.Count, _ = strconv.Atoi(countTmp)
	ack = PersonsAllAck{
		Status:  status,
		Res:     cp,
		ErrInfo: errInfo,
	}
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "GetAllPersonsInfo 处理请求"), zap.Any("处理返回值", ack))
	return
}

// SearchPersonsInfoByName 模糊搜索名字
func (s baseServer) SearchPersonsInfoByName(ctx context.Context, in SearchPersons) (ack SearchPersonsAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "SearchPersonsInfoByName 处理请求"))
	status := true
	res := make([]*entity.Person, 0)
	var errInfo string
	// 模糊查询所有记录
	//sql := fmt.Sprintf("select * from person where person_name like '\%%s\%'")
	sql := "select * from person where person_name like '%" + in.Name + "%'"
	searchRes, err := s.Engine().Query(sql)

	if err != nil {
		status = false
		errInfo = err.Error()
	}

	for _, v := range searchRes {
		//fmt.Println(string(v["person_name"]))
		//countTmp = string(v["person_name"])
		idTmp, _ := strconv.Atoi(string(v["id"]))
		p := entity.Person{
			Id:           idTmp,
			PersonId:     string(v["person_id"]),
			PersonName:   string(v["person_name"]),
			Gender:       string(v["gender"]),
			DeptId:       string(v["dept_id"]),
			DeptName:     string(v["dept_name"]),
			Automobile:   string(v["automobile"]),
			Email:        string(v["email"]),
			Sn:           string(v["sn"]),
			WeixinId:     string(v["weixin_id"]),
			WeixinDeptid: string(v["weixin_deptid"]),
			Position:     string(v["position"]),
			LoginName:    string(v["login_name"]),
			Status:       string(v["status"]),
			UpdateTime:   string(v["update_time"]),
		}
		//fmt.Println(p)
		res = append(res, &p)
		//fmt.Println(res)
	}
	//fmt.Println(res)
	ack = SearchPersonsAck{
		Status:  status,
		Res:     res,
		ErrInfo: errInfo,
	}

	//go func() {
	//	err := async.SendSumTask()
	//	if err != nil {
	//		errInfo = err.Error()
	//	}
	//}()

	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "SearchPersonsInfoByName 处理请求"), zap.Any("处理返回值", ack))
	return
}

func (s baseServer) Login(ctx context.Context, in *pb.Login) (ack *pb.LoginAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "Login 处理请求"))
	if in.Account != "gfj" || in.Password != "123456" {
		err = errors.New("用户信息错误")
		return
	}
	ack = &pb.LoginAck{}
	ack.Token, err = utils.CreateJwtToken(in.Account, 1)
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "Login 处理请求"), zap.Any("处理返回值", ack))
	return
}

func (s baseServer) LoginHTTP(ctx context.Context, in Login) (ack LoginAck, err error) {
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "Login 处理请求"))
	if in.Account != "gfj" || in.Password != "123456" {
		err = errors.New("用户信息错误")
		return
	}

	ack.Token, err = utils.CreateJwtToken(in.Account, 1)
	s.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Service", "Login 处理请求"), zap.Any("处理返回值", ack))
	return
}
