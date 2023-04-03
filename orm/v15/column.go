//go:build v15
package orm

type Column struct {
	name  string
	alias string
}

func (c Column) assign() {}

func (c Column) expr() {}

func (c Column) selectable() {}

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
