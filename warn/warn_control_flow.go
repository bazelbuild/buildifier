// Warnings related to the control flow

package warn

import (
	"fmt"
	"github.com/bazelbuild/buildtools/build"
	"github.com/bazelbuild/buildtools/bzlenv"
	"github.com/bazelbuild/buildtools/edit"
	"strings"
)

// findReturnsWithoutValue searches for return statements without a value, calls `callback` on
// them and returns whether the current list of statements terminates (either by a return or fail()
// statements on the current level in all subranches.
func findReturnsWithoutValue(stmts []build.Expr, callback func(*build.ReturnStmt)) bool {
	if len(stmts) == 0 {
		// May occur in empty else-clauses
		return false
	}
	terminated := false
	for _, stmt := range stmts {
		switch stmt := stmt.(type) {
		case *build.ReturnStmt:
			if stmt.Result == nil {
				callback(stmt)
			}
			terminated = true
		case *build.CallExpr:
			ident, ok := stmt.X.(*build.Ident)
			if ok && ident.Name == "fail" {
				terminated = true
			}
		case *build.ForStmt:
			// Call recursively to find all return statements without a value there.
			// Even if a for-loop is guaranteed to terminate in each iteration, buildifier still can't
			// check whether the loop is not empty, so we can't say that the statement after the ForStmt
			// is unreachable.
			findReturnsWithoutValue(stmt.Body, callback)
		case *build.IfStmt:
			// Save to separate values to avoid short circuit evaluation
			term1 := findReturnsWithoutValue(stmt.True, callback)
			term2 := findReturnsWithoutValue(stmt.False, callback)
			if term1 && term2 {
				terminated = true
			}
		}
	}
	return terminated
}

// missingReturnValueWarning warns if a function returns both explicit and implicit values.
func missingReturnValueWarning(f *build.File, fix bool) []*Finding {
	findings := []*Finding{}

	for _, stmt := range f.Stmt {
		function, ok := stmt.(*build.DefStmt)
		if !ok {
			continue
		}

		var hasNonEmptyReturns bool
		build.Walk(function, func(expr build.Expr, stack []build.Expr) {
			if ret, ok := expr.(*build.ReturnStmt); ok && ret.Result != nil {
				hasNonEmptyReturns = true
			}
		})

		if !hasNonEmptyReturns {
			continue
		}
		explicitReturn := findReturnsWithoutValue(function.Body, func(ret *build.ReturnStmt) {
			start, end := ret.Span()
			findings = append(findings,
				makeFinding(f, start, end, "return-value",
					`Some but not all execution paths of "`+function.Name+`" return a value.`, true, nil))
		})
		if !explicitReturn {
			start, end := function.Span()
			findings = append(findings,
				makeFinding(f, start, end, "return-value",
					`Some but not all execution paths of "`+function.Name+`" return a value.
The function may terminate by an implicit return in the end.`, true, nil))
		}
	}
	return findings
}

// findUnreachableStatements searches for unreachable statements (i.e. statements that immediately
// follow `return`, `break`, `continue`, and `fail()` statements and calls `callback` on them.
// If there are several consequent unreachable statements, it only reports the first of them.
// Returns whether the execution is terminated explicitly.
func findUnreachableStatements(stmts []build.Expr, callback func(build.Expr)) bool {
	unreachable := false
	for _, stmt := range stmts {
		if unreachable {
			callback(stmt)
			return true
		}
		switch stmt := stmt.(type) {
		case *build.ReturnStmt:
			unreachable = true
		case *build.CallExpr:
			ident, ok := stmt.X.(*build.Ident)
			if ok && ident.Name == "fail" {
				unreachable = true
			}
		case *build.Ident:
			switch stmt.Name {
			case "continue", "break":
				unreachable = true
			}
		case *build.ForStmt:
			findUnreachableStatements(stmt.Body, callback)
		case *build.IfStmt:
			// Save to separate values to avoid short circuit evaluation
			term1 := findUnreachableStatements(stmt.True, callback)
			term2 := findUnreachableStatements(stmt.False, callback)
			if term1 && term2 {
				unreachable = true
			}
		}
	}
	return unreachable
}

