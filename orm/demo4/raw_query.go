package orm

import "context"

type RawQuerier[T any] struct {
	core
	sql string
	args []any
	sess Session
	typ string
}

func (r *RawQuerier[T]) Build() (*Query, error) {
	return &Query{
		SQL: r.sql,
		Args: r.args,
	}, nil
}

// RawQuery[TestModel)("SELECT * FROM XXX WHERE XXX").GET(ctx)
func RawQuery[T any](sess Session, sql string, args...any) *RawQuerier[T] {
	return &RawQuerier[T]{
		sql: sql,
		args: args,
		typ: "RAW",
		sess: sess,
		core: sess.getCore(),
	}
}

// r.Type("SELECT")
// func (r *RawQuerier[T]) Type(typ string) *RawQuerier[T] {
//
// }

func (r *RawQuerier[T]) Get(ctx context.Context) (*T, error) {
	model, err := r.r.Get(new(T))
	if err != nil {
		return nil, err
	}
	res := get[T](ctx, r.core, r.sess, &QueryContext{
		Type: r.typ,
		Builder: r,
		Model: model,
		TableName: model.TableName,
	})
	if res.Result != nil {
		return res.Result.(*T), res.Err
	}
	return nil, res.Err
}

func (r *RawQuerier[T]) get(ctx context.Context) *QueryResult {
	panic("implement me")
}

func get[T any](ctx context.Context, c core, sess Session, qc *QueryContext) *QueryResult {
	var root Handler = func(ctx context.Context, qc *QueryContext) *QueryResult {
		q, err := qc.Builder.Build()
		if err != nil {
			return &QueryResult{
				Err: err,
			}
		}
		rows, err := sess.queryContext(ctx, q.SQL, q.Args...)
		if err != nil {
			return &QueryResult{
				Err: err,
			}
		}

		t := new(T)
		m, err := c.r.Get(t)
		if err != nil {
			return &QueryResult{
				Err: err,
			}
		}
		val := c.valCreator(t, m)
		// 在这里灵活切换反射或者 unsafe
		err = val.SetColumns(rows)
		return &QueryResult{
			Result: t,
			Err: err,
		}
	}
	for i := len(c.ms) - 1; i >= 0 ; i-- {
		root = c.ms[i](root)
	}
	return root(ctx, qc)
}