//go:build v16
package orm

type Column struct {
	table TableReference
	name  string
	alias string
}

func (c Column) assign() {}

func (c Column) expr() {}

func (c Column) selectedAlias() string {
	return c.alias
}

func (c Column) fieldName() string {
	return c.name
}

func (c Column) target() TableReference {
	return c.table
}

func (c Column) As(alias string) Column {
	return Column {
		name:  c.name,
		alias: alias,
	}
}

type value struct {
	val any
}

func (c value) expr() {}

func valueOf(val any) value {
	return value{
		val: val,
	}
}

func C(name string) Column {
	return Column{name: name}
}

func (c Column) Add(delta int) MathExpr {
	return MathExpr{
		left: c,
		op: opAdd,
		right: value{val: delta},
	}
}

func (c Column) Multi(delta int) MathExpr {
	return MathExpr{
		left: c,
		op: opAdd,
		right: value{val: delta},
	}
}

// EQ 例如 C("id").Eq(12)
func (c Column) EQ(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opEQ,
		right: exprOf(arg),
	}
}

func (c Column) LT(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opLT,
		right: exprOf(arg),
	}
}

func (c Column) GT(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opGT,
		right: exprOf(arg),
	}
}

// In 有两种输入，一种是 IN 子查询
// 另外一种就是普通的值
// 这里我们可以定义两个方法，如 In  和 InQuery，也可以定义一个方法
// 这里我们使用一个方法
func (c Column) In(vals...any) Predicate {
	return Predicate{
		left:  c,
		op:    opIN,
		right: valueOf(vals),
	}
}

func (c Column) InQuery(sub Subquery) Predicate {
	return Predicate{
		left:  c,
		op:    opIN,
		right: sub,
	}
}
