package redisdb

import (
	"bigSystem/svc/common/utils"
	"fmt"
	"github.com/go-redis/redis"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
)

type RConfig struct {
	RedisServer struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Password string `yaml:"password"`
		db       int    `yaml:"db"`
	}
}

type RedisDrivers struct {
	RClient *redis.Client
	Conf    RConfig
}

var RdsClient *redis.Client
var conf RConfig

// Init 初始化
func Init(rConf string) {
	conf = RConfig{}
	confFile, err := ioutil.ReadFile(rConf)
	utils.CheckErr(err)

	err = yaml.Unmarshal(confFile, &conf)
	utils.CheckErr(err)

	//uri := fmt.Sprintf("mongodb://%s:%s@%s:%s", conf.Mongodb.Username, conf.Mongodb.Password, conf.Mongodb.Host, conf.Mongodb.Port)
	//fmt.Println(uri)
	RdsClient = Connect(conf)
	//MgoDbName = conf.Mongodb.Database
}

func Connect(conf RConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", conf.RedisServer.Host, conf.RedisServer.Port),
		Password: conf.RedisServer.Password,
		DB:       conf.RedisServer.db,
	})
	// 通过 client.Ping() 来检查是否成功连接到了 redis 服务器
	_, err := client.Ping().Result()
	utils.CheckErr(err)
	if err != nil {
		//fmt.Println("redis 连接失败")
		utils.GetLogger().Info("Failed to Redis!")
		return nil
	}
	//fmt.Println(pong, "redis 连接成功！！！")

	//fmt.Println("Connected to MongoDB!")
	utils.GetLogger().Info("Connected to Redis!")
	return client
}
