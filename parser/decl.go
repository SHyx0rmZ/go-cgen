package parser

import (
	"github.com/SHyx0rmZ/cgen/ast"
	"github.com/SHyx0rmZ/cgen/lexer"
	"github.com/SHyx0rmZ/cgen/token"
)

func (p *parser) parseExternDecl() ast.Decl {
	keyword := p.expect(lexer.ItemExtern, "external declaration")
	next := p.peekNonSpace()
	if next.Typ == lexer.ItemString && next.Val == `"C"` {
		p.next()
		curly := p.expect(lexer.ItemOpenCurly, "external declaration")
		return &ast.ExternDecl{
			KeyPos: keyword.Pos,
			Decl: &ast.CDecl{
				Value: &ast.BasicLit{
					ValuePos: next.Pos,
					Kind:     token.STRING,
					Value:    next.Val,
				},
				BodyPos: curly.Pos,
			},
		}
	}
	return &ast.ExternDecl{
		KeyPos: keyword.Pos,
		Decl:   nil,
	}
}
