// This file contains functions to convert from one AST to the other.
// Input: AST from github.com/google/skylark/syntax
// Output: AST from github.com/bazelbuild/buildtools/build

package main

import (
	"strconv"
	"strings"

	"github.com/bazelbuild/buildtools/build"
	"github.com/google/skylark/syntax"
)

func convFile(f *syntax.File) *build.File {
	stmts := []build.Expr{}
	for _, stmt := range f.Stmts {
		stmts = append(stmts, convStmt(stmt))
	}

	return &build.File{
		Stmt:     stmts,
		Comments: convComments(f.Comments()),
	}
}

func convStmt(stmt syntax.Stmt) build.Expr {
	switch stmt := stmt.(type) {
	case *syntax.ExprStmt:
		s := convExpr(stmt.X)
		*s.Comment() = convComments(stmt.Comments())
		return s
	case *syntax.BranchStmt:
		return &build.LiteralExpr{
			Token:    stmt.Token.String(),
			Comments: convComments(stmt.Comments()),
		}
	case *syntax.LoadStmt:
		args := []build.Expr{}
		args = append(args, &build.StringExpr{Value: stmt.Module.Value.(string)})
		for i, _ := range stmt.From {
			if stmt.From[i].Name == stmt.To[i].Name {
				args = append(args, &build.StringExpr{Value: stmt.To[i].Name})
			} else {
				args = append(args, &build.BinaryExpr{
					Op: "=",
					X:  &build.LiteralExpr{Token: stmt.From[i].Name},
					Y:  &build.StringExpr{Value: stmt.To[i].Name}})
			}
		}
		return &build.CallExpr{
			X:            &build.LiteralExpr{Token: "load"},
			List:         args,
			ForceCompact: singleLine(stmt),
		}
	case *syntax.AssignStmt:
		return &build.BinaryExpr{
			Op:       stmt.Op.String(),
			X:        convExpr(stmt.LHS),
			Y:        convExpr(stmt.RHS),
			Comments: convComments(stmt.Comments()),
		}
	case *syntax.IfStmt:
		return &build.IfStmt{
			Cond:     convExpr(stmt.Cond),
			True:     convStmts(stmt.True),
			False:    convStmts(stmt.False),
			Comments: convComments(stmt.Comments()),
		}
	case *syntax.DefStmt:
		return &build.DefStmt{
			Name:     stmt.Name.Name,
			Comments: convComments(stmt.Comments()),
			Function: build.Function{
				Params: convExprs(stmt.Function.Params),
				Body:   convStmts(stmt.Function.Body),
			},
		}
	case *syntax.ForStmt:
		return &build.ForStmt{
			Vars:     convExpr(stmt.Vars),
			X:        convExpr(stmt.X),
			Comments: convComments(stmt.Comments()),
			Body:     convStmts(stmt.Body),
		}
	case *syntax.ReturnStmt:
		return &build.ReturnStmt{
			Comments: convComments(stmt.Comments()),
			Result:   convExpr(stmt.Result),
		}
	}
	panic("unreachable")
}

func convStmts(list []syntax.Stmt) []build.Expr {
	res := []build.Expr{}
	for _, i := range list {
		res = append(res, convStmt(i))
	}
	return res
}

func convExprs(list []syntax.Expr) []build.Expr {
	res := []build.Expr{}
	for _, i := range list {
		res = append(res, convExpr(i))
	}
	return res
}

func convCommentList(list []syntax.Comment, txt string) []build.Comment {
	res := []build.Comment{}
	for _, c := range list {
		res = append(res, build.Comment{Token: c.Text})
	}
	return res
}

func convComments(c *syntax.Comments) build.Comments {
	if c == nil {
		return build.Comments{}
	}
	return build.Comments{
		Before: convCommentList(c.Before, "before"),
		Suffix: convCommentList(c.Suffix, "suffix"),
		After:  convCommentList(c.After, "after"),
	}
}

// singleLine returns true if the node fits on a single line.
func singleLine(n syntax.Node) bool {
	start, end := n.Span()
	return start.Line == end.Line
}