func unreachableStatementWarning(f *build.File, fix bool) []*Finding {
	findings := []*Finding{}

	for _, stmt := range f.Stmt {
		function, ok := stmt.(*build.DefStmt)
		if !ok {
			continue
		}

		findUnreachableStatements(function.Body, func(expr build.Expr) {
			start, end := expr.Span()
			findings = append(findings,
				makeFinding(f, start, end, "unreachable",
					`The statement is unreachable.`, true, nil))
		})
	}
	return findings
}

func noEffectStatementsCheck(f *build.File, body []build.Expr, isTopLevel, isFunc bool, findings []*Finding) []*Finding {
	seenNonComment := false
	for _, stmt := range body {
		start, end := stmt.Span()
		if _, ok := stmt.(*build.StringExpr); ok {
			if !seenNonComment && (isTopLevel || isFunc) {
				// It's a docstring.
				seenNonComment = true
				continue
			}
		}
		if _, ok := stmt.(*build.CommentBlock); !ok {
			seenNonComment = true
		}
		switch s := (stmt).(type) {
		case *build.DefStmt, *build.ForStmt, *build.IfStmt, *build.LoadStmt, *build.ReturnStmt,
			*build.CallExpr, *build.CommentBlock:
			continue
		case *build.BinaryExpr:
			if s.Op != "==" && s.Op != "!=" && strings.HasSuffix(s.Op, "=") {
				continue
			}
		case *build.Ident:
			if s.Name == "break" || s.Name == "continue" || s.Name == "pass" {
				continue
			}
		}
		if comp, ok := stmt.(*build.Comprehension); ok {
			if !isTopLevel || comp.Curly {
				// List comprehensions are allowed on top-level.
				findings = append(findings,
					makeFinding(f, start, end, "no-effect",
						"Expression result is not used. Use a for-loop instead of a list comprehension.", true, nil))
			}
			continue
		}
		findings = append(findings,
			makeFinding(f, start, end, "no-effect",
				"Expression result is not used.", true, nil))
	}
	return findings
}

func noEffectWarning(f *build.File, fix bool) []*Finding {
	findings := []*Finding{}
	findings = noEffectStatementsCheck(f, f.Stmt, true, false, findings)
	build.Walk(f, func(expr build.Expr, stack []build.Expr) {
		// The AST should have a ExprStmt node.
		// Since we don't have that, we match on the nodes that contain a block to get the list of statements.
		switch expr := expr.(type) {
		case *build.ForStmt:
			findings = noEffectStatementsCheck(f, expr.Body, false, false, findings)
		case *build.DefStmt:
			findings = noEffectStatementsCheck(f, expr.Function.Body, false, true, findings)
		case *build.IfStmt:
			findings = noEffectStatementsCheck(f, expr.True, false, false, findings)
			findings = noEffectStatementsCheck(f, expr.False, false, false, findings)
		}
	})
	return findings
}

// unusedVariableCheck checks for unused variables inside a given node `stmt` (either *build.File or
// *build.DefStmt) and reports unused and already defined variables.
func unusedVariableCheck(f *build.File, stmts []build.Expr, findings []*Finding) []*Finding {
	if f.Type == build.TypeDefault {
		// Not applicable to .bzl files, unused symbols may be loaded and used in other files.
		return findings
	}
	usedSymbols := make(map[string]bool)

	for _, stmt := range stmts {
		for key := range edit.UsedSymbols(stmt) {
			usedSymbols[key] = true
		}
	}

	for _, s := range stmts {
		if defStmt, ok := s.(*build.DefStmt); ok {
			findings = unusedVariableCheck(f, defStmt.Body, findings)
			continue
		}

		// look for all assignments in the scope
		as, ok := s.(*build.BinaryExpr)
		if !ok || as.Op != "=" {
			continue
		}
		start, end := as.X.Span()
		left, ok := as.X.(*build.Ident)
		if !ok {
			continue
		}
		if usedSymbols[left.Name] {
			continue
		}
		if edit.ContainsComments(s, "@unused") {
			// To disable the warning, put a comment that contains '@unused'
			continue
		}
		findings = append(findings,
			makeFinding(f, start, end, "unused-variable",
				"Variable \""+left.Name+"\" is unused. Please remove it.\n"+
					"To disable the warning, add '@unused' in a comment.", true, nil))
	}
	return findings
}

func unusedVariableWarning(f *build.File, fix bool) []*Finding {
	return unusedVariableCheck(f, f.Stmt, []*Finding{})
}

