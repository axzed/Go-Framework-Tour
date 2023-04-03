package valuer

import (
	"database/sql"
	"gitee.com/geektime-geekbang/geektime-go/orm/homework1/internal/errs"
	"gitee.com/geektime-geekbang/geektime-go/orm/homework1/model"
	"reflect"
	"unsafe"
)

type unsafeValue struct {
	addr unsafe.Pointer
	meta *model.Model
}

var _ Creator = NewUnsafeValue

func NewUnsafeValue(val interface{}, meta *model.Model) Value {
	return unsafeValue{
		addr: unsafe.Pointer(reflect.ValueOf(val).Pointer()),
		meta: meta,
	}
}


func (u unsafeValue) SetColumns(rows *sql.Rows) error {
	cs, err := rows.Columns()
	if err != nil {
		return err
	}
	if len(cs) > len(u.meta.ColumnMap) {
		return errs.ErrTooManyReturnedColumns
	}

	colValues := make([]interface{}, len(cs))
	for i, c := range cs {
		cm, ok := u.meta.ColumnMap[c]
		if !ok {
			return errs.NewErrUnknownColumn(c)
		}
		ptr := unsafe.Pointer(uintptr(u.addr) + cm.Offset)
		val := reflect.NewAt(cm.Type, ptr)
		colValues[i]=val.Interface()
	}
	return rows.Scan(colValues...)
}
