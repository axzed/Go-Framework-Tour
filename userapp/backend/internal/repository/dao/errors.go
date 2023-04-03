package dao

import (
	"errors"
	"gitee.com/geektime-geekbang/geektime-go/orm"
)

var (
	ErrDuplicateEmail = errors.New("dao: 邮件已经被注册过")
	ErrNoRows = orm.ErrNoRows
)