func redefinedVariableWarning(f *build.File, fix bool) []*Finding {
	findings := []*Finding{}
	definedSymbols := make(map[string]bool)

	for _, s := range f.Stmt {
		// look for all assignments in the scope
		as, ok := s.(*build.BinaryExpr)
		if !ok || as.Op != "=" {
			continue
		}
		start, end := as.X.Span()
		left, ok := as.X.(*build.Ident)
		if !ok {
			continue
		}
		if definedSymbols[left.Name] {
			findings = append(findings,
				makeFinding(f, start, end, "redefined-variable",
					"Variable \""+left.Name+"\" has already been defined. "+
						"Redefining a global value is discouraged and will be forbidden in the future.\n"+
						"Consider using a new variable instead.", true, nil))
			continue
		}
		definedSymbols[left.Name] = true
	}
	return findings
}

func unusedLoadWarning(f *build.File, fix bool) []*Finding {
	findings := []*Finding{}
	loaded := make(map[string]struct{ label, from string })

	symbols := edit.UsedSymbols(f)
	for stmtIndex := 0; stmtIndex < len(f.Stmt); stmtIndex++ {
		load, ok := f.Stmt[stmtIndex].(*build.LoadStmt)
		if !ok {
			continue
		}
		for i := 0; i < len(load.To); i++ {
			from := load.From[i]
			to := load.To[i]
			// Check if the symbol was already loaded
			origin, alreadyLoaded := loaded[to.Name]
			loaded[to.Name] = struct{ label, from string }{load.Module.Token, from.Name}

			if alreadyLoaded {
				if fix && origin.label == load.Module.Token && origin.from == from.Name {
					// Only fix if it's loaded from the label and variable
					load.To = append(load.To[:i], load.To[i+1:]...)
					load.From = append(load.From[:i], load.From[i+1:]...)
					i--
				} else {
					start, end := to.Span()
					findings = append(findings,
						makeFinding(f, start, end, "load",
							"Symbol \""+to.Name+"\" has already been loaded. Please remove it.", true, nil))
				}
				continue
			}
			_, ok := symbols[to.Name]
			if !ok && !edit.ContainsComments(load, "@unused") && !edit.ContainsComments(to, "@unused") && !edit.ContainsComments(from, "@unused") {
				// To disable the warning, put a comment that contains '@unused'
				if fix {
					load.To = append(load.To[:i], load.To[i+1:]...)
					load.From = append(load.From[:i], load.From[i+1:]...)
					i--
				} else {
					start, end := to.Span()
					findings = append(findings,
						makeFinding(f, start, end, "load",
							"Loaded symbol \""+to.Name+"\" is unused. Please remove it.\n"+
								"To disable the warning, add '@unused' in a comment.", true, nil))

				}
			}
		}
		// If there are no loaded symbols left remove the entire load statement
		if fix && len(load.To) == 0 {
			f.Stmt = append(f.Stmt[:stmtIndex], f.Stmt[stmtIndex+1:]...)
		}
	}
	return findings
}

// collectLocalVariables traverses statements (e.g. of a function definition) and returns a map of
// variables defined anywhere inside the function.
func collectLocalVariables(stmts []build.Expr) map[string]bool {
	variables := make(map[string]bool)

	for _, stmt := range stmts {
		switch stmt := stmt.(type) {
		case *build.DefStmt:
			// Don't traverse nested functions
		case *build.ForStmt:
			for _, key := range bzlenv.CollectLValues(stmt.Vars) {
				variables[key.Name] = true
			}
			for key := range collectLocalVariables(stmt.Body) {
				variables[key] = true
			}
		case *build.IfStmt:
			for key := range collectLocalVariables(stmt.True) {
				variables[key] = true
			}
			for key := range collectLocalVariables(stmt.False) {
				variables[key] = true
			}
		case *build.BinaryExpr:
			if stmt.Op == "=" {
				for _, key := range bzlenv.CollectLValues(stmt.X) {
					variables[key.Name] = true
				}
			}
		}
	}
	return variables
}

