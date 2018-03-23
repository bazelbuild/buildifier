// BUILD file parser.

// This is a yacc grammar. Its lexer is in lex.go.
//
// For a good introduction to writing yacc grammars, see
// Kernighan and Pike's book The Unix Programming Environment.
//
// The definitive yacc manual is
// Stephen C. Johnson and Ravi Sethi, "Yacc: A Parser Generator",
// online at http://plan9.bell-labs.com/sys/doc/yacc.pdf.

%{
package build
%}

// The generated parser puts these fields in a struct named yySymType.
// (The name %union is historical, but it is inaccurate for Go.)
%union {
	// input tokens
	tok       string     // raw input syntax
	str       string     // decoding of quoted string
	pos       Position   // position of token
	triple    bool       // was string triple quoted?

	// partial syntax trees
	expr      Expr
	exprs     []Expr
	forc      *ForClause
	ifs       []*IfClause
	forifs    *ForClauseWithIfClausesOpt
	forsifs   []*ForClauseWithIfClausesOpt
	string    *StringExpr
	strings   []*StringExpr
	ifstmt    *IfStmt

	// supporting information
	comma     Position   // position of trailing comma in list, if present
	lastRule  Expr  // most recent rule, to attach line comments to
}

// These declarations set the type for a $ reference ($$, $1, $2, ...)
// based on the kind of symbol it refers to. Other fields can be referred
// to explicitly, as in $<tok>1.
//
// %token is for input tokens generated by the lexer.
// %type is for higher-level grammar rules defined here.
//
// It is possible to put multiple tokens per line, but it is easier to
// keep ordered using a sparser one-per-line list.

%token	<pos>	'%'
%token	<pos>	'('
%token	<pos>	')'
%token	<pos>	'*'
%token	<pos>	'+'
%token	<pos>	','
%token	<pos>	'-'
%token	<pos>	'.'
%token	<pos>	'/'
%token	<pos>	':'
%token	<pos>	'<'
%token	<pos>	'='
%token	<pos>	'>'
%token	<pos>	'['
%token	<pos>	']'
%token	<pos>	'{'
%token	<pos>	'}'

// By convention, yacc token names are all caps.
// However, we do not want to export them from the Go package
// we are creating, so prefix them all with underscores.

%token	<pos>	_AUGM    // augmented assignment
%token	<pos>	_AND     // keyword and
%token	<pos>	_COMMENT // top-level # comment
%token	<pos>	_EOF     // end of file
%token	<pos>	_EQ      // operator ==
%token	<pos>	_FOR     // keyword for
%token	<pos>	_GE      // operator >=
%token	<pos>	_IDENT   // non-keyword identifier or number
%token	<pos>	_IF      // keyword if
%token	<pos>	_ELSE    // keyword else
%token	<pos>	_ELIF    // keyword elif
%token	<pos>	_IN      // keyword in
%token	<pos>	_IS      // keyword is
%token	<pos>	_LAMBDA  // keyword lambda
%token	<pos>	_LOAD    // keyword load
%token	<pos>	_LE      // operator <=
%token	<pos>	_NE      // operator !=
%token	<pos>	_STAR_STAR // operator **
%token	<pos>	_NOT     // keyword not
%token	<pos>	_OR      // keyword or
%token	<pos>	_PYTHON  // uninterpreted Python block
%token	<pos>	_STRING  // quoted string
%token	<pos>	_DEF     // keyword def
%token	<pos>	_RETURN  // keyword return
%token	<pos>	_INDENT  // indentation
%token	<pos>	_UNINDENT // unindentation

