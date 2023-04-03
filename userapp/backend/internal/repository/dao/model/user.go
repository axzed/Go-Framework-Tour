package model

type User struct {
	Id uint64
	Name string
	Avatar string
	Email string
	Password string
	Salt string
	CreateTime uint64  // time second
	UpdateTime uint64 // time second
}

// type UserExtend struct {
// 	Phone string
// }

// func (usr *User) ToPB() *dto.User {
// 	return &dto.User{
// 		Id: usr.Id,
// 		Name: usr.Name,
// 		Avatar: usr.Avatar,
// 		Email: usr.Email,
// 		CreateTime: usr.CreateTime,
// 	}
// }
//
// func (usr *User) ToPBWithSensitive() *dto.User {
// 	return &dto.User{
// 		Id: usr.Id,
// 		Name: usr.Name,
// 		Avatar: usr.Avatar,
// 		Email: usr.Email,
// 		CreateTime: usr.CreateTime,
// 		Password: usr.Password,
// 	}
// }

