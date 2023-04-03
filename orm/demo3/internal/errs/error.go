package errs

import (
	"errors"
	"fmt"
)

var (
	ErrInputNil = errors.New("orm: 不支持 nil")
	ErrPointerOnly = errors.New("orm: 只支持一级指针作为输入，例如 *User")
	ErrTooManyReturnedColumns = errors.New("orm: 过多列")

	// errUnsupportedExpressionType = errors.New("orm: 不支持的表达式")
	ErrEmptyTableName = errors.New("orm: 表名为空")
	ErrNoRows = errors.New("orm: 未找到数据")

	// ErrInsertZeroRow 代表插入 0 行
	ErrInsertZeroRow = errors.New("orm: 插入 0 行")
)

type MyErr struct {
	code string
	msg string
}

func (m MyErr) Error() string {
	return "orm: " + m.code + m.msg
}

func NewErrUnknownColumn(name string) error {
	return fmt.Errorf("orm: 未知列 %s", name)
}

func NewErrUnknownField(name string) error {
	return fmt.Errorf("orm: 未知字段 %s", name)
}

// NewErrUnsupportedExpressionType 返回一个不支持该 expression 错误信息
func NewErrUnsupportedExpressionType(exp any) error {
	return fmt.Errorf("orm: 不支持的表达式 %v", exp)
}

func NewErrUnsupportedExpressionTypeV2(exp any) error {
	return MyErr{
		code: "50001",
		msg: "不支持的表达式",
	}
}

// func NewErrUnsupportedExpressionTypeV1(exp any) error {
// 	return fmt.Errorf("%w %v", errUnsupportedExpressionType, exp)
// }