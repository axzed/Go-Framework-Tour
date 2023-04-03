package orm

// op 代表操作符
type op string

// 后面可以每次支持新的操作符就加一个
const (
	opEQ  = "="
	opLT  = "<"
	opGT  = ">"
	opAND = "AND"
	opOR  = "OR"
	opNOT = "NOT"
)

func (o op) String() string {
	return string(o)
}

// Expression 代表语句，或者语句的部分
// 暂时没想好怎么设计方法，所以直接做成标记接口
type Expression interface {
	expr()
}

func exprOf(e any) Expression {
	switch exp := e.(type) {
	case Expression:
		return exp
	default:
		return valueOf(exp)
	}
}

// Predicate 代表一个查询条件
// Predicate 可以通过和 Predicate 组合构成复杂的查询条件
type Predicate struct {
	left  Expression
	op    op
	right Expression
}

func (Predicate) expr() {}


func Not(p Predicate) Predicate {
	return Predicate{
		op:    opNOT,
		right: p,
	}
}

func (p Predicate) And(r Predicate) Predicate {
	return Predicate{
		left:  p,
		op:    opAND,
		right: r,
	}
}

func (p Predicate) Or(right Predicate) Predicate {
	return Predicate{
		left:  p,
		op:    opOR,
		right: right,
	}
}