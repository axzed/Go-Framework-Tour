package annotation

import (
	"go/ast"
)

type SingleFileEntryVisitor struct {
	file *fileVisitor
}

func (s *SingleFileEntryVisitor) Get() File {
	if s.file != nil {
		return s.file.Get()
	}
	return File{}
}

func (s *SingleFileEntryVisitor) Visit(node ast.Node) ast.Visitor {
	file, ok := node.(*ast.File)
	if ok {
		s.file = &fileVisitor{
			ans: newAnnotations[*ast.File](file, file.Doc),
		}
		return s.file
	}
	return s
}

type fileVisitor struct {
	ans   Annotations[*ast.File]
	types []*typeVisitor
	visited bool
}

func (f *fileVisitor) Get() File {
	types := make([]Type, 0, len(f.types))
	for _, t := range f.types {
		types = append(types, t.Get())
	}
	return File{
		Annotations: f.ans,
		Types:       types,
	}
}

func (f *fileVisitor) Visit(node ast.Node) ast.Visitor {
	typ, ok := node.(*ast.TypeSpec)
	if ok {
		res := &typeVisitor{
			ans:    newAnnotations[*ast.TypeSpec](typ, typ.Doc),
			fields: make([]Field, 0, 0),
		}
		f.types = append(f.types, res)
		return res
	}
	return f
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
	return Type{
		Annotations: t.ans,
		Fields:      t.fields,
	}
}

func (t *typeVisitor) Visit(node ast.Node) (w ast.Visitor) {
	fd, ok := node.(*ast.Field)
	if ok {
		t.fields = append(t.fields, Field{Annotations: newAnnotations[*ast.Field](fd, fd.Doc)})
		return nil
	}
	return t
}

type Type struct {
	Annotations[*ast.TypeSpec]
	Fields []Field
}

type Field struct {
	Annotations[*ast.Field]
}
