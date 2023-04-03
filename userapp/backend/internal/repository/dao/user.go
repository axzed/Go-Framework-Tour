package dao

import (
	"context"
	"gitee.com/geektime-geekbang/geektime-go/orm"
	"gitee.com/geektime-geekbang/geektime-go/userapp/backend/internal/repository/dao/model"
	"github.com/go-sql-driver/mysql"
	"github.com/opentracing/opentracing-go"
)

//go:generate mockgen -source=user.go -destination=mocks/user_mock.gen.go -package=daomocks UserDAO
type UserDAO interface {
	InsertUser(ctx context.Context, u *model.User) error
	UpdateUser(ctx context.Context, u *model.User) error
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserById(ctx context.Context, id uint64) (*model.User, error)
}

func NewUserDAO(sess orm.Session) UserDAO {
	return &userDAO{
		sess: sess,
	}
}

type userDAO struct {
	sess orm.Session
}

const operationNamePrefix = "dao."

func (dao *userDAO) UpdateUser(ctx context.Context, u *model.User) error {
	// 这里你可以考虑使用你在作业里面支持的，只更新非零值的那个特性
	// 我一般喜欢显式地写，也就是我在课堂上说的，更新非零值这个东西，看代码你是看不出来哪个字段是零值，哪个字段不是
	return orm.NewUpdater[model.User](dao.sess).Set(orm.Assign("Name", u.Name)).Exec(ctx).Err()
}

func (dao *userDAO) GetUserById(ctx context.Context, id uint64) (*model.User, error) {
	return orm.NewSelector[model.User](dao.sess).
		Where(orm.C("Id").EQ(id)).Get(ctx)
}

func (dao *userDAO) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return orm.NewSelector[model.User](dao.sess).
		Where(orm.C("Email").EQ(email)).Get(ctx)
}

func (dao *userDAO) InsertUser(ctx context.Context, u *model.User) error {
	name := operationNamePrefix + "InsertUser"
	span, _ := opentracing.StartSpanFromContext(ctx, name)
	defer span.Finish()
	err := orm.NewInserter[model.User](dao.sess).Values(u).Exec(ctx).Err()
	if err != nil {
		me, ok := err.(*mysql.MySQLError)
		if ok && me.Number == 1062 {
			return ErrDuplicateEmail
		}
	}
	return err
}