%type	<pos>		comma_opt
%type	<expr>		argument
%type	<exprs>		arguments
%type	<exprs>		arguments_opt
%type	<expr>		parameter
%type	<exprs>		parameters
%type	<exprs>		parameters_opt
%type	<expr>		test
%type	<expr>		test_opt
%type	<exprs>		tests_opt
%type	<expr>		primary_expr
%type	<expr>		expr
%type	<expr>		expr_opt
%type <exprs>		tests
%type	<exprs>		exprs
%type	<exprs>		exprs_opt
%type	<exprs>		loop_vars
%type	<forc>		for_clause
%type	<forifs>	for_clause_with_if_clauses_opt
%type	<forsifs>	for_clauses_with_if_clauses_opt
%type	<expr>		ident
%type	<ifs>		if_clauses_opt
%type	<exprs>		stmts
%type	<exprs>		stmt          // a simple_stmt or a for/if/def block
%type	<expr>		block_stmt    // a single for/if/def statement
%type	<ifstmt>	if_else_block // a complete if-elif-else block
%type	<ifstmt>	if_block      // a single if block
%type	<ifstmt>	else_block    // a single else block
%type	<ifstmt>	elif_chain    // an elif-elif-else chain
%type <pos>		elif          // `elif` or `else if` token(s)
%type	<exprs>		simple_stmt   // One or many small_stmts on one line, e.g. 'a = f(x); return str(a)'
%type	<expr>		small_stmt    // A single statement, e.g. 'a = f(x)'
%type <exprs>		small_stmts_continuation  // A sequence of `';' small_stmt`
%type	<expr>		keyvalue
%type	<exprs>		keyvalues
%type	<exprs>		keyvalues_no_comma
%type	<string>	string
%type	<strings>	strings
%type	<exprs>		suite
%type	<exprs>		comments

// Operator precedence.
// Operators listed lower in the table bind tighter.

// We tag rules with this fake, low precedence to indicate
// that when the rule is involved in a shift/reduce
// conflict, we prefer that the parser shift (try for a longer parse).
// Shifting is the default resolution anyway, but stating it explicitly
// silences yacc's warning for that specific case.
%left	ShiftInstead

%left	'\n'
%left	_ASSERT
// '=' and augmented assignments have the lowest precedence
// e.g. "x = a if c > 0 else 'bar'"
// followed by
// 'if' and 'else' which have lower precedence than all other operators.
// e.g. "a, b if c > 0 else 'foo'" is either a tuple of (a,b) or 'foo'
// and not a tuple of "(a, (b if ... ))"
%left  '=' _AUGM
%left  _IF _ELSE _ELIF
%left  ','
%left  ':'
%left  _IS
%left  _OR
%left  _AND
%left  '<' '>' _EQ _NE _LE _GE _NOT _IN
%left  '+' '-'
%left  '*' '/' '%'
%left  '.' '[' '('
%right _UNARY
%left  _STRING

%%

// Grammar rules.
//
// A note on names: if foo is a rule, then foos is a sequence of foos
// (with interleaved commas or other syntax as appropriate)
// and foo_opt is an optional foo.

file:
	stmts _EOF
	{
		yylex.(*input).file = &File{Stmt: $1}
		return 0
	}

suite:
	'\n' comments _INDENT stmts _UNINDENT
	{
		statements := $4
		if $2 != nil {
			// $2 can only contain *CommentBlock objects, each of them contains a non-empty After slice
			cb := $2[len($2)-1].(*CommentBlock)
			// $4 can't be empty and can't start with a comment
			stmt := $4[0]
			start, _ := stmt.Span()
			if start.Line - cb.After[len(cb.After)-1].Start.Line == 1 {
				// The first statement of $4 starts on the next line after the last comment of $2.
				// Attach the last comment to the first statement
				stmt.Comment().Before = cb.After
				$2 = $2[:len($2)-1]
			}
			statements = append($2, $4...)
		}
		$$ = statements
	}
|	simple_stmt
	{
		$$ = $1
	}

comments:
	{
		$$ = nil
		$<lastRule>$ = nil
	}
