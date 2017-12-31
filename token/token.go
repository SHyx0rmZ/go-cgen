package token

type Token int

const (
	ILLEGAL_TOKEN Token = iota
	EOF

	literal_beg
	INT
	literal_end
)
