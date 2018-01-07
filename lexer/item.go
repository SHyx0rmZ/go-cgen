package lexer

import (
	"fmt"
	"github.com/SHyx0rmZ/cgen/token"
)

type Item struct {
	Pos  token.Pos
	Val  string
	Tok  token.Token
	Line int
}

func (i Item) String() string {
	switch {
	case i.Tok == token.EOF:
		return "EOF"
	case i.Tok == token.ILLEGAL:
		return i.Val
	case len(i.Val) > 10:
		return fmt.Sprintf("%s(%.30q...)", i.Tok, i.Val)
	}
	return fmt.Sprintf("%s(%q)", i.Tok, i.Val)
}