|	comments _COMMENT '\n'
	{
		$$ = $1
		$<lastRule>$ = $<lastRule>1
		if $<lastRule>$ == nil {
			cb := &CommentBlock{Start: $2}
			$$ = append($$, cb)
			$<lastRule>$ = cb
		}
		com := $<lastRule>$.Comment()
		com.After = append(com.After, Comment{Start: $2, Token: $<tok>2})
	}
|	comments '\n'
	{
		$$ = $1
		$<lastRule>$ = nil
	}

stmts:
	{
		$$ = nil
		$<lastRule>$ = nil
	}
|	stmts stmt
	{
		// If this statement follows a comment block,
		// attach the comments to the statement.
		if cb, ok := $<lastRule>1.(*CommentBlock); ok {
			$$ = append($1[:len($1)-1], $2...)
			$2[0].Comment().Before = cb.After
			$<lastRule>$ = $2[len($2)-1]
			break
		}

		// Otherwise add to list.
		$$ = append($1, $2...)
		$<lastRule>$ = $2[len($2)-1]

		// Consider this input:
		//
		//	foo()
		//	# bar
		//	baz()
		//
		// If we've just parsed baz(), the # bar is attached to
		// foo() as an After comment. Make it a Before comment
		// for baz() instead.
		if x := $<lastRule>1; x != nil {
			com := x.Comment()
			// stmt is never empty
			$2[0].Comment().Before = com.After
			com.After = nil
		}
	}
|	stmts '\n'
	{
		// Blank line; sever last rule from future comments.
		$$ = $1
		$<lastRule>$ = nil
	}
|	stmts _COMMENT '\n'
	{
		$$ = $1
		$<lastRule>$ = $<lastRule>1
		if $<lastRule>$ == nil {
			cb := &CommentBlock{Start: $2}
			$$ = append($$, cb)
			$<lastRule>$ = cb
		}
		com := $<lastRule>$.Comment()
		com.After = append(com.After, Comment{Start: $2, Token: $<tok>2})
	}

stmt:
	simple_stmt
	{
		$$ = $1
	}
|	block_stmt
	{
		$$ = []Expr{$1}
	}

block_stmt:
	_DEF _IDENT '(' parameters_opt ')' ':' suite
	{
		$$ = &DefStmt{
			StartPos: $1,
			Name: $<tok>2,
			Params: $4,
			Body: $7,
			ForceCompact: forceCompact($3, $4, $5),
			ForceMultiLine: forceMultiLine($3, $4, $5),
		}
	}
|	_FOR loop_vars _IN expr ':' suite
	{
		$$ = &ForStmt{
			For: $1,
			Vars: $2,
			X: $4,
			Body: $6,
		}
	}
|	if_else_block
	{
		$$ = $1
	}

// A single else-statement
else_block:
	_ELSE ':' suite
	{
		$$ = &IfStmt{
			ElsePos: $1,
			False: $3,
		}
	}

// One or several elif-elif-else statements
elif_chain:
	else_block
|	elif expr ':' suite elif_chain
	{
		inner := $5
		inner.If = $1
		inner.Cond = $2
		inner.True = $4
		$$ = &IfStmt{
			ElsePos: $1,
			False: []Expr{inner},
		}
	}

// A single if-block
if_block:
	_IF expr ':' suite
	{
		$$ = &IfStmt{
			If: $1,
			Cond: $2,
			True: $4,
		}
	}

// A complete if-elif-elif-else chain
if_else_block:
	if_block
|	if_block elif_chain
	{
		$$ = $1
		$$.ElsePos = $2.ElsePos
		$$.False = $2.False
	}

elif:
	_ELSE _IF
|	_ELIF

simple_stmt:
	small_stmt small_stmts_continuation semi_opt '\n'
	{
		$$ = append([]Expr{$1}, $2...)
		$<lastRule>$ = $$[len($$)-1]
	}

small_stmts_continuation:
	{
		$$ = []Expr{}
	}
