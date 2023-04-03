package service

import (
	"errors"
	"gitee.com/geektime-geekbang/geektime-go/userapp/backend/internal/repository"
)

var (
	ErrInvalidNewUser = errors.New("新用户数据错误")
	ErrInvalidUserOrPassword = errors.New("错误的登录信息")
	ErrDuplicateEmail = repository.ErrDuplicateEmail
)