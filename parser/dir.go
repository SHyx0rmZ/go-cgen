package parser

import (
	"github.com/SHyx0rmZ/cgen/ast"
	"github.com/SHyx0rmZ/cgen/token"
	"strings"
)

func (p *parser) parseArgList() *ast.ArgList {
	if p.peek().Tok != token.LPAREN {
		return nil
	}

	open := p.next()
	var list []*ast.Ident
	if p.peek().Tok == token.IDENT {
		id := p.next()
		list = append(list, &ast.Ident{
			NamePos: id.Pos,
			Name:    id.Val,
		})

		for p.peek().Tok == token.COMMA {
			id = p.expect(token.IDENT, "macro argument list")
			list = append(list, &ast.Ident{
				NamePos: id.Pos,
				Name:    id.Val,
			})
		}
	}
	closing := p.expect(token.RPAREN, "macro argument list")

	return &ast.ArgList{
		Opening: open.Pos,
		List:    list,
		Closing: closing.Pos,
	}
}

func (p *parser) parseMacroDir() ast.Dir {
	keyword := p.expect(token.DEFINE, "macro definition")
	name := p.expect(token.IDENT, "macro definition")
	args := p.parseArgList()
	switch p.peek().Tok {
	case token.WHITESPACE, token.EOF:
		if p.peek().Val == "" || strings.Contains(p.peek().Val, "\n") {
			return &ast.MacroDir{
				DirPos: keyword.Pos,
				Name: &ast.Ident{
					NamePos: name.Pos,
					Name:    name.Val,
				},
				Args:  args,
				Value: nil,
			}
		}
		fallthrough
	default:
		p.next()
		return &ast.MacroDir{
			DirPos: keyword.Pos,
			Name: &ast.Ident{
				NamePos: name.Pos,
				Name:    name.Val,
			},
			Args:  args,
			Value: p.parseExpr(),
		}
	}
}

func (p *parser) parseIncludeDir() ast.Dir {
	keyword := p.expect(token.INCLUDE, "include directive")
	path := p.expect(token.INCLUDE_PATH, "include directive")
	return &ast.IncludeDir{
		DirPos:  keyword.Pos,
		PathPos: path.Pos,
		Path:    path.Val,
	}
}

func (p *parser) parseIfDefDir(cond ast.IfDefCond) ast.Dir {
	keyword := p.expectOneOf(token.IFDEF, token.IFNDEF, "conditional directive")
	identifier := p.expect(token.IDENT, "conditional directive")
	return &ast.IfDefDir{
		DirPos: keyword.Pos,
		Cond:   cond,
		Name: &ast.Ident{
			NamePos: identifier.Pos,
			Name:    identifier.Val,
		},
	}
}
