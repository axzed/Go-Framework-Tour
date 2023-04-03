package demo

import (
	"context"
	"errors"
	"reflect"
	"strings"
)

type Selector[T any] struct {
	sb strings.Builder
	tbl string
	where []Predicate

	args []any
}

func (s *Selector[T]) Where(ps...Predicate) *Selector[T] {
	s.where = ps
	return s
}

// func (s *Selector[T]) Where(where string, args...any) *Selector[T] {
//
// }

// 为了链式调用
func (s *Selector[T]) From(tbl string) *Selector[T] {
	s.tbl = tbl
	return s
}

func (s *Selector[T]) Build() (*Query, error) {
	s.sb.WriteString("SELECT * FROM ")

	if s.tbl == "" {
		var t T
		typ := reflect.TypeOf(t)
		goName := typ.Name()
		s.sb.WriteByte('`')
		s.sb.WriteString(goName)
		s.sb.WriteByte('`')
	} else {
		s.sb.WriteString(s.tbl)
	}

	if len(s.where) >0 {
		// 构造 WHERE
		s.sb.WriteString(" WHERE ")
		pred := s.where[0]
		for i := 1; i < len(s.where); i++ {
			pred=pred.And(s.where[i])
		}
		err := s.buildExpression(pred)
		if err != nil {
			return nil, err
		}
	}

	s.sb.WriteByte(';')
	return &Query{
		SQL: s.sb.String(),
		Args: s.args,
	}, nil
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Selector[T]) buildExpression(expression Expression) error {
	switch expr := expression.(type) {
	case nil:
		return nil
	case Column:
		s.sb.WriteByte('`')
		s.sb.WriteString(expr.name)
		s.sb.WriteByte('`')
	case Value:
		s.sb.WriteByte('?')
		if s.args == nil {
			s.args = make([]any, 0, 8)
		}
		s.args = append(s.args, expr.val)
	case Predicate:
		_, ok := expr.left.(Predicate)
		if ok {
			s.sb.WriteByte('(')
		}
		if err := s.buildExpression(expr.left); err != nil {
			return err
		}
		if ok {
			s.sb.WriteByte(')')
		}

		s.sb.WriteByte(' ')
		s.sb.WriteString(expr.op.String())
		s.sb.WriteByte(' ')
		_, ok = expr.right.(Predicate)
		if ok {
			s.sb.WriteByte('(')
		}
		if err:= s.buildExpression(expr.right); err != nil {
			return err
		}
		if ok {
			s.sb.WriteByte(')')
		}
	default:
		return errors.New("不支持的表达式")
	}
	return nil
}


func NewSelector[T any]() *Selector[T] {
	return &Selector[T]{}
}


