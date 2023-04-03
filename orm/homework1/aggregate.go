package orm

// Aggregate 代表聚合函数，例如 AVG, MAX, MIN 等
type Aggregate struct {
	fn    string
	arg   string
	alias string
}

func (a Aggregate) selectable() {}

func (a Aggregate) expr() {}

func (a Aggregate) As(alias string) Aggregate {
	return Aggregate{
		fn:    a.fn,
		arg:   a.arg,
		alias: alias,
	}
}

// EQ 例如 C("id").Eq(12)
func (a Aggregate) EQ(arg any) Predicate {
	return Predicate{
		left:  a,
		op:    opEQ,
		right: exprOf(arg),
	}
}

func (a Aggregate) LT(arg any) Predicate {
	return Predicate{
		left:  a,
		op:    opLT,
		right: exprOf(arg),
	}
}

func (a Aggregate) GT(arg any) Predicate {
	return Predicate{
		left:  a,
		op:    opGT,
		right: exprOf(arg),
	}
}


func Avg(c string) Aggregate {
	return Aggregate{
		fn:  "AVG",
		arg: c,
	}
}

func Max(c string) Aggregate {
	return Aggregate{
		fn:  "MAX",
		arg: c,
	}
}

func Min(c string) Aggregate {
	return Aggregate{
		fn:  "MIN",
		arg: c,
	}
}

func Count(c string) Aggregate {
	return Aggregate{
		fn:  "COUNT",
		arg: c,
	}
}

func Sum(c string) Aggregate {
	return Aggregate{
		fn:  "SUM",
		arg: c,
	}
}