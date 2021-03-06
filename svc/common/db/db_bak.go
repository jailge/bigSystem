package db

import (
	"fmt"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	//"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

// 定义一个全局对象db
//var db *gorm.DB
var db *sqlx.DB
var redisClient *redis.Client

func Init() {
	initDB()
	initRedis()
}

// 定义一个初始化数据库的函数
//func initDB() (err error) {
//	// DSN:Data Source Name
//	db, err = gorm.Open("mysql", "root:root@(127.0.0.1:3306)/mdm?charset=utf8mb4&parseTime=True&loc=Local")
//	db.LogMode(true)
//	db.SetLogger(Logger())
//	if err != nil {
//		panic(err)
//	}
//	db.SingularTable(true)
//	fmt.Println(db)
//	return nil
//}
func initDB() (err error) {
	database, err := sqlx.Open("mysql", "root:root@tcp(127.0.0.1:3306)/mdm")
	if err != nil {
		fmt.Println("open mysql failed,", err)
		return
	}

	db = database
	//defer db.Close()  // 注意这行代码要写在上面err判断的下面
	return nil
}

func initRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	// 通过 cient.Ping() 来检查是否成功连接到了 redis 服务器
	pong, err := redisClient.Ping().Result()
	if err != nil {
		fmt.Println("redis 连接失败。。。。。")
	}
	fmt.Println(pong, "redis 连接成功！！！")
}

func Close() {
	db.Close()
}

//func GetInstance() *gorm.DB{
//	return db
//}

func GetInstance() *sqlx.DB {
	return db
}

func Logger() *logrus.Logger {
	now := time.Now()
	logFilePath := ""
	if dir, err := os.Getwd(); err == nil {
		logFilePath = dir + "/logs/"
	}
	if err := os.MkdirAll(logFilePath, 0777); err != nil {
		fmt.Println(err.Error())
	}
	logFileName := now.Format("2006-01-02") + "-db.log"
	//日志文件
	fileName := path.Join(logFilePath, logFileName)
	if _, err := os.Stat(fileName); err != nil {
		if _, err := os.Create(fileName); err != nil {
			fmt.Println(err.Error())
		}
	}
	//写入文件
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}

	//实例化
	logger := logrus.New()

	//设置输出
	logger.Out = src

	//设置日志级别
	logger.SetLevel(logrus.DebugLevel)

	//设置日志格式
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	return logger
}
