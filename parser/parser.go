package parser

import (
	"fmt"
	//"strings"

	"github.com/SHyx0rmZ/cgen/ast"
	"github.com/SHyx0rmZ/cgen/lexer"
	"github.com/SHyx0rmZ/cgen/token"
)

type parser struct {
	lex       lexer.Lexer
	token     [3]lexer.Item
	peekCount int
	name      string
	indent    int
	trace     bool

	pos   token.Pos
	tok   token.Token
	lit   string
	error chan error
}

func NewParser(name, input string) *parser {
	return &parser{lex: lexer.NewLexer(name, input), name: name, trace: true, error: make(chan error, 1)}
}

func (p *parser) Err() error {
	select {
	case err := <-p.error:
		return err
	default:
		return nil
	}
}

func (p *parser) Nodes() []ast.Node {
	var nodes []ast.Node
	for node := range p.Parse() {
		nodes = append(nodes, node)
	}
	return nodes
}

func (p *parser) Parse() chan ast.Node {
	var m = map[token.Token]func() ast.Node{
		token.ENDIF:   func() ast.Node { return &ast.EndIfDir{DirPos: p.next().Pos} },
		token.ELSE:    func() ast.Node { return &ast.ElseDir{DirPos: p.next().Pos} },
		token.DEFINE:  func() ast.Node { return p.parseMacroDir() },
		token.INCLUDE: func() ast.Node { return p.parseIncludeDir() },
		token.IFDEF:   func() ast.Node { return p.parseIfDefDir(ast.DEFINED) },
		token.IFNDEF:  func() ast.Node { return p.parseIfDefDir(ast.NOT_DEFINED) },
		token.COMMENT: func() ast.Node {
			comment := p.next()
			return &ast.Comment{
				Slash: comment.Pos,
				Text:  comment.Val,
			}
		},
		token.EXTERN: func() ast.Node { return p.parseExternDecl() },
	}
	c := make(chan ast.Node)
	go func() {
		defer close(c)
		for {
			i := p.peek()
			f, ok := m[i.Tok]
			if ok {
				c <- f()
				continue
			}
			switch i.Tok {
			case token.EOF:
				return
			case token.ILLEGAL:
				p.errorf(i.Val)
			//case lexer.ItemIdentifier:
			//if i.Val == "typedef" {
			//case lexer.ItemSpace:
			//	if strings.Contains(i.Val, "\n") {
			//		continue
			//	}
			//	fallthrough
			default:
				c <- p.parseExpr()
			}
		}
	}()
	return c
}

func (p *parser) printTrace(a ...interface{}) {
	const dots = ". . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . "
	const n = len(dots)
	fmt.Printf("%5d:???: ", p.pos)
	i := 2 * p.indent
	for i > n {
		fmt.Print(dots)
		i -= n
	}
	fmt.Print(dots[0:i])
	fmt.Println(a...)
}

func trace(p *parser, msg string) *parser {
	if p.trace {
		p.printTrace(msg, "(")
		p.indent++
	}
	return p
}

func un(p *parser) {
	if p.trace {
		p.indent--
		p.printTrace(")")
	}
}

func (p *parser) parseOperand() ast.Expr {
	/*switch p.tok {
	case token.INT:
		x := &BasicLit{
			ValuePos: p.Pos,
			Kind:     p.tok,
			Value:    p.lit,
		}
		p.next()
		return x
	}*/
	return &ast.BadExpr{
		From: p.pos,
		To:   p.pos,
	}
}

func (p *parser) tokPrec() (lexer.Item, token.Token, int) {
	tok := p.nextNonSpace()
	switch tok.Tok {
	case token.SUB:
		return tok, token.SUB, 4
	case token.QUO:
		return tok, token.QUO, 5
	}
	p.backup()
	return tok, token.ILLEGAL, 0
}

func (p *parser) next() lexer.Item {
	if p.peekCount > 0 {
		p.peekCount--
	} else {
		p.token[0] = p.lex.NextItem()
	}
	p.pos, p.tok, p.lit = p.token[p.peekCount].Pos, token.INT, p.token[p.peekCount].Val
	return p.token[p.peekCount]
}

func (p *parser) backup() {
	p.peekCount = 1
}

func (p *parser) peek() lexer.Item {
	if p.peekCount > 0 {
		return p.token[p.peekCount-1]
	}
	p.peekCount = 1
	p.token[0] = p.lex.NextItem()
	return p.token[0]
}

func (p *parser) nextNonSpace() lexer.Item {
	var t lexer.Item
	for {
		t = p.next()
		if t.Tok != token.WHITESPACE {
			break
		}
	}
	return t
}

func (p *parser) peekNonSpace() (t lexer.Item) {
	for {
		t = p.next()
		if t.Tok != token.WHITESPACE {
			break
		}
	}
	p.backup()
	return t
}

func (p *parser) errorf(format string, args ...interface{}) {
	format = fmt.Sprintf("cgen: %s:%d: %s", p.name, p.token[0].Line, format)
	go func() {
		for {
			p.error <- fmt.Errorf(format, args...)
		}
	}()
	panic(fmt.Errorf(format, args...))
}

func (p *parser) expect(expected token.Token, context string) lexer.Item {
	token := p.nextNonSpace()
	if token.Tok != expected {
		p.unexpected(token, context)
	}
	return token
}

func (p *parser) expectOneOf(expected1, expected2 token.Token, context string) lexer.Item {
	token := p.nextNonSpace()
	if token.Tok != expected1 && token.Tok != expected2 {
		p.unexpected(token, context)
	}
	return token
}

func (p *parser) unexpected(token lexer.Item, context string) {
	p.errorf("unexpected %s in %s", token, context)
}
