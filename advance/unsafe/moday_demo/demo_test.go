package moday_demo

import (
	"errors"
	"fmt"
	"io"
	"testing"
)

func TestMyService_Serve(t *testing.T) {
	ms := &MyService{
		handler: &mockHandler{},
	}
	ms.Serve()
}

type mockHandler struct {
}

func (m *mockHandler) Handle() {
	fmt.Println("hello")
}

type Number interface {
	int32 | int64
}

func Sum[T Number](i T, j T) T {
	return i + j
}

func SumV1(i any, j any) any {
	var sum int64 = 0
	switch iv := i.(type) {
	case int32:
		sum = sum + int64(iv)
	case int64:

	}

	switch iv := j.(type) {
	case int32:
		sum = sum + int64(iv)
	case int64:

	}
	return sum
}

func Write[T io.Writer](i T, j T) {
	i.Write([]byte("aa"))
	j.Write([]byte("aa"))
}

func WriteV1(i io.Writer, j io.Writer) {
	i.Write([]byte("aa"))
	j.Write([]byte("aa"))
}

func WriteV2(i any, j any) {
	i.(io.Writer).Write([]byte("aa"))
	j.(io.Writer).Write([]byte("aa"))
}

type Base struct {
}

func (b *Base) A() {
	fmt.Println("A")
	b.B()
}

func (b *Base) B() {
	fmt.Println("B")
}

type Sub struct {
	Base
}

func (s *Sub) B() {
	fmt.Println("Sub")
}

func TestInherit(t *testing.T) {
	s := &Sub{}
	// 真正的继承，应该输出 A Sub
	// 组合输出什么 A, B
	s.A()
	// 等价于 s.Base.A(), s.Base.B()
}

func UseSub() {
	// 用户只用指针，那么 Base 直接组合
	s := &Sub{}
	fmt.Println(s)
	// 用户用结构体（比如说为了不可变对象），那么 Base 组合
	s1 := Sub{}
	fmt.Println(s1)
	// 用户用结构体，但是 Base 不能被复制（例如 sync.XXX），那么用 *Base 组合
	// type Sub struct {
	//	 *Base
	// }
	//
	//
	s2 := Sub{}
	fmt.Println(s2)
}

// 用组合来实现装饰器模式

type BaseDecorator struct {
	Base
}

func (p BaseDecorator) A() {
	// 做点事情

	p.Base.A()

	// 再做点事情
}

// 假如你有一个结构体
// 某些字段是必须传的
// 某些字段不是必须的，或者全部字段都是非必须的
type SomeStruct struct {
	// 必须的
	field1 string
	field2 int

	// 非必须的
	field3 bool
	field4 []byte
}

type SomeStructOption func(s *SomeStruct)

func WithSomeStructField3(fd3 bool) SomeStructOption {
	return func(s *SomeStruct) {
		s.field3 = fd3
	}
}

func NewSomeStruct(fd1 string, fd2 int, opts ...SomeStructOption) *SomeStruct {
	res := &SomeStruct{
		field1: fd1,
		field2: fd2,

		// 非必选字段默认值设置好
		field3: true,
		field4: []byte("hello"),
	}

	for _, opt := range opts {
		opt(res)
	}
	return res
}

type SomeStructBuilder struct {
	// 必须的
	field1 string
	field2 int

	// 非必须的
	field3 bool
	field4 []byte
}

func NewSomeStructBuilder(fd1 string, fd2 int) *SomeStructBuilder {
	return &SomeStructBuilder{
		field1: fd1,
		field2: fd2,

		// 非必选字段默认值设置好
		field3: true,
		field4: []byte("hello"),
	}
}

// 中间方法，依旧返回 builder
func (s *SomeStructBuilder) BuildField3(fd bool) *SomeStructBuilder {
	s.field3 = fd
	return s
}

// 终结方法，返回最终的成果
func (*SomeStructBuilder) Build() *SomeStruct {
	// 也可以设置默认值
	// 在这里完成构建，
	return &SomeStruct{}
}

type SomeStructErr struct {
	// 必须的
	field1 string
	field2 int

	// 非必须的
	field3 bool
	field4 []byte
}

type SomeStructErrOption func(s *SomeStructErr) error

func WithSomeStructErrField4(fd4 []byte) SomeStructErrOption {
	return func(s *SomeStructErr) error {
		if len(fd4) == 0 {
			return errors.New("fd4 不能为空")
		}
		s.field4 = fd4
		return nil
	}
}

func NewSomeStructErr(fd1 string, fd2 int, opts ...SomeStructErrOption) (*SomeStructErr, error) {
	res := &SomeStructErr{
		field1: fd1,
		field2: fd2,

		// 非必选字段默认值设置好
		field3: true,
		field4: []byte("hello"),
	}

	for _, opt := range opts {
		err := opt(res)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}
