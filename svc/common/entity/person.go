package entity

type Person struct {
	Id           int    `json:"id"`
	PersonId     string `json:"person_id"`
	PersonName   string `json:"person_name"`
	Gender       string `json:"gender"`
	DeptName     string `json:"dept_name"`
	DeptId       string `json:"dept_id"`
	Automobile   string `json:"automobile"`
	Email        string `json:"email"`
	Sn           string `json:"sn"`
	WeixinId     string `json:"weixin_id"`
	WeixinDeptid string `json:"weixin_deptid"`
	Position     string `json:"position"`
	LoginName    string `json:"login_name"`
	Status       string `json:"status"`
	UpdateTime   string `json:"update_time"`
}

//type Person struct {
//	Id int `json:"id" db:"id"`
//	PersonId sql.NullString `json:"person_id" db:"person_id"`
//	PersonName sql.NullString `json:"person_name" db:"person_name"`
//	Gender sql.NullString `json:"gender" db:"gender"`
//	DeptName sql.NullString `json:"dept_name" db:"dept_name"`
//	DeptId sql.NullString `json:"dept_id" db:"dept_id"`
//	Automobile sql.NullString `json:"automobile" db:"automobile"`
//	Email sql.NullString `json:"email" db:"email"`
//	Sn string `json:"sn" db:"sn"`
//	WeixinId sql.NullString `json:"weixin_id" db:"weixin_id"`
//	Weixin_DeptId sql.NullString `json:"weixin_deptid" db:"weixin_deptid"`
//	Position sql.NullString `json:"position" db:"position"`
//	LoginName sql.NullString `json:"login_name" db:"login_name"`
//	Status sql.NullString `json:"status" db:"status"`
//	UpdateTime sql.NullString `json:"update_time" db:"update_time"`
//}

//type Person struct {
//	Id int `json:"id" gorm:"column:id"`
//	PersonId string `json:"person_id" gorm:"column:person_id"`
//	PersonName string `json:"person_name" gorm:"column:person_name"`
//	Gender string `json:"gender" gorm:"column:gender"`
//	DeptName string `json:"dept_name" gorm:"column:dept_name"`
//	DeptId string `json:"dept_id" gorm:"column:dept_id"`
//	Automobile string `json:"automobile" gorm:"column:automobile"`
//	Email string `json:"email" gorm:"column:email"`
//	Sn string `json:"sn" gorm:"column:sn"`
//	WeixinId string `json:"weixin_id" gorm:"column:weixin_id"`
//	Weixin_DeptId string `json:"weixin_deptid" gorm:"column:weixin_deptid"`
//	Position string `json:"position" gorm:"column:position"`
//	LoginName string `json:"login_name" gorm:"column:login_name"`
//	Status string `json:"status" gorm:"column:status"`
//	UpdateTime string `json:"update_time" gorm:"column:update_time"`
//}

func (Person) TableName() string {
	return "person"
}
