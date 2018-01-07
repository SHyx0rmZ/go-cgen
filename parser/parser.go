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
	var m = map[lexer.ItemType]func() ast.Node{
		lexer.ItemEndIf:        func() ast.Node { return &ast.EndIfDir{DirPos: p.next().Pos} },
		lexer.ItemElseDir:      func() ast.Node { return &ast.ElseDir{DirPos: p.next().Pos} },
		lexer.ItemDefine:       func() ast.Node { return p.parseMacroDir() },
		lexer.ItemInclude:      func() ast.Node { return p.parseIncludeDir() },
		lexer.ItemIfDefined:    func() ast.Node { return p.parseIfDefDir(ast.DEFINED) },
		lexer.ItemIfNotDefined: func() ast.Node { return p.parseIfDefDir(ast.NOT_DEFINED) },
		lexer.ItemComment: func() ast.Node {
			comment := p.next()
			return &ast.Comment{
				Slash: comment.Pos,
				Text:  comment.Val,
			}
		},
		lexer.ItemExtern: func() ast.Node { return p.parseExternDecl() },
	}
	c := make(chan ast.Node)
	go func() {
		defer close(c)
		for {
			i := p.peek()
			f, ok := m[i.Typ]
			if ok {
				c <- f()
				continue
			}
			switch i.Typ {
			case lexer.ItemEOF:
				return
			case lexer.ItemError:
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
	switch tok.Typ {
	case lexer.ItemMinus:
		return tok, token.SUB, 4
	case lexer.ItemSlash:
		return tok, token.QUO, 5
	}
	p.backup()
	return tok, token.ILLEGAL_TOKEN, 0
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

func (p *parser) nextNonSpace() (token lexer.Item) {
	for {
		token = p.next()
		if token.Typ != lexer.ItemSpace {
			break
		}
	}
	return token
}

func (p *parser) peekNonSpace() (token lexer.Item) {
	for {
		token = p.next()
		if token.Typ != lexer.ItemSpace {
			break
		}
	}
	p.backup()
	return token
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

func (p *parser) expect(expected lexer.ItemType, context string) lexer.Item {
	token := p.nextNonSpace()
	if token.Typ != expected {
		p.unexpected(token, context)
	}
	return token
}

func (p *parser) expectOneOf(expected1, expected2 lexer.ItemType, context string) lexer.Item {
	token := p.nextNonSpace()
	if token.Typ != expected1 && token.Typ != expected2 {
		p.unexpected(token, context)
	}
	return token
}

func (p *parser) unexpected(token lexer.Item, context string) {
	p.errorf("unexpected %s in %s", token, context)
}
