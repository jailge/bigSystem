package mongodb

import (
	"context"
	"fmt"
	"github.com/qiniu/qmgo"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Mongodb struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	}
}

type MdbService struct {
	cli  *qmgo.QmgoClient
	mCtx context.Context
	conf Config
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func (mSvc *MdbService) InitMongo(dbConf string) error {
	confFile, err := ioutil.ReadFile(dbConf)
	checkErr(err)

	err = yaml.Unmarshal(confFile, &mSvc.conf)
	checkErr(err)

	uri := fmt.Sprintf("mongodb://%s:%s", mSvc.conf.Mongodb.Host, mSvc.conf.Mongodb.Port)

	mCtx := context.Background()
	auth := qmgo.Credential{
		Username: mSvc.conf.Mongodb.Username,
		Password: mSvc.conf.Mongodb.Password,
	}
	mSvc.cli, err = qmgo.Open(mCtx, &qmgo.Config{
		Uri:      uri,
		Database: mSvc.conf.Mongodb.Database,
		Coll:     "record",
		Auth:     &auth,
	})

	defer func() {
		if err = mSvc.cli.Close(mCtx); err != nil {
			checkErr(err)
		}
	}()
	return err

}

func (mSvc *MdbService) Cli() *qmgo.QmgoClient {
	//context.Background()
	return mSvc.cli
}
