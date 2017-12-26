package cgen

type parser struct {
	lex *lexer
}

func NewParser(name, input string) *parser {
	return &parser{lex: NewLexer(name, input)}
}

func (p *parser) Parse() chan item {
	return p.lex.items
}