func convClauses(list []syntax.Node) []*build.ForClauseWithIfClausesOpt {
	res := []*build.ForClauseWithIfClausesOpt{}
	var current *build.ForClauseWithIfClausesOpt
	for _, c := range list {
		switch stmt := c.(type) {
		case *syntax.ForClause:
			if current != nil {
				res = append(res, current)
			}
			f := &build.ForClause{
				Var:  convExpr(stmt.Vars),
				Expr: convExpr(stmt.X),
			}
			current = &build.ForClauseWithIfClausesOpt{For: f}
		case *syntax.IfClause:
			current.Ifs = append(current.Ifs, &build.IfClause{Cond: convExpr(stmt.Cond)})
		}
	}
	res = append(res, current)
	return res
}

func convExpr(e syntax.Expr) build.Expr {
	if e == nil {
		return nil
	}
	switch e := e.(type) {
	case *syntax.Literal:
		switch e.Token {
		case syntax.INT:
			return &build.LiteralExpr{
				Token:    strconv.FormatInt(e.Value.(int64), 10),
				Comments: convComments(e.Comments())}
		case syntax.FLOAT:
			panic("float")
		case syntax.STRING:
			return &build.StringExpr{
				Value:       e.Value.(string),
				TripleQuote: strings.HasPrefix(e.Raw, "\"\"\""),
				Comments:    convComments(e.Comments())}
		}
	case *syntax.Ident:
		return &build.LiteralExpr{Token: e.Name, Comments: convComments(e.Comments())}
	case *syntax.BinaryExpr:
		_, lhsEnd := e.X.Span()
		rhsBegin, _ := e.Y.Span()
		return &build.BinaryExpr{
			X:         convExpr(e.X),
			Y:         convExpr(e.Y),
			Op:        e.Op.String(),
			LineBreak: lhsEnd.Line != rhsBegin.Line,
			Comments:  convComments(e.Comments())}
	case *syntax.UnaryExpr:
		return &build.UnaryExpr{Op: e.Op.String(), X: convExpr(e.X)}
	case *syntax.SliceExpr:
		return &build.SliceExpr{X: convExpr(e.X), From: convExpr(e.Lo), To: convExpr(e.Hi), Step: convExpr(e.Step)}
	case *syntax.DotExpr:
		return &build.DotExpr{X: convExpr(e.X), Name: e.Name.Name}
	case *syntax.CallExpr:
		args := []build.Expr{}
		for _, a := range e.Args {
			args = append(args, convExpr(a))
		}
		return &build.CallExpr{
			X:            convExpr(e.Fn),
			List:         args,
			ForceCompact: singleLine(e),
		}
	case *syntax.ListExpr:
		list := []build.Expr{}
		for _, i := range e.List {
			list = append(list, convExpr(i))
		}
		return &build.ListExpr{List: list, Comments: convComments(e.Comments())}
	case *syntax.DictExpr:
		list := []build.Expr{}
		for i, _ := range e.List {
			entry := e.List[i].(*syntax.DictEntry)
			list = append(list, &build.KeyValueExpr{
				Key:      convExpr(entry.Key),
				Value:    convExpr(entry.Value),
				Comments: convComments(entry.Comments()),
			})
		}
		return &build.DictExpr{List: list, Comments: convComments(e.Comments())}
	case *syntax.CondExpr:
		return &build.ConditionalExpr{
			Then:     convExpr(e.True),
			Test:     convExpr(e.Cond),
			Else:     convExpr(e.False),
			Comments: convComments(e.Comments()),
		}
	case *syntax.Comprehension:
		brack := "[]"
		if e.Curly {
			brack = "{}"
		}
		return &build.ListForExpr{
			X:        convExpr(e.Body),
			For:      convClauses(e.Clauses),
			Comments: convComments(e.Comments()),
			Brack:    brack,
		}
	case *syntax.ParenExpr:
		return &build.ParenExpr{
			X:        convExpr(e.X),
			Comments: convComments(e.Comments()),
		}
	case *syntax.TupleExpr:
		return &build.TupleExpr{
			List:         convExprs(e.List),
			Comments:     convComments(e.Comments()),
			ForceCompact: singleLine(e),
		}
	case *syntax.IndexExpr:
		return &build.IndexExpr{
			X:        convExpr(e.X),
			Y:        convExpr(e.Y),
			Comments: convComments(e.Comments()),
		}
	}
	panic("other expr")
}
