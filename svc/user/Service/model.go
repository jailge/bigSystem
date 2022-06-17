package Service

import "bigSystem/svc/common/entity"

type User struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserAck struct {
	Status  bool         `json:"status"`
	Res     *entity.User `json:"res"`
	ErrInfo string       `json:"err_info"`
}

// PersonSn collects the request parameters for the GetPersonBySn method.
type PersonSn struct {
	PersonSn string `json:"person_sn"`
}

// PersonSnAck collects the response parameters for the GetPersonBySn method.
type PersonSnAck struct {
	Status  bool           `json:"status"`
	Res     *entity.Person `json:"res"`
	ErrInfo string         `json:"err_info"`
}

type PersonName struct {
	PersonName string `json:"person_name"`
}

// PersonsAck collects the response parameters for the GetPersonsByName method.
type PersonsAck struct {
	Status  bool            `json:"status"`
	Res     []entity.Person `json:"res"`
	ErrInfo string          `json:"err_info"`
}

// GetAllPersonsInfo

type AllPerson struct {
	PageNum  int `json:"page_num"`
	PageSize int `json:"page_size"`
}

type CountAndPerson struct {
	Persons []entity.Person `json:"persons"`
	Count   int             `json:"count"`
}

type PersonsAllAck struct {
	Status  bool           `json:"status"`
	Res     CountAndPerson `json:"res"`
	ErrInfo string         `json:"err_info"`
}

// SearchPersonsInfo

type SearchPersons struct {
	Name string `json:"name"`
}

type SearchPersonsAck struct {
	Status  bool             `json:"status"`
	Res     []*entity.Person `json:"res"`
	ErrInfo string           `json:"err_info"`
}

// Department

type Add struct {
	A int `json:"a"`
	B int `json:"b"`
}

type AddAck struct {
	Status bool   `json:"status"`
	Res    int    `json:"res"`
	Msg    string `json:"msg"`
}

type Login struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type LoginAck struct {
	Token string `json:"token"`
}
