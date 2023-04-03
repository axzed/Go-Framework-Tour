package valuer

import (
	"database/sql"
	"gitee.com/geektime-geekbang/geektime-go/orm/demo3/internal/errs"
	orm "gitee.com/geektime-geekbang/geektime-go/orm/demo3/model"
	"reflect"
)

type reflectValue struct {
	val reflect.Value
	model *orm.Model
}

func NewReflectValue(t any, model *orm.Model) Value {
	return reflectValue{
		val: reflect.ValueOf(t).Elem(),
		model: model,
	}
}

func (u reflectValue) Field(name string) (any, error) {
	val := u.val
	typ := val.Type()
	_, ok := typ.FieldByName(name)
	if !ok {
		return nil, errs.NewErrUnknownField(name)
	}

	// val = val.FieldByName(name)
	// if val == (reflect.Value{}) {
	// 	return nil, errs.NewErrUnknownField(name)
	// }

	return val.FieldByName(name), nil
}

func (u reflectValue) SetColumns(rows *sql.Rows) error {
	if !rows.Next() {
		return errs.ErrNoRows
	}
	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	if len(cols) > len(u.model.FieldMap) {
		return errs.ErrTooManyReturnedColumns
	}

	colValues := make([]any, 0, len(cols))
	colElemVals := make([]reflect.Value, 0, len(cols))
	for _, col := range cols {
		fd, ok := u.model.ColumnMap[col]
		if !ok {
			return errs.NewErrUnknownColumn(col)
		}
		// fd.Type 是 int，那么  reflect.New(fd.typ) 是 *int
		fdVal := reflect.New(fd.Type)
		colElemVals = append(colElemVals, fdVal.Elem())

		// 因为 Scan 要指针，所以我们在这里，不需要调用 Elem
		colValues = append(colValues, fdVal.Interface())
	}
	// 要把 cols 映射过去字段

	err = rows.Scan(colValues...)
	if err != nil {
		return err
	}
	// 咋办呢？我已经有 colValues 了，接下来咋办？ colValues= [123, "Ming", 18, "Deng"]

	// 反射放回去 t 里面

	tVal := u.val
	for i, col := range cols {
		fd := u.model.ColumnMap[col]
		tVal.FieldByName(fd.GoName).Set(colElemVals[i])
	}
	return nil
}
