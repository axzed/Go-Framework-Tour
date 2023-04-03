package main

type FindByUserIdReq struct {
	Id uint64
}

type FindByUserIdResp struct {
	User *User
}

type User struct {
	Id         uint64
	Name       string
	Avatar     string
	Email      string
	Password   string
	CreateTime int
}
