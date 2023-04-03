package orm

import (
	"testing"
)

func TestDBUseReflectValuer(t *testing.T) {
	Open("sqlite3", "file:test.db?cache=shared&mode=memory", DBUseReflectValuer())
}
