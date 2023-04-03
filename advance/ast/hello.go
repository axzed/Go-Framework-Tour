package ast

import (
	"fmt"
)

func Hello(name string) string {
	return "hello, " + name
}

func Hello1(firstName, lastName string, age int) string {
	return fmt.Sprintf("hello, %s %s, %d", firstName, lastName, age)
}
