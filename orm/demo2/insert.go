package orm

import (
	"gitee.com/geektime-geekbang/geektime-go/orm/demo2/internal/errs"
	"gitee.com/geektime-geekbang/geektime-go/orm/demo2/model"
	"reflect"
)

type Inserter[T any] struct {
	builder
	db *DB
	values []*T
	columns []string
	onDuplicate *OnDuplicateKey
}

func (i *Inserter[T]) Build() (*Query, error) {
	if len(i.values) == 0 {
		return nil, errs.ErrInsertZeroRow
	}
	i.sb.WriteString("INSERT INTO ")
	m, err := i.db.r.Get(i.values[0])
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
		refVal := reflect.ValueOf(val).Elem()
		for j, c := range fields {
			if j > 0 {
				i.sb.WriteByte(',')
			}
			i.sb.WriteByte('?')
			i.args = append(i.args, refVal.FieldByIndex(c.Index).Interface())
		}
		i.sb.WriteByte(')')
	}

	if i.onDuplicate != nil {
		err = i.dialect.buildDuplicateKey(&i.builder, i.onDuplicate)
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

func (i *Inserter[T]) OnDuplicateKey() *OnDuplicateKeyBuilder[T] {
	return &OnDuplicateKeyBuilder[T]{
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

func NewInserter[T any](db *DB) *Inserter[T]{
	return &Inserter[T]{
		builder: builder{
			dialect: db.dialect,
		},
		db: db,
	}
}

type OnDuplicateKeyBuilder[T any] struct {
	i *Inserter[T]
}

func (o *OnDuplicateKeyBuilder[T]) Update(assigns...Assignable) *Inserter[T] {
	o.i.onDuplicate = &OnDuplicateKey{
		assigns: assigns,
	}
	return o.i
}

type OnDuplicateKey struct {
	assigns []Assignable
}


