package v1

import (
	"context"
	"encoding/json"
	"reflect"
)

func InitClientProxy(address string, val Service) error {
	c, err := NewClient(address)
	if err != nil {
		return err
	}
	setFuncField(val, c)
	return nil
}

// 这个单独的拆出来，就是为了测试，我们可以考虑传入一个 mock proxy
func setFuncField(val Service, c Proxy) {
	v := reflect.ValueOf(val)
	ele := v.Elem()
	t := ele.Type()

	numField := t.NumField()
	for i := 0; i < numField; i++ {
		field := t.Field(i)
		fieldValue := ele.Field(i)
		if fieldValue.CanSet() {
			fn := func(args []reflect.Value) (results []reflect.Value) {
				in := args[1].Interface()
				out := reflect.New(field.Type.Out(0).Elem()).Interface()
				inData, err := json.Marshal(in)
				if err != nil {
					return []reflect.Value{reflect.ValueOf(out), reflect.ValueOf(err)}
				}

				req := &Request{
					ServiceName: val.ServiceName(),
					Method:      field.Name,
					Data:        inData,
				}

				// 要在下面考虑发过去

				resp, err := c.Invoke(args[0].Interface().(context.Context), req)
				if err != nil {
					return []reflect.Value{reflect.ValueOf(out), reflect.ValueOf(err)}
				}

				err = json.Unmarshal(resp.Data, out)

				if err != nil {
					return []reflect.Value{reflect.ValueOf(out), reflect.ValueOf(err)}
				}

				return []reflect.Value{reflect.ValueOf(out), reflect.Zero(reflect.TypeOf(new(error)).Elem())}
			}
			fieldValue.Set(reflect.MakeFunc(field.Type, fn))
		}
	}
}
