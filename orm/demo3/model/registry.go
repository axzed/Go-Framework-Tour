package model

import (
	"gitee.com/geektime-geekbang/geektime-go/orm/demo3/internal/errs"
	"reflect"
	"strings"
	"sync"
	"unicode"
)

type Registry interface {
	Get(val any) (*Model, error)
	Register(val any, opts...Option) (*Model, error)
}

type registry struct {
	models sync.Map
}

func NewRegistry() Registry {
	return &registry{}
}

func (r *registry) Get(val any) (*Model, error) {
	typ := reflect.TypeOf(val)
	m, ok := r.models.Load(typ)
	if ok {
		return m.(*Model), nil
	}
	return r.Register(val)
}

// Register 输入不能为 nil
func (r *registry) Register(val any, opts...Option) (*Model, error) {
	if val == nil {
		return nil, errs.ErrInputNil
	}
	typ := reflect.TypeOf(val)
	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
		return nil,  errs.ErrPointerOnly
	}
	typ = typ.Elem()
	numField := typ.NumField()
	fieldMap := make(map[string]*Field, numField)
	colMap := make(map[string]*Field, numField)
	columns := make([]*Field, numField)
	for i := 0; i < numField; i++ {
		fd := typ.Field(i)
		ormTagKVs := r.parseTag(fd.Tag)
		colName, ok := ormTagKVs["column"]
		if !ok || colName == "" {
			colName = underscoreName(fd.Name)
		}
		fdMeta := &Field{
			ColName: colName,
			Type:    fd.Type,
			GoName:  fd.Name,
			Offset:  fd.Offset,
			Index: fd.Index,
		}
		fieldMap[fd.Name] = fdMeta
		colMap[colName] = fdMeta
		columns[i]=fdMeta
	}

	var tableName string
	if tn, ok := val.(TableName); ok {
		tableName = tn.TableName()
	}
	if tableName == "" {
		tableName = underscoreName(typ.Name())
	}

	res := &Model{
		TableName: tableName,
		FieldMap:  fieldMap,
		ColumnMap: colMap,
		Fields:    columns,
	}
	for _, opt := range opts {
		if err := opt(res); err != nil {
			return nil, err
		}

	}
	r.models.Store(typ, res)
	return res, nil
}

// column => id
func (r *registry)parseTag(tag reflect.StructTag) map[string]string {
	ormTag := tag.Get("orm")
	kvs := strings.Split(ormTag, ",")
	res := make(map[string]string, len(kvs))
	for _, kv := range kvs {
		segs := strings.Split(kv, "=")
		key := segs[0]
		var v = ""
		if len(segs) > 1 {
			v = segs[1]
		}
		res[key]= v
	}
	return res
}

// underscoreName 驼峰转字符串命名
func underscoreName(tableName string) string {
	var buf []byte
	for i, v := range tableName {
		if unicode.IsUpper(v) {
			if i != 0 {
				buf = append(buf, '_')
			}
			buf = append(buf, byte(unicode.ToLower(v)))
		} else {
			buf = append(buf, byte(v))
		}

	}
	return string(buf)
}