|	small_stmts_continuation ';' small_stmt
	{
		$$ = append($1, $3)
	}

small_stmt:
	expr %prec ShiftInstead
|	_RETURN expr
	{
		$$ = &ReturnStmt{
			Return: $1,
			Result: $2,
		}
	}
|	_RETURN
	{
		$$ = &ReturnStmt{
			Return: $1,
		}
	}
|	expr '=' expr      { $$ = binary($1, $2, $<tok>2, $3) }
|	expr _AUGM expr    { $$ = binary($1, $2, $<tok>2, $3) }
|	_PYTHON
	{
		$$ = &PythonBlock{Start: $1, Token: $<tok>1}
	}

semi_opt:
|	';'

primary_expr:
	ident
|	primary_expr '.' _IDENT
	{
		$$ = &DotExpr{
			X: $1,
			Dot: $2,
			NamePos: $3,
			Name: $<tok>3,
		}
	}
|	_LOAD '(' arguments_opt ')'
	{
		$$ = &CallExpr{
			X: &LiteralExpr{Start: $1, Token: "load"},
			ListStart: $2,
			List: $3,
			End: End{Pos: $4},
			ForceCompact: forceCompact($2, $3, $4),
			ForceMultiLine: forceMultiLine($2, $3, $4),
		}
	}
|	primary_expr '(' arguments_opt ')'
	{
		$$ = &CallExpr{
			X: $1,
			ListStart: $2,
			List: $3,
			End: End{Pos: $4},
			ForceCompact: forceCompact($2, $3, $4),
			ForceMultiLine: forceMultiLine($2, $3, $4),
		}
	}
|	primary_expr '[' expr ']'
	{
		$$ = &IndexExpr{
			X: $1,
			IndexStart: $2,
			Y: $3,
			End: $4,
		}
	}
|	primary_expr '[' expr_opt ':' test_opt ']'
	{
		$$ = &SliceExpr{
			X: $1,
			SliceStart: $2,
			From: $3,
			FirstColon: $4,
			To: $5,
			End: $6,
		}
	}
|	primary_expr '[' expr_opt ':' test_opt ':' test_opt ']'
	{
		$$ = &SliceExpr{
			X: $1,
			SliceStart: $2,
			From: $3,
			FirstColon: $4,
			To: $5,
			SecondColon: $6,
			Step: $7,
			End: $8,
		}
	}
|	primary_expr '(' expr for_clauses_with_if_clauses_opt ')'  // TODO: remove, not supported
	{
		$$ = &CallExpr{
			X: $1,
			ListStart: $2,
			List: []Expr{
				&ListForExpr{
					Brack: "",
					Start: $2,
					X: $3,
					For: $4,
					End: End{Pos: $5},
				},
			},
			End: End{Pos: $5},
		}
	}
|	strings %prec ShiftInstead
	{
		if len($1) == 1 {
			$$ = $1[0]
			break
		}
		$$ = $1[0]
		for _, x := range $1[1:] {
			_, end := $$.Span()
			$$ = binary($$, end, "+", x)
		}
	}
|	'[' tests_opt ']'
	{
		$$ = &ListExpr{
			Start: $1,
			List: $2,
			End: End{Pos: $3},
			ForceMultiLine: forceMultiLine($1, $2, $3),
		}
	}
|	'[' test for_clauses_with_if_clauses_opt ']'
	{
		exprStart, _ := $2.Span()
		$$ = &ListForExpr{
			Brack: "[]",
			Start: $1,
			X: $2,
			For: $3,
			End: End{Pos: $4},
			ForceMultiLine: $1.Line != exprStart.Line,
		}
	}
|	'(' test for_clauses_with_if_clauses_opt ')'
	{
		exprStart, _ := $2.Span()
		$$ = &ListForExpr{
			Brack: "()",
			Start: $1,
			X: $2,
			For: $3,
			End: End{Pos: $4},
			ForceMultiLine: $1.Line != exprStart.Line,
		}
	}
