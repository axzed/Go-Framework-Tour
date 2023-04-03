//go:build v4

package orm

import (
	"gitee.com/geektime-geekbang/geektime-go/orm/internal/errs"
	"reflect"
	"sync"
	"unicode"
)

// 这种包变量对测试不友好，缺乏隔离
// var defaultRegistry = &registry{
// 	models: make(map[reflect.Type]*model, 16),
// }

type registry struct {
	// models key 是类型名
	// 这种定义方式是不行的
	// 1. 类型名冲突，例如都是 User，但是一个映射过去 buyer_t
	// 一个映射过去 seller_t
	// 2. 并发不安全
	// models map[string]*model

	// lock sync.RWMutex
	// models map[reflect.Type]*model
	models sync.Map
}

// 使用 sync.Map
func (r *registry) get(val any) (*model, error) {
	typ := reflect.TypeOf(val)
	m, ok := r.models.Load(typ)
	if !ok {
		var err error
		if m, err = r.parseModel(typ); err != nil {
			return nil, err
		}
	}
	r.models.Store(typ, m)
	return m.(*model), nil
}

// 直接 map
// func (r *registry) get(val any) (*model, error) {
// 	typ := reflect.TypeOf(val)
// 	m, ok := r.models[typ]
// 	if !ok {
// 		var err error
// 		if m, err = r.parseModel(typ); err != nil {
// 			return nil, err
// 		}
// 	}
// 	r.models[typ] = m
// 	return m, nil
// }

// 使用读写锁的并发安全解决思路
// func (r *registry) get1(val any) (*model, error) {
// 	r.lock.RLock()
// 	typ := reflect.TypeOf(val)
// 	m, ok := r.models[typ]
// 	r.lock.RUnlock()
// 	if ok {
// 		return m, nil
// 	}
// 	r.lock.Lock()
// 	defer r.lock.Unlock()
// 	m, ok = r.models[typ]
// 	if ok {
// 		return m, nil
// 	}
// 	var err error
// 	if m, err = r.parseModel(typ); err != nil {
// 		return nil, err
// 	}
// 	r.models[typ] = m
// 	return m, nil
// }

func (r *registry) parseModel(typ reflect.Type) (*model, error) {
	if typ.Kind() != reflect.Ptr ||
		typ.Elem().Kind() != reflect.Struct {
		return nil, errs.ErrPointerOnly
	}
	typ = typ.Elem()

	// 获得字段的数量
	numField := typ.NumField()
	fds := make(map[string]*field, numField)
	for i := 0; i < numField; i++ {
		fdType := typ.Field(i)
		fds[fdType.Name] = &field{
			colName: underscoreName(fdType.Name),
		}
	}
	return &model{
		tableName: underscoreName(typ.Name()),
		fieldMap:  fds,
	}, nil
}

// underscoreName 驼峰转字符串命名
func underscoreName(tableName string) string {
	var buf []byte
	for i, v := range tableName {
		if unicode.IsUpper(v) {
			if i != 0 {
				buf = append(buf, '_')
			}
			buf = append(buf, byte(unicode.ToLower(v)))
		} else {
			buf = append(buf, byte(v))
		}

	}
	return string(buf)
}

var models = map[reflect.Type]*model{}
