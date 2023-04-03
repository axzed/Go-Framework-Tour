package orm

import (
	"context"
	"gitee.com/geektime-geekbang/geektime-go/orm/demo3/internal/errs"
	"gitee.com/geektime-geekbang/geektime-go/orm/demo3/model"
)

type Inserter[T any] struct {
	builder
	core
	sess Session
	values []*T
	columns []string
	onDuplicate *Upsert
}

func (i *Inserter[T]) Exec(ctx context.Context) Result {
	q, err := i.Build()
	if err != nil {
		return Result{
			err: err,
		}
	}
	res, err := i.sess.execContext(ctx, q.SQL, q.Args...)
	return Result{
		res: res,
		err: err,
	}
}

// type MySQLInserter[T any] struct {
// 	Inserter[T]
// }
//
// func (m *MySQLInserter[T]) OnDuplicateKey() *MySQLInserter[T] {
//
// }
//
// type PostgreSQL[T any] struct {
// 	Inserter[T]
// }
//
// func (m *PostgreSQL[T]) OnConflict(cols...string) *PostgreSQL[T] {
//
// }


func (i *Inserter[T]) Build() (*Query, error) {
	if len(i.values) == 0 {
		return nil, errs.ErrInsertZeroRow
	}
	i.sb.WriteString("INSERT INTO ")
	m, err := i.r.Get(i.values[0])
	if err != nil {
		return nil, err
	}
	i.model = m
	i.quote(m.TableName)
	fields := m.Fields
	// 你指定了
	if len(i.columns) != 0 {
		fields = make([]*model.Field, 0, len(i.columns))
		for _, c := range i.columns {
			fd, ok := m.FieldMap[c]
			if !ok {
				return nil, errs.NewErrUnknownField(c)
			}
			fields = append(fields, fd)
		}
	}

	i.sb.WriteByte('(')
	for idx, c := range fields {
		if idx > 0 {
			i.sb.WriteByte(',')
		}
		i.quote(c.ColName)
	}

	i.sb.WriteByte(')')
	i.sb.WriteString(" VALUES")
	i.args = make([]any, 0, len(i.values) * len(m.Fields))
	for idx, val := range i.values {
		if idx > 0 {
			i.sb.WriteByte(',')
		}
		i.sb.WriteByte('(')
		refVal := i.valCreator(val, i.model)
		for j, c := range fields {
			if j > 0 {
				i.sb.WriteByte(',')
			}
			i.sb.WriteByte('?')
			fdVal, err := refVal.Field(c.GoName)
			if err != nil {
				return nil, err
			}
			i.args = append(i.args, fdVal)
		}
		i.sb.WriteByte(')')
	}

	if i.onDuplicate != nil {
		err = i.core.dialect.buildDuplicateKey(&i.builder, i.onDuplicate)
		if err != nil {
			return nil, err
		}
	}

	i.sb.WriteByte(';')

	return &Query{
		SQL: i.sb.String(),
		Args: i.args,
	}, nil
}

func (i *Inserter[T]) Update() *UpsertBuilder[T] {
	return &UpsertBuilder[T]{
		i: i,
	}
}

func (i *Inserter[T]) Columns(cs...string) *Inserter[T]{
	i.columns = cs
	return i
}

// Values 指定 INSERT INTO xxx VALUES 的 VALUES 部分
func (i *Inserter[T]) Values(vals...*T) *Inserter[T]{
	i.values = vals
	return i
}

func NewInserter[T any](sess Session) *Inserter[T]{
	c := sess.getCore()
	return &Inserter[T]{
		builder: builder{
			dialect: c.dialect,
		},
		core: c,
		sess: sess,
	}
}

type UpsertBuilder[T any] struct {
	i *Inserter[T]
	// where []Predicate
	conflictColumns []string
}

func (o *UpsertBuilder[T]) ConflictColumns(cols ...string) *UpsertBuilder[T] {
	o.conflictColumns = cols
	return o
}

// func (o *UpsertBuilder[T]) Where(ps...Predicate) *UpsertBuilder[T] {
// 	o.where = ps
// 	return o
// }

func (o *UpsertBuilder[T]) Update(assigns...Assignable) *Inserter[T] {
	o.i.onDuplicate = &Upsert{
		conflictColumns: o.conflictColumns,
		assigns: assigns,
	}
	return o.i
}

// func (o *UpsertBuilder[T]) DoNothing(assigns...Assignable) *Inserter[T] {
// 	o.i.onDuplicate = &Update{
// 		doNothing: true,
// 	}
// 	return o.i
// }

type Upsert struct {
	assigns []Assignable
	conflictColumns []string
	// doNothing bool
}


