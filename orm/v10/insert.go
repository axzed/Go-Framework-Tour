//go:build v10

package orm

import (
	"gitee.com/geektime-geekbang/geektime-go/orm/internal/errs"
	"gitee.com/geektime-geekbang/geektime-go/orm/model"
	"reflect"
	"strings"
)

type OnDuplicateKeyBuilder[T any] struct {
	i *Inserter[T]
}

type OnDuplicateKey struct {
	assigns []Assignable
}

func (o *OnDuplicateKeyBuilder[T]) Update(assigns ...Assignable) *Inserter[T] {
	o.i.onDuplicate = &OnDuplicateKey{
		assigns: assigns,
	}
	return o.i
}

type Inserter[T any] struct {
	values  []*T
	db      *DB
	columns []string
	sb      strings.Builder
	args    []any
	model   *model.Model
	// 方案二
	onDuplicate *OnDuplicateKey

	// 方案一
	// onDuplicate []Assignable
}

func NewInserter[T any](db *DB) *Inserter[T] {
	return &Inserter[T]{
		db: db,
	}
}

func (i *Inserter[T]) Values(vals ...*T) *Inserter[T] {
	i.values = vals
	return i
}

func (i *Inserter[T]) OnDuplicateKey() *OnDuplicateKeyBuilder[T] {
	return &OnDuplicateKeyBuilder[T]{
		i: i,
	}
}

// 方案一
// func (i *Inserter[T]) OnDuplicateKeyBuilder(assigns...Assignable) *Inserter[T] {
// 	i.onDuplicate = assigns
// 	return i
// }

// Fields 指定要插入的列
// TODO 目前我们只支持指定具体的列，但是不支持复杂的表达式
// 例如不支持 VALUES(..., now(), now()) 这种在 MySQL 里面常用的
func (i *Inserter[T]) Columns(cols ...string) *Inserter[T] {
	i.columns = cols
	return i
}

func (i *Inserter[T]) Build() (*Query, error) {
	if len(i.values) == 0 {
		return nil, errs.ErrInsertZeroRow
	}
	m, err := i.db.r.Get(i.values[0])
	if err != nil {
		return nil, err
	}
	i.model = m

	i.sb.WriteString("INSERT INTO `")
	i.sb.WriteString(m.TableName)
	i.sb.WriteString("`(")

	fields := m.Fields
	if len(i.columns) != 0 {
		fields = make([]*model.Field, 0, len(i.columns))
		for _, c := range i.columns {
			field, ok := m.FieldMap[c]
			if !ok {
				return nil, errs.NewErrUnknownField(c)
			}
			fields = append(fields, field)
		}
	}

	// (len(i.values) + 1) 中 +1 是考虑到 UPSERT 语句会传递额外的参数
	i.args = make([]any, 0, len(fields)*(len(i.values)+1))
	for idx, fd := range fields {
		if idx > 0 {
			i.sb.WriteByte(',')
		}
		i.sb.WriteByte('`')
		i.sb.WriteString(fd.ColName)
		i.sb.WriteByte('`')
	}

	i.sb.WriteString(") VALUES")
	for vIdx, val := range i.values {
		if vIdx > 0 {
			i.sb.WriteByte(',')
		}
		refVal := reflect.ValueOf(val).Elem()
		i.sb.WriteByte('(')
		for fIdx, field := range fields {
			if fIdx > 0 {
				i.sb.WriteByte(',')
			}
			i.sb.WriteByte('?')
			fdVal := refVal.Field(field.Index)
			i.addArgs(fdVal.Interface())
		}
		i.sb.WriteByte(')')
	}

	if i.onDuplicate != nil {
		i.sb.WriteString(" ON DUPLICATE KEY UPDATE ")
		for idx, assign := range i.onDuplicate.assigns {
			if idx > 0 {
				i.sb.WriteByte(',')
			}
			if err = i.buildAssignment(assign); err != nil {
				return nil, err
			}
		}
	}

	i.sb.WriteString(";")
	return &Query{
		SQL:  i.sb.String(),
		Args: i.args,
	}, nil
}

func (i *Inserter[T]) buildAssignment(a Assignable) error {
	switch assign := a.(type) {
	case Column:
		i.sb.WriteByte('`')
		fd, ok := i.model.FieldMap[assign.name]
		if !ok {
			return errs.NewErrUnknownField(assign.name)
		}
		i.sb.WriteString(fd.ColName)
		i.sb.WriteString("`=VALUES(`")
		i.sb.WriteString(fd.ColName)
		i.sb.WriteString("`)")
	case Assignment:
		i.sb.WriteByte('`')
		fd, ok := i.model.FieldMap[assign.column]
		if !ok {
			return errs.NewErrUnknownField(assign.column)
		}
		i.sb.WriteString(fd.ColName)
		i.sb.WriteByte('`')
		i.sb.WriteString("=?")
		i.addArgs(assign.val)
	default:
		return errs.NewErrUnsupportedAssignableType(a)
	}
	return nil
}

func (i *Inserter[T]) addArgs(args ...any) {
	i.args = append(i.args, args...)
}
