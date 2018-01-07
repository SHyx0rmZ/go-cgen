package token

import "fmt"

type Token int

const (
	ILLEGAL Token = iota
	EOF
	COMMENT
	WHITESPACE

	INCLUDE_PATH

	literal_beg
	IDENT
	INT
	STRING
	literal_end

	operator_beg
	// Operators and delimiters
	ADD // +
	SUB // -
	MUL // *
	QUO // /
	REM // %

	AND // &
	OR  // |
	XOR // ^
	SHL // <<
	SHR // >>

	LAND // &&
	LOR  // ||
	INC  // ++
	DEC  // --

	ASSIGN // =

	LPAREN // (
	LBRACE // {
	COMMA  // ,

	RPAREN    // )
	RBRACE    // }
	SEMICOLON // ;
	operator_end

	keyword_beg
	DEFINE  // #define
	ELSE    // #else
	ENDIF   // #endif
	IFDEF   // #ifdef
	IFNDEF  // #ifndef
	INCLUDE // #include
	EXTERN  // extern
	keyword_end
)

var tokens = [...]string{
	ILLEGAL:    "ILLEGAL",
	EOF:        "EOF",
	COMMENT:    "COMMENT",
	WHITESPACE: "WHITESPACE",

	INCLUDE_PATH: "INCLUDE_PATH",

	IDENT:  "IDENT",
	INT:    "INT",
	STRING: "STRING",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	QUO: "/",
	REM: "%",

	AND: "&",
	OR:  "|",
	XOR: "^",
	SHL: "<<",
	SHR: ">>",

	LAND: "&&",
	LOR:  "||",
	INC:  "++",
	DEC:  "--",

	ASSIGN: "=",

	LPAREN: "(",
	LBRACE: "{",
	COMMA:  ",",

	RPAREN:    ")",
	RBRACE:    "}",
	SEMICOLON: ";",

	DEFINE:  "#define",
	ELSE:    "#else",
	ENDIF:   "#endif",
	IFDEF:   "#ifdef",
	IFNDEF:  "#ifndef",
	INCLUDE: "#include",
	EXTERN:  "extern",
}

func (t Token) String() string {
	s := ""
	if 0 <= t && t < Token(len(tokens)) {
		s = tokens[t]
	}
	if s == "" {
		s = fmt.Sprintf("token(%d)", t)
	}
	return s
}

func (t Token) Precedence() int {
	switch t {
	case LOR:
		return 1
	case LAND:
		return 2
	case ADD, SUB, OR, XOR:
		return 4
	case MUL, QUO, REM, SHL, SHR, AND:
		return 5
	}
	return 0
}
