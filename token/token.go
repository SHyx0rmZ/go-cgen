package token

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

/*

	ItemHexValue
	ItemExternC
	ItemOpenCurly
	ItemCloseCurly
	ItemOpenParen
	ItemCloseParen
	ItemIdentifier
	ItemEndIf
	ItemLogicalOr
	ItemBitwiseOr
	ItemExtern
	ItemSemicolon
	ItemStar
	ItemComma
	ItemEqualSign
	ItemLogicalAnd
	ItemBitwiseAnd
	ItemSpace
	ItemIncludePath
	ItemIncludePathSystem
	ItemDecrement
	ItemMinus
	ItemSlash
	ItemString
*/