|	'{' keyvalue for_clauses_with_if_clauses_opt '}'
	{
		exprStart, _ := $2.Span()
		$$ = &ListForExpr{
			Brack: "{}",
			Start: $1,
			X: $2,
			For: $3,
			End: End{Pos: $4},
			ForceMultiLine: $1.Line != exprStart.Line,
		}
	}
|	'{' keyvalues '}'
	{
		$$ = &DictExpr{
			Start: $1,
			List: $2,
			End: End{Pos: $3},
			ForceMultiLine: forceMultiLine($1, $2, $3),
		}
	}
|	'{' tests_opt '}'  // TODO: remove, not supported
	{
		$$ = &SetExpr{
			Start: $1,
			List: $2,
			End: End{Pos: $3},
			ForceMultiLine: forceMultiLine($1, $2, $3),
		}
	}
|	'(' tests_opt ')'
	{
		if len($2) == 1 && $<comma>2.Line == 0 {
			// Just a parenthesized expression, not a tuple.
			$$ = &ParenExpr{
				Start: $1,
				X: $2[0],
				End: End{Pos: $3},
				ForceMultiLine: forceMultiLine($1, $2, $3),
			}
		} else {
			$$ = &TupleExpr{
				Start: $1,
				List: $2,
				End: End{Pos: $3},
				ForceCompact: forceCompact($1, $2, $3),
				ForceMultiLine: forceMultiLine($1, $2, $3),
			}
		}
	}

arguments_opt:
	{
		$$ = nil
	}
|	arguments comma_opt
	{
		$$ = $1
	}

arguments:
	argument
	{
		$$ = []Expr{$1}
	}
|	arguments ',' argument
	{
		$$ = append($1, $3)
	}

argument:
	test
|	ident '=' test
	{
		$$ = binary($1, $2, $<tok>2, $3)
	}
|	'*' test
	{
		$$ = unary($1, $<tok>1, $2)
	}
|	_STAR_STAR test
	{
		$$ = unary($1, $<tok>1, $2)
	}

parameters_opt:
	{
		$$ = nil
	}
|	parameters comma_opt
	{
		$$ = $1
	}

parameters:
	parameter
	{
		$$ = []Expr{$1}
	}
|	parameters ',' parameter
	{
		$$ = append($1, $3)
	}

parameter:
	ident
|	ident '=' test
	{
		$$ = binary($1, $2, $<tok>2, $3)
	}
|	'*' ident
	{
		$$ = unary($1, $<tok>1, $2)
	}
|	_STAR_STAR ident
	{
		$$ = unary($1, $<tok>1, $2)
	}

expr:
	test
|	expr ',' test
	{
		tuple, ok := $1.(*TupleExpr)
		if !ok || !tuple.Start.IsValid() {
			tuple = &TupleExpr{
				List: []Expr{$1},
				ForceCompact: true,
				ForceMultiLine: false,
			}
		}
		tuple.List = append(tuple.List, $3)
		$$ = tuple
	}

expr_opt:
	{
		$$ = nil
	}
|	expr

exprs:
	expr
	{
		$$ = []Expr{$1}
	}
|	exprs ',' expr
	{
		$$ = append($1, $3)
	}

exprs_opt:
	{
		$$ = nil
	}
|	exprs comma_opt
	{
		$$ = $1
	}

test:
	primary_expr
|	_LAMBDA exprs_opt ':' expr  // TODO: remove, not supported
	{
		$$ = &LambdaExpr{
			Lambda: $1,
			Var: $2,
			Colon: $3,
			Expr: $4,
		}
	}
