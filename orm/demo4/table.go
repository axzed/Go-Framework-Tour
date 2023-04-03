package orm

type TableReference interface {
	tableAlias() string
}

// 普通表
type Table struct {
	entity any
	alias string
}

// A JOIN B
func TableOf(entity any) Table {
	return Table{
		entity: entity,
	}
}

func (t Table) As(alias string) Table {
	return Table{
		entity: t.entity,
		alias: alias,
	}
}

// TableOf(Order{}).C("Id") => 表达的是定义在 Order 上的 Id
// TableOf(OrderDetail{}).C("Id") => 表达的是定义在 OrderDetail 上的 Id
// TableOf(Item{}).C("Id") => 表达的是定义在 Item 上的 Id
// 解决不同表有同名列的问题
// 在 JOIN 查询里面，C("Id") 你是不知道是哪一个表的，
func (t Table) C(col string) Column {
	return Column{
		name: col,
		tbl: t,
	}
}

// func (t *Table) As(alias string) {
//
// }

// A JOIN B
// A JOIN (B JOIN C)
// (A JOIN B) JOIN (C JOIN D)
func (t Table) Join(right TableReference) *JoinBuilder {
	return &JoinBuilder{
		left: t,
		typ: "JOIN",
		right: right,
	}
}



func (t Table) LeftJoin(right TableReference) *JoinBuilder {
	return &JoinBuilder{
		left: t,
		typ: "LEFT JOIN",
		right: right,
	}
}

func (t Table) tableAlias() string {
	return t.alias
}

// func (t Table) C(name string) Column {
//
// }


// JOIN 查询
type Join struct {
	left TableReference
	// JOIN, LEFT JOIN, LEFT OUTER JOIN
	typ string
	right TableReference
	on []Predicate
	using []string
}


func (t Join) tableAlias() string {
	return ""
}

type JoinBuilder struct {
	left TableReference
	// JOIN, LEFT JOIN, LEFT OUTER JOIN
	typ string
	right TableReference

}

func (jb *JoinBuilder) On(ps...Predicate) Join {
	return Join{
		left: jb.left,
		typ: jb.typ,
		right: jb.right,
		on: ps,
	}
}

func (jb *JoinBuilder) Using(cols...string) Join {
	return Join{
		left: jb.left,
		typ: jb.typ,
		right: jb.right,
		using: cols,
	}
}

// 子查询
type Subquery struct {

}

// func (t Subquery) C(name string) Column {
//
// }
