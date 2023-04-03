package errs

import "fmt"

func NewErrKeyNotFound(key string) error{
	return fmt.Errorf("cache: 找不到 key %s", key)

}