|	_NOT test %prec _UNARY { $$ = unary($1, $<tok>1, $2) }
|	'-' test  %prec _UNARY { $$ = unary($1, $<tok>1, $2) }
|	test '*' test      { $$ = binary($1, $2, $<tok>2, $3) }
|	test '%' test      { $$ = binary($1, $2, $<tok>2, $3) }
|	test '/' test      { $$ = binary($1, $2, $<tok>2, $3) }
|	test '+' test      { $$ = binary($1, $2, $<tok>2, $3) }
|	test '-' test      { $$ = binary($1, $2, $<tok>2, $3) }
|	test '<' test      { $$ = binary($1, $2, $<tok>2, $3) }
|	test '>' test      { $$ = binary($1, $2, $<tok>2, $3) }
|	test _EQ test      { $$ = binary($1, $2, $<tok>2, $3) }
|	test _LE test      { $$ = binary($1, $2, $<tok>2, $3) }
|	test _NE test      { $$ = binary($1, $2, $<tok>2, $3) }
|	test _GE test      { $$ = binary($1, $2, $<tok>2, $3) }
|	test _IN test      { $$ = binary($1, $2, $<tok>2, $3) }
|	test _NOT _IN test { $$ = binary($1, $2, "not in", $4) }
|	test _OR test      { $$ = binary($1, $2, $<tok>2, $3) }
|	test _AND test     { $$ = binary($1, $2, $<tok>2, $3) }
|	test _IS test
	{
		if b, ok := $3.(*UnaryExpr); ok && b.Op == "not" {
			$$ = binary($1, $2, "is not", b.X)
		} else {
			$$ = binary($1, $2, $<tok>2, $3)
		}
	}
|	test _IF test _ELSE test
	{
		$$ = &ConditionalExpr{
			Then: $1,
			IfStart: $2,
			Test: $3,
			ElseStart: $4,
			Else: $5,
		}
	}

tests:
	test
	{
		$$ = []Expr{$1}
	}
|	tests ',' test
	{
		$$ = append($1, $3)
	}

test_opt:
	{
		$$ = nil
	}
|	test

tests_opt:
	{
		$$, $<comma>$ = nil, Position{}
	}
|	tests comma_opt
	{
		$$, $<comma>$ = $1, $2
	}

// comma_opt is an optional comma. If the comma is present,
// the rule's value is the position of the comma. Otherwise
// the rule's value is the zero position. Tracking this
// lets us distinguish (x) and (x,).
comma_opt:
	{
		$$ = Position{}
	}
|	','

keyvalue:
	test ':' test  {
		$$ = &KeyValueExpr{
			Key: $1,
			Colon: $2,
			Value: $3,
		}
	}

keyvalues_no_comma:
	keyvalue
	{
		$$ = []Expr{$1}
	}
|	keyvalues_no_comma ',' keyvalue
	{
		$$ = append($1, $3)
	}

keyvalues:
	keyvalues_no_comma
	{
		$$ = $1
	}
|	keyvalues_no_comma ','
	{
		$$ = $1
	}

loop_vars:
	primary_expr
	{
		$$ = []Expr{$1}
	}
|	loop_vars ',' primary_expr
	{
		$$ = append($1, $3)
	}

string:
	_STRING
	{
		$$ = &StringExpr{
			Start: $1,
			Value: $<str>1,
			TripleQuote: $<triple>1,
			End: $1.add($<tok>1),
			Token: $<tok>1,
		}
	}

strings:
	string
	{
		$$ = []*StringExpr{$1}
	}
|	strings string
	{
		$$ = append($1, $2)
	}

ident:
	_IDENT
	{
		$$ = &LiteralExpr{Start: $1, Token: $<tok>1}
	}

for_clause:
	_FOR loop_vars _IN test
	{
		$$ = &ForClause{
			For: $1,
			Var: $2,
			In: $3,
			Expr: $4,
		}
	}

for_clause_with_if_clauses_opt:
	for_clause if_clauses_opt {
		$$ = &ForClauseWithIfClausesOpt{
			For: $1,
			Ifs: $2,
		}
	}

