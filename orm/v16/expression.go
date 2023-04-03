//go:build v16
package orm

// RawExpr 代表一个原生表达式
// 意味着 ORM 不会对它进行任何处理
type RawExpr struct {
	raw  string
	args []interface{}
}

func (r RawExpr) selectedAlias() string {
	return ""
}

func (r RawExpr) fieldName() string {
	return ""
}

func (r RawExpr) target() TableReference {
	return nil
}

func (r RawExpr) assign() {}

func (r RawExpr) expr() {}

func (r RawExpr) AsPredicate() Predicate {
	return Predicate{
		left: r,
	}
}

// Raw 创建一个 RawExpr
func Raw(expr string, args ...interface{}) RawExpr {
	return RawExpr{
		raw:  expr,
		args: args,
	}
}

type binaryExpr struct {
	left  Expression
	op    op
	right Expression
}

func (binaryExpr) expr() {}

type MathExpr binaryExpr

func (m MathExpr) Add(val interface{}) MathExpr {
	return MathExpr{
		left:  m,
		op:    opAdd,
		right: valueOf(val),
	}
}

func (m MathExpr) Multi(val interface{}) MathExpr {
	return MathExpr{
		left:  m,
		op:    opMulti,
		right: valueOf(val),
	}
}

func (m MathExpr) expr() {}

// SubqueryExpr 注意，这个谓词这种不是在所有的数据库里面都支持的
// 这里采取的是和 Upsert 不同的做法
// Upsert 里面我们是属于用 dialect 来区别不同的实现
// 这里我们采用另外一种方案，就是直接生成，依赖于数据库来报错
// 实际中两种方案你可以自由替换
type SubqueryExpr struct {
	s Subquery
	// 谓词，ALL，ANY 或者 SOME
	pred string
}

func (SubqueryExpr) expr() {}

func Any(sub Subquery) SubqueryExpr {
	return SubqueryExpr{
		s: sub,
		pred: "ANY",
	}
}

func All(sub Subquery) SubqueryExpr {
	return SubqueryExpr{
		s: sub,
		pred: "ALL",
	}
}

func Some(sub Subquery) SubqueryExpr {
	return SubqueryExpr{
		s: sub,
		pred: "SOME",
	}
}