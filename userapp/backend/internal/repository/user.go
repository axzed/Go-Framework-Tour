package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gitee.com/geektime-geekbang/geektime-go/cache"
	"gitee.com/geektime-geekbang/geektime-go/userapp/backend/internal/domainobject/entity"
	"gitee.com/geektime-geekbang/geektime-go/userapp/backend/internal/repository/dao"
	"gitee.com/geektime-geekbang/geektime-go/userapp/backend/internal/repository/dao/model"
	"go.uber.org/zap"
	"time"
)

const cacheByIdFormat = "user-%d"

//go:generate mockgen -source=user.go -destination=mocks/user_mock.gen.go -package=repomocks UserRepository
type UserRepository interface {
	CreateUser(ctx context.Context, user entity.User) (entity.User, error)
	UpdateUser(ctx context.Context, user entity.User) error
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	GetUserById(ctx context.Context, id uint64) (entity.User, error)
}

func NewUserRepository(dao dao.UserDAO, c cache.Cache) UserRepository {
	return &cacheUserRepository{
		dao: dao,
		cache: c,
	}
}

type cacheUserRepository struct {
	dao dao.UserDAO
	cache cache.Cache
}

func (c *cacheUserRepository) UpdateUser(ctx context.Context, user entity.User) error {
	err := c.dao.UpdateUser(ctx, &model.User{
		Name: user.Name,
		Email: user.Email,
	})
	if err != nil {
		return err
	}
	return c.cache.Delete(ctx, fmt.Sprintf(cacheByIdFormat, user.Id))
}

func (c *cacheUserRepository) GetUserById(ctx context.Context, id uint64) (entity.User, error) {
	u, err := c.getFromCache(ctx, id)
	if err == nil {
		return u, nil
	}

	mu, err := c.dao.GetUserById(ctx, id)
	if errors.Is(err, dao.ErrNoRows) {
		// 实际上，如果你们公司对安全和数据的要求很高，那么这里你并不能打印 email
		// 而是要打印 email 加了 * 号遮蔽一部分的
		return entity.User{}, fmt.Errorf("repo: %w, id %d", ErrUserNotFound, id)
	}
	if err != nil {
		return entity.User{}, err
	}
	return entity.User{
		Id: mu.Id,
		Name: mu.Name,
		Avatar: mu.Avatar,
		Email: mu.Email,
		Password: mu.Password,
		Salt: mu.Salt,
	}, nil
}

func (c *cacheUserRepository) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	// 理论上来说你也可以同时在缓存里面缓存 email 到用户信息的映射。
	// 但是这个方法只会被 GetUserByEmail 使用，它本身是一个低频行为，那么缓存与否就不是关键了
	u, err := c.dao.GetUserByEmail(ctx, email)
	if errors.Is(err, dao.ErrNoRows) {
		// 实际上，如果你们公司对安全和数据的要求很高，那么这里你并不能打印 email
		// 而是要打印 email 加了 * 号遮蔽一部分的
		return entity.User{}, fmt.Errorf("repo: %w, 邮箱 %s", ErrUserNotFound, email)
	}
	if err != nil {
		return entity.User{}, err
	}
	return entity.User{
		Id: u.Id,
		Name: u.Name,
		Avatar: u.Avatar,
		Email: u.Email,
		Password: u.Password,
		Salt: u.Salt,
	}, nil
}

// CreateUser 创建用户
// 如果邮箱已经存在，将会返回 ErrDuplicateEmail
func (c *cacheUserRepository) CreateUser(ctx context.Context, user entity.User) (entity.User, error) {
	now := uint64(time.Now().Unix())
	um := &model.User{
		Name:       user.Name,
		Avatar:     user.Avatar,
		Email:      user.Email,
		Salt: user.Salt,
		Password:   user.Password,
		CreateTime: now,
		UpdateTime: now,
	}
	err := c.dao.InsertUser(ctx, um)
	if err != nil {
		return entity.User{}, err
	}
	user.Id = um.Id
	c.cacheUser(ctx, um.Id, user)
	return user, nil
}

func (c *cacheUserRepository) getFromCache(ctx context.Context, id uint64) (entity.User, error) {
	key := fmt.Sprintf(cacheByIdFormat, id)
	val, err := c.cache.Get(ctx, key)
	if err != nil {
		return entity.User{}, err
	}
	var res entity.User
	err = json.Unmarshal([]byte(val.(string)), &res)
	return res, err
}

func (c *cacheUserRepository) cacheUser(ctx context.Context, id uint64, value entity.User) {
	key := fmt.Sprintf(cacheByIdFormat, id)
	val, err := json.Marshal(value)
	if err != nil {
		zap.L().Error("repo: 序列化 User 失败", zap.Error(err))
		return
	}
	if err = c.cache.Set(ctx, key, val, time.Hour); err != nil {
		zap.L().Error("repo: 缓存 user 失败", zap.Error(err))
	}
}
