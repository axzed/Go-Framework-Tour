package demo


type op string

const (
	opEQ = "="
	opLT = "<"
	opGT = ">"

	opNOT = "NOT"
	opAND = "AND"
	opOR = "OR"
)

func (o op) String() string {
	return string(o)
}

type Predicate struct {
	left Expression
	op op
	right Expression
}

func (Predicate) expr() {}

// Eq("id", 12)
// func Eq(column string, arg any) Predicate {
// 	return Predicate{
// 		Column: column,
// 		Op: "=",
// 		Arg: arg,
// 	}
// }

type Column struct {
	name string
}

func (c Column) expr() {}

func C(name string) Column {
	return Column{name: name}
}

// C("id").EQ(12)
func (c Column) Eq(val any) Predicate {
	return Predicate{
		left: c,
		op: opEQ,
		right: Value{val: val},
	}
}

func (c Column) GT(val any) Predicate {
	return Predicate{
		left: c,
		op: opGT,
		right: Value{val: val},
	}
}

func (c Column) LT(val any) Predicate {
	return Predicate{
		left: c,
		op: opLT,
		right: Value{val: val},
	}
}

// Not(C("id").Eq(12))
// NOT (id = ?), 12
func Not(p Predicate) Predicate {
	return Predicate{
		op: opNOT,
		right: p,
	}
}

// And( C("id").Eq(12), C("name").Eq("Tom"))
// func And(p1 Predicate, p2 Predicate) Predicate {
// 	return Predicate{
// 		left: p1,
// 		op: opAND,
// 		right: p2,
// 	}
// }

func (p1 Predicate) Or(p2 Predicate) Predicate {
	return Predicate{
		left: p1,
		op: opOR,
		right: p2,
	}
}

// C("id").Eq(12).And(C("name").Eq("Tom"))
func (p1 Predicate) And(p2 Predicate) Predicate {
	return Predicate{
		left: p1,
		op: opAND,
		right: p2,
	}
}

type Value struct {
	val any
}

func (Value) expr() {}