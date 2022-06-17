package db

import (
	"fmt"
	//"fmt"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	DbServer struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Schema   string `yaml:"schema"`
	}

	RedisServer struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Password string `yaml:"password"`
		db       int    `yaml:"db"`
	}
}

//type dbConfig struct {
//	Host string `yaml:"host"`
//	Port string `yaml:"port"`
//	Username string `yaml:"username"`
//	Password string `yaml:"password"`
//	Schema string `yaml:"schema"`
//}

type DbService struct {
	conf   Config
	engine *xorm.Engine
}

//type redisConfig struct {
//	Host string `yaml:"host"`
//	Port string `yaml:"port"`
//	Password string `yaml:"password"`
//	db string `yaml:"db"`
//}

type RedisService struct {
	conf   Config
	client *redis.Client
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func (dbsvc *DbService) Bind(dbconf string) error {
	//dir, _ := os.Getwd()
	//fmt.Println(dir)
	//dbconf = dir + "/svc/common/config/conf.db.yml"
	confFile, err := ioutil.ReadFile(dbconf)
	checkErr(err)

	err = yaml.Unmarshal(confFile, &dbsvc.conf)
	checkErr(err)
	dataSourceName := dbsvc.conf.DbServer.Username + ":" + dbsvc.conf.DbServer.Password + "@tcp(" + dbsvc.conf.DbServer.Host + ":" + dbsvc.conf.DbServer.Port + ")/" + dbsvc.conf.DbServer.Schema + "?charset=utf8"
	dbsvc.engine, err = xorm.NewEngine("mysql", dataSourceName)
	dbsvc.engine.ShowSQL(true)
	checkErr(err)
	if err != nil {
		fmt.Println("Mysql 连接失败")
	}
	fmt.Println("Mysql 连接成功！！！")
	return err
}

func (dbsvc *DbService) Engine() *xorm.Engine {
	return dbsvc.engine
}

func (rdsSvc *RedisService) RdsBind(rdsconfPath string) error {
	confFile, err := ioutil.ReadFile(rdsconfPath)
	checkErr(err)

	err = yaml.Unmarshal(confFile, &rdsSvc.conf)
	checkErr(err)

	rdsSvc.client = redis.NewClient(&redis.Options{
		Addr:     rdsSvc.conf.RedisServer.Host,
		Password: "",
		DB:       rdsSvc.conf.RedisServer.db,
	})
	// 通过 cient.Ping() 来检查是否成功连接到了 redis 服务器
	pong, err := redisClient.Ping().Result()
	checkErr(err)
	if err != nil {
		fmt.Println("redis 连接失败")
	}
	fmt.Println(pong, "redis 连接成功！！！")
	return err
}

func (rdsSvc *RedisService) RedisClient() *redis.Client {
	return rdsSvc.client
}

//func (rdsSvc *RedisService) RedisClient() *redis.Client {
//	rdsSvc.client = redis.NewClient(&redis.Options{
//		Addr:     "localhost:6379",
//		Password: "",
//		DB:       0,
//	})
//	// 通过 cient.Ping() 来检查是否成功连接到了 redis 服务器
//	pong, err := redisClient.Ping().Result()
//	if err != nil {
//		fmt.Println("redis 连接失败。。。。。")
//	}
//	fmt.Println(pong, "redis 连接成功！！！")
//}
