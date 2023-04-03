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
	sess  Session
}

func NewUpdater[T any](sess Session) *Updater[T]{
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
	if u.model == nil {
		model, err := u.r.Get(u.val)
		if err != nil {
			return nil, err
		}
		u.model = model
	}

	u.sb.WriteString("UPDATE ")
	u.quote(u.model.TableName)
	u.sb.WriteString(" SET ")
	val := u.valCreator(u.val, u.model)
	for i, a := range u.assigns {
		if i > 0 {
			u.sb.WriteByte(',')
		}
		switch assign := a.(type) {
		case Column:
			if err := u.buildColumn(assign.table, assign.name); err != nil {
				return nil, err
			}
			u.sb.WriteString("=?")
			arg, err := val.Field(assign.name)
			if err != nil {
				return nil, err
			}
			u.addArgs(arg)
		case Assignment:
			if err := u.buildAssignment(assign); err != nil {
				return nil, err
			}
		default:
			return nil, errs.NewErrUnsupportedAssignableType(a)
		}
	}
	if len(u.where) > 0 {
		u.sb.WriteString(" WHERE ")
		if err := u.buildPredicates(u.where);err != nil {
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
	if u.model == nil {
		m, err := u.r.Get(new(T))
		if err != nil {
			return Result{
				err: err,
			}
		}
		u.model=m
	}
	return exec(ctx, u.sess, u.core, &QueryContext{
		builder: u,
		Type:    "UPDATE",
		Model:   u.model,
	})
}
