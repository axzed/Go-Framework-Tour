package annotation

import (
	"go/ast"
)

// SingleFileEntryVisitor 这部分和课堂演示差不多，但是我建议你们自己试着写一些
type SingleFileEntryVisitor struct {
	file *fileVisitor
}

func (s *SingleFileEntryVisitor) Get() File {
	panic("implement me")
}

func (s *SingleFileEntryVisitor) Visit(node ast.Node) ast.Visitor {
	panic("implement me")
}

type fileVisitor struct {
	ans     Annotations[*ast.File]
	types   []*typeVisitor
	visited bool
}

func (f *fileVisitor) Get() File {
	panic("implement me")
}

func (f *fileVisitor) Visit(node ast.Node) ast.Visitor {
	panic("implement me")
}

type File struct {
	Annotations[*ast.File]
	Types []Type
}

type typeVisitor struct {
	ans    Annotations[*ast.TypeSpec]
	fields []Field
}

func (t *typeVisitor) Get() Type {
	panic("implement me")
}

func (t *typeVisitor) Visit(node ast.Node) (w ast.Visitor) {
	panic("implement me")
}

type Type struct {
	Annotations[*ast.TypeSpec]
	Fields []Field
}

type Field struct {
	Annotations[*ast.Field]
}
