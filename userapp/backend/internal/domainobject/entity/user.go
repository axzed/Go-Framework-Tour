package entity

import (
	"errors"
	"regexp"
)


type User struct {
	Id uint64
	Name string
	Avatar string
	Email string
	Password string
	Salt string
}

func (u User) Check() error {
	const emailPattern = "(.+)@(.+){2,}\\.(.+){2,}"
	if len(u.Password) < 8 {
		return errors.New("密码长度过短")
	}

	ok, err := regexp.Match(emailPattern, []byte(u.Email))
	if !ok || err !=nil {
		return errors.New("不正确的邮箱地址")
	}
	return nil
}
