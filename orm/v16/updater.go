//go:build v16
package orm

import (
	"context"
	"gitee.com/geektime-geekbang/geektime-go/orm/internal/errs"
)

type Updater[T any] struct {
	builder
	assigns []Assignable
	val *T
	where []Predicate
	sess session
}

func NewUpdater[T any](sess session) *Updater[T]{
	c := sess.getCore()
	return &Updater[T]{
		builder: builder{
			core: c,
			dialect: c.dialect,
			quoter: c.dialect.quoter(),
		},
		sess: sess,
	}
}

func (u *Updater[T]) Update(t *T) *Updater[T] {
	u.val = t
	return u
}

func (u *Updater[T]) Set(assigns...Assignable) *Updater[T] {
	u.assigns = assigns
	return u
}

func (u *Updater[T]) Build() (*Query, error) {
	if len(u.assigns) == 0 {
		return nil, errs.ErrNoUpdatedColumns
	}
	if u.val == nil {
		u.val = new(T)
	}
	model, err := u.r.Get(u.val)
	if err != nil {
		return nil, err
	}
	u.model = model
	u.sb.WriteString("UPDATE ")
	u.quote(model.TableName)
	u.sb.WriteString(" SET ")
	val := u.valCreator(u.val, model)
	for i, a := range u.assigns {
		if i > 0 {
			u.sb.WriteByte(',')
		}
		switch assign := a.(type) {
		case Column:
			if err = u.buildColumn(assign.table, assign.name); err != nil {
				return nil, err
			}
			u.sb.WriteString("=?")
			arg, err := val.Field(assign.name)
			if err != nil {
				return nil, err
			}
			u.addArgs(arg)
		case Assignment:
			if err = u.buildAssignment(assign); err != nil {
				return nil, err
			}
		default:
			return nil, errs.NewErrUnsupportedAssignableType(a)
		}
	}
	if len(u.where) > 0 {
		u.sb.WriteString(" WHERE ")
		if err = u.buildPredicates(u.where);err != nil {
			return nil, err
		}
	}
	u.sb.WriteByte(';')
	return &Query{
		SQL:  u.sb.String(),
		Args: u.args,
	}, nil
}

func (u *Updater[T]) buildAssignment(assign Assignment) error {
	if err := u.buildColumn(nil, assign.column); err != nil {
		return err
	}
	u.sb.WriteByte('=')
	return u.buildExpression(assign.val)
}

func (u *Updater[T]) Where(ps ...Predicate) *Updater[T] {
	u.where = ps
	return u
}

func (u *Updater[T]) Exec(ctx context.Context) Result {
	q, err := u.Build()
	if err != nil {
		return Result{err: err}
	}
	res, err := u.sess.execContext(ctx, q.SQL, q.Args...)
	return Result{err: err, res: res}
}