// searchUninitializedVariables takes a list of statements (e.g. body of a block statement)
// and a map of previously initialized statements, and calls `callback` on all idents that are not
// initialized. An ident is considered initialized if it's initialized by every possible execution
// path (before or by `stmts`).
// Returns variables that are guaranteed to be defined by `stmts`.
func findUninitializedVariables(stmts []build.Expr, previouslyInitialized map[string]bool, callback func(*build.Ident)) map[string]bool {
	// Variables that are guaranteed to be de initialized
	locallyInitialized := make(map[string]bool) // in the local block of `stmts`
	initialized := make(map[string]bool)        // anywhere before the current line
	for key := range previouslyInitialized {
		initialized[key] = true
	}

	// findUninitializedIdents traverses an expression (simple statement or a part of it), and calls
	// `callback` on every *build.Ident that's not mentioned in the map of initialized variables
	findUninitializedIdents := func(expr build.Expr, callback func(ident *build.Ident)) {
		build.Walk(expr, func(expr build.Expr, stack []build.Expr) {
			if ident, ok := expr.(*build.Ident); ok && !initialized[ident.Name] {
				callback(ident)
			}
		})
	}

	for _, stmt := range stmts {
		switch stmt := stmt.(type) {
		case *build.DefStmt:
			// Don't traverse nested functions
		case *build.ForStmt:
			// Although loop variables are defined as local variables, buildifier doesn't know whether
			// the loop will be empty or not

			// Traverse but ignore the result. Even if something is defined inside a for-loop, the loop
			// may be empty and the variable initialization may not happen.
			findUninitializedIdents(stmt.X, callback)
			findUninitializedVariables(stmt.Body, initialized, callback)
		case *build.IfStmt:
			findUninitializedIdents(stmt.Cond, callback)
			// If a variable is defined in both if- and else-clauses, it's considered as defined
			definedInTrue := findUninitializedVariables(stmt.True, initialized, callback)
			for key := range findUninitializedVariables(stmt.False, initialized, callback) {
				if definedInTrue[key] {
					locallyInitialized[key] = true
					initialized[key] = true
				}
			}
		case *build.BinaryExpr:
			if stmt.Op == "=" {
				// Assignment expression. Collect all definitions from the lhs (they shouldn't be taken into
				// account while checking for undefined usages.
				lValues := make(map[*build.Ident]bool)
				for _, ident := range bzlenv.CollectLValues(stmt.X) {
					lValues[ident] = true
				}
				// Traverse the statement and report all undefined idents expect LValues
				findUninitializedIdents(stmt, func(ident *build.Ident) {
					if !lValues[ident] {
						callback(ident)
					}
				})
				// Update locallyInitialized and defined with newly defined variables
				for ident := range lValues {
					locallyInitialized[ident.Name] = true
					initialized[ident.Name] = true
				}
			} else {
				findUninitializedIdents(stmt, callback)
			}
		default:
			findUninitializedIdents(stmt, callback)
		}
	}
	return locallyInitialized
}

// uninitializedVariableWarning warns about usages of values that may not have been initialized.
func uninitializedVariableWarning(f *build.File, _ bool) []*Finding {
	findings := []*Finding{}
	for _, stmt := range f.Stmt {
		def, ok := stmt.(*build.DefStmt)
		if !ok {
			continue
		}

		// Get all variables defined in the function body.
		// If a variable is not defined there, it can be builtin, global, or loaded.
		localVars := collectLocalVariables(def.Body)

		// Function parameters are guaranteed to be defined everywhere in the function, even if they
		// are redefined inside the function body. They shouldn't be taken into consideration.
		for _, node := range def.Params {
			switch node := node.(type) {
			case *build.Ident:
				delete(localVars, node.Name)
			case *build.UnaryExpr:
				// either *args or **kwargs
				if ident, ok := node.X.(*build.Ident); ok {
					delete(localVars, ident.Name)
				}
			case *build.BinaryExpr:
				// x = value
				if ident, ok := node.X.(*build.Ident); ok {
					delete(localVars, ident.Name)
				}
			}
		}

		// Search for all potentially initialized variables in the function body
		findUninitializedVariables(def.Body, make(map[string]bool), func(ident *build.Ident) {
			// Check that the found ident represents a local variable
			if localVars[ident.Name] {
				start, end := ident.Span()
				findings = append(findings,
					makeFinding(f, start, end, "uninitialized",
						fmt.Sprintf(`Variable "%s" may not have been initialized.`, ident.Name), true, nil))
			}
		})
	}
	return findings
}
