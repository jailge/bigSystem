package db

import (
	"bigSystem/svc/common/entity"
	"github.com/go-xorm/xorm"
	"time"
)

// InsertUser 通过ok
func InsertUser(e *xorm.Engine, userName string, password string, email string) error {
	//fmt.Println(userName)
	today := time.Now().Format("2006-01-02 15:04:05")
	user := new(entity.User)
	user.UserName = userName
	user.Password = password
	user.CreateAt = today
	user.Email = email

	_, err := e.Insert(user)

	if err != nil {
		return err
	}
	return nil
}

func CheckUserExists(e *xorm.Engine, userName string) (bool, error) {
	has, err := e.Exist(&entity.User{
		UserName: userName,
	})
	//fmt.Println(has) //true
	return has, err
}