for_clauses_with_if_clauses_opt:
	for_clause_with_if_clauses_opt
	{
		$$ = []*ForClauseWithIfClausesOpt{$1}
	}
|	for_clauses_with_if_clauses_opt for_clause_with_if_clauses_opt {
		$$ = append($1, $2)
	}

if_clauses_opt:
	{
		$$ = nil
	}
|	if_clauses_opt _IF test
	{
		$$ = append($1, &IfClause{
			If: $2,
			Cond: $3,
		})
	}

%%

// Go helper code.

// unary returns a unary expression with the given
// position, operator, and subexpression.
func unary(pos Position, op string, x Expr) Expr {
	return &UnaryExpr{
		OpStart: pos,
		Op:      op,
		X:       x,
	}
}

// binary returns a binary expression with the given
// operands, position, and operator.
func binary(x Expr, pos Position, op string, y Expr) Expr {
	_, xend := x.Span()
	ystart, _ := y.Span()
	return &BinaryExpr{
		X:       x,
		OpStart: pos,
		Op:      op,
		LineBreak: xend.Line < ystart.Line,
		Y:       y,
	}
}

// isSimpleExpression returns whether an expression is simple and allowed to exist in
// compact forms of sequences.
// The formal criteria are the following: an expression is considered simple if it's
// a literal (variable, string or a number), a literal with a unary operator or an empty sequence.
func isSimpleExpression(expr *Expr) bool {
	switch x := (*expr).(type) {
	case *LiteralExpr, *StringExpr:
		return true
	case *UnaryExpr:
		_, ok := x.X.(*LiteralExpr)
		return ok
	case *ListExpr:
		return len(x.List) == 0
	case *TupleExpr:
		return len(x.List) == 0
	case *DictExpr:
		return len(x.List) == 0
	case *SetExpr:
		return len(x.List) == 0
	default:
		return false
	}
}

// forceCompact returns the setting for the ForceCompact field for a call or tuple.
//
// NOTE 1: The field is called ForceCompact, not ForceSingleLine,
// because it only affects the formatting associated with the call or tuple syntax,
// not the formatting of the arguments. For example:
//
//	call([
//		1,
//		2,
//		3,
//	])
//
// is still a compact call even though it runs on multiple lines.
//
// In contrast the multiline form puts a linebreak after the (.
//
//	call(
//		[
//			1,
//			2,
//			3,
//		],
//	)
//
// NOTE 2: Because of NOTE 1, we cannot use start and end on the
// same line as a signal for compact mode: the formatting of an
// embedded list might move the end to a different line, which would
// then look different on rereading and cause buildifier not to be
// idempotent. Instead, we have to look at properties guaranteed
// to be preserved by the reformatting, namely that the opening
// paren and the first expression are on the same line and that
// each subsequent expression begins on the same line as the last
// one ended (no line breaks after comma).
func forceCompact(start Position, list []Expr, end Position) bool {
	if len(list) <= 1 {
		// The call or tuple will probably be compact anyway; don't force it.
		return false
	}

	// If there are any named arguments or non-string, non-literal
	// arguments, cannot force compact mode.
	line := start.Line
	for _, x := range list {
		start, end := x.Span()
		if start.Line != line {
			return false
		}
		line = end.Line
		if !isSimpleExpression(&x) {
			return false
		}
	}
	return end.Line == line
}

// forceMultiLine returns the setting for the ForceMultiLine field.
func forceMultiLine(start Position, list []Expr, end Position) bool {
	if len(list) > 1 {
		// The call will be multiline anyway, because it has multiple elements. Don't force it.
		return false
	}

	if len(list) == 0 {
		// Empty list: use position of brackets.
		return start.Line != end.Line
	}

	// Single-element list.
	// Check whether opening bracket is on different line than beginning of
	// element, or closing bracket is on different line than end of element.
	elemStart, elemEnd := list[0].Span()
	return start.Line != elemStart.Line || end.Line != elemEnd.Line
}
