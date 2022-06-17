package db

//
//import (
//	"fmt"
//	"github.com/go-redis/redis"
//	_ "github.com/go-sql-driver/mysql"
//	"github.com/go-xorm/xorm"
//	"log"
//	"os"
//	"time"
//	yaml "gopkg.in/yaml.v2"
//	"io/ioutil"
//	"github.com/sirupsen/logrus"
//	"path"
//)
//
//
//type Config struct {
//	DbServer struct{
//		Host string `yaml:"host"`
//		Port string `yaml:"port"`
//		Username string `yaml:"username"`
//		Password string `yaml:"password"`
//		Schema string `yaml:"schema"`
//	}
//
//	RedisServer struct{
//		Host string `yaml:"host"`
//		Port string `yaml:"port"`
//		Password string `yaml:"password"`
//		db int `yaml:"db"`
//	}
//}
//
//type DbService struct {
//	conf   Config
//	engine *xorm.Engine
//}
//
//
//type RedisService struct {
//	conf Config
//	client *redis.Client
//}
//
//
//func checkErr(err error)  {
//	if err != nil {
//		log.Println(err)
//	}
//}
//
//
//func (dbsvc *DbService) Bind(dbconf string) error {
//	confFile, err := ioutil.ReadFile(dbconf)
//	checkErr(err)
//
//	err = yaml.Unmarshal(confFile, &dbsvc.conf)
//	checkErr(err)
//	dataSourceName := dbsvc.conf.DbServer.Username + ":" + dbsvc.conf.DbServer.Password + "@tcp(" + dbsvc.conf.DbServer.Host + ":" + dbsvc.conf.DbServer.Port + ")/" + dbsvc.conf.DbServer.Schema + "?charset=utf8"
//	dbsvc.engine, err = xorm.NewEngine("mysql", dataSourceName)
//	dbsvc.engine.ShowSQL(true)
//	//dbsvc.engine.SetLogger(Logger())
//	checkErr(err)
//	if err != nil {
//		fmt.Println("Mysql 连接失败")
//	}
//	fmt.Println("Mysql 连接成功！！！")
//	return err
//}
//
//func (dbsvc *DbService) Engine() *xorm.Engine {
//	return dbsvc.engine
//}
//
//
//func Logger() *logrus.Logger {
//	now := time.Now()
//	logFilePath := ""
//	if dir, err := os.Getwd(); err == nil {
//		logFilePath = dir + "/logs/"
//	}
//	if err := os.MkdirAll(logFilePath, 0777); err != nil {
//		fmt.Println(err.Error())
//	}
//	logFileName := now.Format("2022-05-25") + ".log"
//	//日志文件
//	fileName := path.Join(logFilePath, logFileName)
//	if _, err := os.Stat(fileName); err != nil {
//		if _, err := os.Create(fileName); err != nil {
//			fmt.Println(err.Error())
//		}
//	}
//	//写入文件
//	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
//	if err != nil {
//		fmt.Println("err", err)
//	}
//
//	//实例化
//	logger := logrus.New()
//
//	//设置输出
//	logger.Out = src
//
//	//设置日志级别
//	logger.SetLevel(logrus.DebugLevel)
//
//	//设置日志格式
//	logger.SetFormatter(&logrus.TextFormatter{
//		TimestampFormat: "2006-01-02 15:04:05",
//	})
//	return logger
//}
