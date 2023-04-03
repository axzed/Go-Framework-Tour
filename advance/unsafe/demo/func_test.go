package demo

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestIterateFuncs(t *testing.T) {
	type args struct {
		val any
	}
	tests := []struct {
		name string

		args args

		want    map[string]*FuncInfo
		wantErr error
	}{
		{
			name:    "nil",
			wantErr: errors.New("输入 nil"),
		},
		{
			name: "basic types",
			args: args{
				val: 123,
			},
			wantErr: errors.New("不支持类型"),
		},
		{
			name: "struct type",
			args: args{
				val: Order{
					buyer:  18,
					seller: 100,
				},
			},
			want: map[string]*FuncInfo{
				"GetBuyer": {
					Name:   "GetBuyer",
					In:     []reflect.Type{reflect.TypeOf(Order{})},
					Out:    []reflect.Type{reflect.TypeOf(int64(0))},
					Result: []any{int64(18)},
				},
			},
		},
		{
			name: "struct type but input ptr",
			args: args{
				val: &Order{
					buyer:  18,
					seller: 100,
				},
			},
			want: map[string]*FuncInfo{
				"GetBuyer": {
					Name:   "GetBuyer",
					In:     []reflect.Type{reflect.TypeOf(&Order{})},
					Out:    []reflect.Type{reflect.TypeOf(int64(0))},
					Result: []any{int64(18)},
				},
			},
		},
		{
			name: "pointer type",
			args: args{
				val: &OrderV1{
					buyer:  18,
					seller: 100,
				},
			},
			want: map[string]*FuncInfo{
				"GetBuyer": {
					Name:   "GetBuyer",
					In:     []reflect.Type{reflect.TypeOf(&OrderV1{})},
					Out:    []reflect.Type{reflect.TypeOf(int64(0))},
					Result: []any{int64(18)},
				},
			},
		},

		{
			name: "pointer type but input struct",
			args: args{
				val: OrderV1{
					buyer:  18,
					seller: 100,
				},
			},
			want: map[string]*FuncInfo{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IterateFuncs(tt.args.val)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

type Order struct {
	buyer  int64
	seller int64
}

// 反射层面
// func GetBuyer(o Order) int64 {

func (o Order) GetBuyer() int64 {
	return o.buyer
}

// func (o Order) getSeller() int64 {
// 	return o.seller
// }

type OrderV1 struct {
	buyer  int64
	seller int64
}

// 反射层面
// func GetBuyer(o Order) int64 {

func (o *OrderV1) GetBuyer() int64 {
	// 注入一些代码
	res := o.buyer
	// 注入一些代码
	return res
}

type MyInterface interface {
	Abc()
}

// var _ MyInterface = abcImpl{}
var _ MyInterface = &abcImpl{}

type abcImpl struct {
}

func (a *abcImpl) Abc() {
	// TODO implement me
	panic("implement me")
}

type MyService struct {
	GetById func()
}

// 数据传输
type UserDTO struct {
	Name string
}

// 数据库直接对应
type UserEntity struct {
	Name string
}

// func Copy(src any, dst any) error {
// 	// 反射操作，一个个字段复制过去
// }
//
// // ignoreFields 忽略一些字段，不复制
// func CopyV1(src any, dst any, ignoreFields ...string) error {
// 	// 反射操作，一个个字段复制过去
// }
//
// type Copier struct {
// 	src          any
// 	dst          any
// 	ignoreFields []string
// }
//
// func (c Copier) Copy() error {
//
// }
//
// type User struct {
// 	Name  string
// 	Email string
// }
