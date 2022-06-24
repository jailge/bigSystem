package redisdb

import (
	"bigSystem/svc/common/utils"
	"time"
)

func GetLock(k string, v string) (bool, error) {
	res, err := RdsClient.SetNX(k, v, time.Minute).Result()
	if err != nil {
		utils.GetLogger().Info("GetLock error")
	}
	return res, err
}

func RunEvalDel(k string, v string) (int, error) {
	//执行Lua脚本
	//:param key: key
	//:param value: value
	//:return: 删除成功1，失败0
	res, err := RdsClient.Eval("if redis.call('get', KEYS[1])==ARGV[1] then return redis.call('del', KEYS[1]) else return 0 end", []string{k}, v).Result()
	if err != nil {
		utils.GetLogger().Info("RunEvalDel error")
	}
	r, _ := res.(int)
	return r, err
}

func Pttl(k string) (int64, error) {
	//获取key对应剩余过期时间
	//:param key:
	//:return: ms
	pttl, err := RdsClient.PTTL(k).Result()
	if err != nil {
		utils.GetLogger().Info("Pttl error")
	}
	return int64(pttl), err
}
