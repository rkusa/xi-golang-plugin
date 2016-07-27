// This code is partially based on Go's RPC implementation
// https://golang.org/src/go/ast/walk.go

package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"reflect"
)

func walk(node ast.Node, spans chan *Span) {
	log.Println(node, reflect.TypeOf(node), node.Pos(), node.End())

	// walk children
	// (the order of the cases matches the order
	// of the corresponding node types in ast.go)
	switch n := node.(type) {
	// Comments and fields
	case *ast.Comment:
		// nothing to do

	case *ast.CommentGroup:
		for _, c := range n.List {
			walk(c, spans)
		}

	case *ast.Field:
		if n.Doc != nil {
			walk(n.Doc, spans)
		}

		// walkIdentList(n.Names, spans)
		for _, x := range n.Names {
			spans <- span(x, color3)
		}

		walk(n.Type, spans)
		if n.Tag != nil {
			walk(n.Tag, spans)
		}
		if n.Comment != nil {
			walk(n.Comment, spans)
		}

	case *ast.FieldList:
		for _, f := range n.List {
			walk(f, spans)
		}

	// Expressions
	case *ast.BadExpr, *ast.Ident:
		spans <- &Span{node.Pos(), node.End(), color1}

	case *ast.BasicLit:
		spans <- &Span{node.Pos(), node.End(), color2}

	case *ast.Ellipsis:
		if n.Elt != nil {
			walk(n.Elt, spans)
		}

	case *ast.FuncLit:
		walk(n.Type, spans)
		walk(n.Body, spans)

	case *ast.CompositeLit:
		if n.Type != nil {
			walk(n.Type, spans)
		}
		walkExprList(n.Elts, spans)

	case *ast.ParenExpr:
		walk(n.X, spans)

	case *ast.SelectorExpr:
		walk(n.X, spans)
		walk(n.Sel, spans)

	case *ast.IndexExpr:
		walk(n.X, spans)
		walk(n.Index, spans)

	case *ast.SliceExpr:
		walk(n.X, spans)
		if n.Low != nil {
			walk(n.Low, spans)
		}
		if n.High != nil {
			walk(n.High, spans)
		}
		if n.Max != nil {
			walk(n.Max, spans)
		}

	case *ast.TypeAssertExpr:
		walk(n.X, spans)
		if n.Type != nil {
			walk(n.Type, spans)
		}

	case *ast.CallExpr:
		walk(n.Fun, spans)
		walkExprList(n.Args, spans)

	case *ast.StarExpr:
		walk(n.X, spans)

	case *ast.UnaryExpr:
		walk(n.X, spans)

	case *ast.BinaryExpr:
		walk(n.X, spans)
		walk(n.Y, spans)

	case *ast.KeyValueExpr:
		walk(n.Key, spans)
		walk(n.Value, spans)

	// Types
	case *ast.ArrayType:
		if n.Len != nil {
			walk(n.Len, spans)
		}
		walk(n.Elt, spans)

	case *ast.StructType:
		// highlight `struct`
		spans <- &Span{n.Struct, n.Struct + 6, color5}

		walk(n.Fields, spans)

	case *ast.FuncType:
		if n.Params != nil {
			walk(n.Params, spans)
		}
		if n.Results != nil {
			walk(n.Results, spans)
		}

	case *ast.InterfaceType:
		// highlight `interface`
		spans <- &Span{n.Interface, n.Interface + 9, color5}

		walk(n.Methods, spans)

	case *ast.MapType:
		walk(n.Key, spans)
		walk(n.Value, spans)

	case *ast.ChanType:
		walk(n.Value, spans)

	// Statements
	case *ast.BadStmt:
		// nothing to do

	case *ast.DeclStmt:
		walk(n.Decl, spans)

	case *ast.EmptyStmt:
		// nothing to do

	case *ast.LabeledStmt:
		walk(n.Label, spans)
		walk(n.Stmt, spans)

	case *ast.ExprStmt:
		walk(n.X, spans)

	case *ast.SendStmt:
		walk(n.Chan, spans)
		walk(n.Value, spans)

	case *ast.IncDecStmt:
		walk(n.X, spans)

	case *ast.AssignStmt:
		walkExprList(n.Lhs, spans)
		walkExprList(n.Rhs, spans)

	case *ast.GoStmt:
		walk(n.Call, spans)

	case *ast.DeferStmt:
		walk(n.Call, spans)

	case *ast.ReturnStmt:
		walkExprList(n.Results, spans)

	case *ast.BranchStmt:
		if n.Label != nil {
			walk(n.Label, spans)
		}

	case *ast.BlockStmt:
		walkStmtList(n.List, spans)

	case *ast.IfStmt:
		if n.Init != nil {
			walk(n.Init, spans)
		}
		walk(n.Cond, spans)
		walk(n.Body, spans)
		if n.Else != nil {
			walk(n.Else, spans)
		}

	case *ast.CaseClause:
		walkExprList(n.List, spans)
		walkStmtList(n.Body, spans)

	case *ast.SwitchStmt:
		if n.Init != nil {
			walk(n.Init, spans)
		}
		if n.Tag != nil {
			walk(n.Tag, spans)
		}
		walk(n.Body, spans)

	case *ast.TypeSwitchStmt:
		if n.Init != nil {
			walk(n.Init, spans)
		}
		walk(n.Assign, spans)
		walk(n.Body, spans)

	case *ast.CommClause:
		if n.Comm != nil {
			walk(n.Comm, spans)
		}
		walkStmtList(n.Body, spans)

	case *ast.SelectStmt:
		walk(n.Body, spans)

	case *ast.ForStmt:
		if n.Init != nil {
			walk(n.Init, spans)
		}
		if n.Cond != nil {
			walk(n.Cond, spans)
		}
		if n.Post != nil {
			walk(n.Post, spans)
		}
		walk(n.Body, spans)

	case *ast.RangeStmt:
		if n.Key != nil {
			walk(n.Key, spans)
		}
		if n.Value != nil {
			walk(n.Value, spans)
		}
		walk(n.X, spans)
		walk(n.Body, spans)

	// Declarations
	case *ast.ImportSpec:
		if n.Doc != nil {
			walk(n.Doc, spans)
		}
		if n.Name != nil {
			walk(n.Name, spans)
		}
		walk(n.Path, spans)
		if n.Comment != nil {
			walk(n.Comment, spans)
		}

	case *ast.ValueSpec:
		if n.Doc != nil {
			walk(n.Doc, spans)
		}
		walkIdentList(n.Names, spans)
		if n.Type != nil {
			walk(n.Type, spans)
		}
		walkExprList(n.Values, spans)
		if n.Comment != nil {
			walk(n.Comment, spans)
		}

	case *ast.TypeSpec:
		if n.Doc != nil {
			walk(n.Doc, spans)
		}

		// walk(n.Name, spans)
		spans <- span(n.Name, color3)

		walk(n.Type, spans)
		if n.Comment != nil {
			walk(n.Comment, spans)
		}

	case *ast.BadDecl:
		// nothing to do

	case *ast.GenDecl:
		// highlight `type`
		var l token.Pos = 0
		switch n.Tok {
		case token.IMPORT:
			l = 6
		case token.CONST:
			l = 5
		case token.TYPE:
			l = 4
		case token.VAR:
			l = 3
		}
		spans <- &Span{n.TokPos, n.TokPos + l, color5}

		if n.Doc != nil {
			walk(n.Doc, spans)
		}
		for _, s := range n.Specs {
			walk(s, spans)
		}

	case *ast.FuncDecl:
		if n.Doc != nil {
			walk(n.Doc, spans)
		}
		if n.Recv != nil {
			walk(n.Recv, spans)
		}
		walk(n.Name, spans)
		walk(n.Type, spans)
		if n.Body != nil {
			walk(n.Body, spans)
		}

	// Files and packages
	case *ast.File:
		// highlight `package`
		spans <- &Span{node.Pos(), node.Pos() + 7, color5}

		if n.Doc != nil {
			walk(n.Doc, spans)
		}
		walk(n.Name, spans)
		walkDeclList(n.Decls, spans)
		// don't walk n.Comments - they have been
		// visited already through the individual
		// nodes

	case *ast.Package:
		for _, f := range n.Files {
			walk(f, spans)
		}

	default:
		panic(fmt.Sprintf("walk: unexpected node type %T", n))
	}
}

// Helper functions for common node lists. They may be empty.

func walkIdentList(list []*ast.Ident, spans chan *Span) {
	for _, x := range list {
		walk(x, spans)
	}
}

func walkExprList(list []ast.Expr, spans chan *Span) {
	for _, x := range list {
		walk(x, spans)
	}
}

func walkStmtList(list []ast.Stmt, spans chan *Span) {
	for _, x := range list {
		walk(x, spans)
	}
}

func walkDeclList(list []ast.Decl, spans chan *Span) {
	for _, x := range list {
		walk(x, spans)
	}
}

func span(n ast.Node, c int64) *Span {
	return &Span{n.Pos(), n.End(), c}
}
