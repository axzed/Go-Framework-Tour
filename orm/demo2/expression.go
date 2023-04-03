package orm

type RawExpr struct {
	raw string
	args []any
}

func (RawExpr) selectable() {}

func Raw(expr string, args...any) RawExpr {
	return RawExpr{
		raw: expr,
		args: args,
	}
}

// s.Where(Raw("a =? and b = ? and c = ?", 1, 2, 3))
func (RawExpr) AsPredicate() Predicate {
	return Predicate{}
}