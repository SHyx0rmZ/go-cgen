package token

type Token int

const (
	ILLEGAL_TOKEN Token = iota
	EOF

	literal_beg
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
	operator_end
)
