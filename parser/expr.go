package parser

import (
	"fmt"

	"github.com/SHyx0rmZ/cgen/ast"
	"github.com/SHyx0rmZ/cgen/token"
)

func (p *parser) parsePrimaryExpr() ast.Expr {
	x := p.parseOperand()
	return x
}

func (p *parser) parseUnaryExpr() ast.Expr {
	switch p.peekNonSpace().Tok {
	case token.INT:
		number := p.next()
		return &ast.BasicLit{
			ValuePos: number.Pos,
			Kind:     token.INT,
			Value:    number.Val,
		}
	case token.IDENT:
		identifier := p.next()
		return &ast.Ident{
			NamePos: identifier.Pos,
			Name:    identifier.Val,
		}
	//case token.ADD, token.SUB, token.AND, token.OR, token.XOR:
	case token.SUB:
		operator := p.next()
		expr := p.parseUnaryExpr()
		return &ast.UnaryExpr{
			OpPos: operator.Pos,
			Op:    operator.Tok,
			X:     expr,
		}
	case token.LPAREN:
		//opening := p.next()
		//expr := p.parseExpr()
		//fmt.Printf("%#v\n", expr)
		//closing := p.expect(token.RPAREN, "parentheses expression")
		//return &ast.ParenExpr{
		//	Opening: opening.Pos,
		//	Expr:    expr,
		//	Closing: closing.Pos,
		//}
		return p.parseParenExpr()
	}

	fmt.Printf("B %s\n", p.peekNonSpace())

	return p.parsePrimaryExpr()
}

func (p *parser) parseBinaryExpr(prec1 int) ast.Expr {
	x := p.parseUnaryExpr()
	for {
		op := p.nextNonSpace()
		oprec := op.Tok.Precedence()
		p.backup()
		if oprec < prec1 {
			return x
		}
		p.next()
		y := p.parseBinaryExpr(oprec + 1)
		x = &ast.BinaryExpr{
			X:     x,
			OpPos: op.Pos,
			Op:    op.Tok,
			Y:     y,
		}
	}
}

func (p *parser) parseExpr() ast.Expr {
	return p.parseBinaryExpr(1)
}

func (p *parser) parseParenExpr() ast.Expr {
	opening := p.expect(token.LPAREN, "parentheses expression")
	var expr ast.Expr
	if p.peek().Tok != token.RPAREN {
		expr = p.parseExpr()
	}
	closing := p.expect(token.RPAREN, "parentheses expression")
	return &ast.ParenExpr{
		Opening: opening.Pos,
		Expr:    expr,
		Closing: closing.Pos,
	}
}
