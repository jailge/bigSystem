package mongodb

import (
	"bigSystem/svc/common/entity"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/bson"
)

//func InsertOneRecord(cli *qmgo.QmgoClient) error {
//
//}

func (mSvc *MdbService) FindAllDocument(cli *qmgo.QmgoClient) ([]*entity.WeightRecord, error) {
	var batch []*entity.WeightRecord
	err := cli.Find(mSvc.mCtx, bson.M{}).All(&batch)
	if err != nil {
		return batch, err
	}
	//cli.
	return batch, nil
}
