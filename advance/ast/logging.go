package ast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
	"sync"
)

type loggingVisitor struct {
	file *ast.File
	addLogOnce sync.Once
}

func (l *loggingVisitor) Visit(node ast.Node) (w ast.Visitor) {
	switch n := node.(type) {
	case *ast.File:
		l.file = n
	case *ast.FuncDecl:
		input := n.Type.Params.List
		if len(input) == 0 {
			return l
		}
		logging, ok := l.logStmt(input)
		if !ok {
			return l
		}
		expr, err := parser.ParseExpr(logging.String())
		if err != nil {
			fmt.Println(err)
			return l
		}
		l.addLogOnce.Do(func() {
			l.addImport(l.file)
		})
		if len(n.Body.List) > 0 {
			res := make([]ast.Stmt, 0, len(n.Body.List)+1)
			res = append(res, &ast.ExprStmt{X: expr})
			res = append(res, n.Body.List...)
			n.Body.List = res
		}
	}
	return l
}

func (l *loggingVisitor) logStmt(input []*ast.Field) (strings.Builder, bool) {
	var logging strings.Builder
	logging.WriteString("log.Println(")
	names := make([]string, 0, len(input))
	for _, fd := range input {
		if len(fd.Names) == 0 {
			return strings.Builder{}, false
		}
		for _, name := range fd.Names {
			names = append(names, name.Name)
		}
	}
	for i, name := range names {
		if i > 0 {
			logging.WriteByte(',')
		}
		logging.WriteString(name)
	}

	logging.WriteByte(')')
	return logging, true
}

func (l *loggingVisitor) addImport(n *ast.File) {
	newImport := &ast.ImportSpec{
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: strconv.Quote("log"),
		},
	}

	var importDecl *ast.GenDecl
	for _, decl := range n.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			if d.Tok == token.IMPORT {
				importDecl = d
			}
		}
	}

	if importDecl == nil {
		importDecl = &ast.GenDecl{
			Tok: token.IMPORT,
		}
		n.Decls = append([]ast.Decl{importDecl}, n.Decls...)
	}
	if len(n.Decls) == 0 {
		n.Decls = make([]ast.Decl, 0, 1)
	}
	importDecl.Specs = append(importDecl.Specs, newImport)
}
