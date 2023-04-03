
package errs

import (
	"errors"
	"fmt"
)

var (
	// ErrPointerOnly 只支持一级指针作为输入
	// 看到这个 error 说明你输入了其它的东西
	// 我们并不希望用户能够直接使用 err == ErrPointerOnly
	// 所以放在我们的 internal 包里
	ErrPointerOnly = errors.New("orm: 只支持一级指针作为输入，例如 *User")
	ErrNoRows                 = errors.New("orm: 未找到数据")
	ErrTooManyReturnedColumns = errors.New("orm: 过多列")
	// ErrInsertZeroRow 代表插入 0 行
	ErrInsertZeroRow = errors.New("orm: 插入 0 行")
	ErrNoUpdatedColumns = errors.New("orm: 未指定更新的列")
)

// NewErrUnknownField 返回代表未知字段的错误
// 一般意味着你可能输入的是列名，或者输入了错误的字段名
// 注意和 NewErrUnknownColumn 区别
func NewErrUnknownField(fd string) error {
	return fmt.Errorf("orm: 未知字段 %s", fd)
}

// NewErrUnknownColumn 返回代表未知列的错误
// 一般意味着你使用了错误的列名
// 注意和 NewErrUnknownField 区别
func NewErrUnknownColumn(col string) error {
	return fmt.Errorf("orm: 未知列 %s", col)
}

func NewErrUnsupportedAssignableType(exp any) error {
	return fmt.Errorf("orm: 不支持的 Assignable 表达式 %v", exp)
}

// NewErrUnsupportedExpressionType 返回一个不支持该 expression 错误信息
func NewErrUnsupportedExpressionType(exp any) error {
	return fmt.Errorf("orm: 不支持的表达式 %v", exp)
}

// NewErrUnsupportedTableType 返回一个不支持该 TableReference 错误信息
func NewErrUnsupportedTableType(exp any) error {
	return fmt.Errorf("orm: 不支持的 TableReference %v", exp)
}

// NewErrUnsupportedSelectable 返回一个不支持该 selectable 的错误信息
// 即 exp 不能作为 SELECT xxx 的一部分
func NewErrUnsupportedSelectable(exp any) error {
	return fmt.Errorf("orm: 不支持的目标列 %v", exp)
}

// 后面可以考虑支持错误码
// func NewErrUnsupportedExpressionType(exp any) error {
// 	return fmt.Errorf("orm-50001: 不支持的表达式 %v", exp)
// }

// 后面还可以考虑用 AST 分析源码，生成错误排除手册，例如
// @ErrUnsupportedExpressionType 40001
// 发生该错误，主要是因为传入了不支持的 Expression 的实际类型
// 一般来说，这是因为中间件

func NewErrInvalidTagContent(tag string) error {
	return fmt.Errorf("orm: 错误的标签设置: %s", tag)
}

func NewErrFailToRollbackTx(bizErr error, rbErr error, panicked bool) error {
	return fmt.Errorf("orm: 回滚事务失败, 业务错误 %w, 回滚错误 %s, panic: %t",
		bizErr, rbErr.Error(), panicked)
}