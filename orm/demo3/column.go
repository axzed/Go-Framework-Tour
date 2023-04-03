package orm

// SELECT xxx
// WHERE xxx
// GROUP BY
// ORDER BY
type Column struct {
	name string
}

func (c Column) selectable() {}
// 这种就是 GORM 的路
// func (c Column) selectable(s *builder) error {
// 	fd, ok := s.model.FieldMap[c.name]
// 	if !ok {
// 		return errs.NewErrUnknownField(c.name)
// 	}
// 	s.sb.WriteByte('`')
// 	s.sb.WriteString(fd.ColName)
// 	s.sb.WriteByte('`')
// }

func (c Column) assign() {}

func (c Column) expr() {}

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
