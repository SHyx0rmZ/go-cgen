package parser

import (
	"github.com/SHyx0rmZ/cgen/ast"
	"github.com/SHyx0rmZ/cgen/lexer"
	"strings"
)

func (p *parser) parseArgList() *ast.ArgList {
	if p.peek().Typ != lexer.ItemOpenParen {
		return nil
	}

	open := p.next()
	var list []*ast.Ident
	if p.peek().Typ == lexer.ItemIdentifier {
		id := p.next()
		list = append(list, &ast.Ident{
			NamePos: id.Pos,
			Name:    id.Val,
		})

		for p.peek().Typ == lexer.ItemComma {
			id = p.expect(lexer.ItemIdentifier, "macro argument list")
			list = append(list, &ast.Ident{
				NamePos: id.Pos,
				Name:    id.Val,
			})
		}
	}
	closing := p.expect(lexer.ItemCloseParen, "macro argument list")

	return &ast.ArgList{
		Opening: open.Pos,
		List:    list,
		Closing: closing.Pos,
	}
}

func (p *parser) parseMacroDir() ast.Dir {
	keyword := p.expect(lexer.ItemDefine, "macro definition")
	name := p.expect(lexer.ItemIdentifier, "macro definition")
	args := p.parseArgList()
	switch p.peek().Typ {
	case lexer.ItemSpace, lexer.ItemEOF:
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
	keyword := p.expect(lexer.ItemInclude, "include directive")
	path := p.expectOneOf(lexer.ItemIncludePath, lexer.ItemIncludePathSystem, "include directive")
	return &ast.IncludeDir{
		DirPos:  keyword.Pos,
		PathPos: path.Pos,
		Path:    path.Val,
	}
}

func (p *parser) parseIfDefDir(cond ast.IfDefCond) ast.Dir {
	keyword := p.expectOneOf(lexer.ItemIfDefined, lexer.ItemIfNotDefined, "conditional directive")
	identifier := p.expect(lexer.ItemIdentifier, "conditional directive")
	return &ast.IfDefDir{
		DirPos: keyword.Pos,
		Cond:   cond,
		Name: &ast.Ident{
			NamePos: identifier.Pos,
			Name:    identifier.Val,
		},
	}
}
