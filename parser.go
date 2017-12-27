package cgen

import "fmt"

type parser struct {
	lex       *lexer
	token     [3]item
	peekCount int
	name      string
}

func NewParser(name, input string) *parser {
	return &parser{lex: NewLexer(name, input), name: name}
}

func (p *parser) Parse() chan Node {
	c := make(chan Node)
	go func() {
		for {
			i := p.next()
			if i.typ == itemEOF || i.typ == itemError {
				break
			}
			fmt.Printf("%s\n", p.peek())
			if i.typ == itemDefine {
				if p.peek().typ != itemIdentifier {
					p.errorf("expected identifier but found %s", p.peek().typ)
				}
				i1 := p.next()
				switch p.peek().typ {
				case itemHexValue:
					i2 := p.next()
					c <- DefineStmt{
						Name: Ident{
							NamePos: i1.pos,
							Name:    i1.val,
						},
						Value: BasicLit{
							ValuePos: i2.pos,
							Kind:     Number,
							Value:    i2.val,
						},
					}
				default:
					//p.errorf("unexpected %s", p.peek().typ)
				}
			}
			if i.typ == itemInclude {
				c <- IncludeStmt{
					PathPos: i.pos,
					Path:    i.val,
				}
			}
			if i.typ == itemComment {
				c <- Comment{
					Slash: i.pos,
					Text:  i.val,
				}
			}
		}
		close(c)
	}()
	return c
}

func (p *parser) errorf(format string, args ...interface{}) {
	format = fmt.Sprintf("cgen: %s:%d: %s", p.name, p.token[0].line, format)
	panic(fmt.Errorf(format, args...))
}

func (p *parser) next() item {
	if p.peekCount > 0 {
		p.peekCount--
	} else {
		p.token[0] = p.lex.nextItem()
	}
	return p.token[p.peekCount]
}

func (p *parser) backup() {
	p.peekCount++
}

func (p *parser) peek() item {
	if p.peekCount > 0 {
		return p.token[p.peekCount-1]
	}
	p.peekCount = 1
	p.token[0] = p.lex.nextItem()
	return p.token[0]
}

type Error struct {
}

func (Error) String() string { return "Error" }

type Node interface {
}

type Expr interface {
	Node
	exprNode()
}

type Stmt interface {
	Node
	stmtNode()
}

type Token int

const (
	Number Token = iota
)

type BasicLit struct {
	ValuePos Pos
	Kind     Token
	Value    string
}

func (BasicLit) exprNode() {}

type Ident struct {
	NamePos Pos
	Name    string
}

type DefineStmt struct {
	Name  Ident
	Value Expr
}

type IncludeStmt struct {
	PathPos Pos
	Path    string
}

type Comment struct {
	Slash Pos
	Text  string
}
