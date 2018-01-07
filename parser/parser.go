package parser

import (
	"fmt"
	"strings"

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
			if !ok {
				switch i.Typ {
				case lexer.ItemEOF:
					return
				case lexer.ItemError:
					p.errorf(i.Val)
				case lexer.ItemComment:
					c <- &ast.Comment{
						Slash: i.Pos,
						Text:  i.Val,
					}
				case lexer.ItemIdentifier:
					if i.Val == "typedef" {
						for {
							i2 := p.next()
							if i2.Typ == lexer.ItemSpace && strings.Contains(i2.Val, "\n") {
								c <- &ast.BadExpr{
									From: i.Pos,
									To:   i2.Pos,
								}
								break
							}
						}
						continue
						i1 := p.peekNonSpace()
						switch {
						case i1.Typ == lexer.ItemIdentifier && i1.Val == "struct":
							p.nextNonSpace()
							switch p.peekNonSpace().Typ {
							/*
								case lexer.ItemOpenCurly:
									var ns []Expr
									for {
										ic := p.nextNonSpace()
										if ic.typ == lexer.ItemCloseCurly {
											break
										}
									}
									i2 := p.expect(lexer.ItemIdentifier, "typedef")
									p.expect(lexer.ItemSemicolon, "typedef")
									c <- &TypeDecl{
										Name: Ident{
											NamePos: i2.Pos,
											Name:    i2.Val,
										},
										Expr: StructDecl{
											Nodes: ns,
										},
									}
							*/
							/*
								case lexer.ItemIdentifier:
									i2 := p.nextNonSpace()
									i3 := p.expect(lexer.ItemIdentifier, "Typedef")
									p.expect(lexer.ItemSemicolon, "typedef")
									c <- TypeDecl{
										Name: Ident{
											NamePos: i3.Pos,
											Name:    i3.Val,
										},
										Expr: StructType{
											Name: Ident{
												NamePos: i2.Pos,
												Name:    i2.Val,
											},
										},
									}
							*/
							default:
								p.unexpected(p.peekNonSpace(), "typedef")
							}
							continue
							/*
								case i1.Typ == lexer.ItemIdentifier && i1.Val == "enum":
									p.nextNonSpace()
									p.expect(lexer.ItemOpenCurly, "typedef")
									var ds []EnumSpec
									for {
										ic := p.peekNonSpace()
										if ic.Typ == lexer.ItemCloseCurly {
											break
										}
										if ic.Typ == lexer.ItemComment {
											p.nextNonSpace()
											continue
										}
										ic = p.expect(lexer.ItemIdentifier, "enum")
										ic1 := p.peekNonSpace()
										if ic1.Typ != lexer.ItemComma && ic1.Typ != lexer.ItemCloseCurly && ic1.Typ != lexer.ItemComment {
											p.expect(lexer.ItemEqualSign, "enum")
											switch p.peekNonSpace().Typ {
											case lexer.ItemHexValue:
												ic3 := p.expect(lexer.ItemHexValue, "enum")
												if p.peekNonSpace().Typ == lexer.ItemComment {
													p.nextNonSpace()
												}
												ds = append(ds, EnumValue{
													Name: Ident{
														NamePos: ic.Pos,
														Name:    ic.Val,
													},
													Value: &BasicLit{
														ValuePos: ic3.Pos,
														Kind:     Number,
														Value:    ic3.Val,
													},
												})
											case lexer.ItemOpenParen:
												p.expect(lexer.ItemOpenParen, "enum hack")
												ic3 := p.expect(lexer.ItemIdentifier, "enum hack")
												ic4 := p.expect(lexer.ItemBitwiseOr, "enum hack")
												ic5 := p.expect(lexer.ItemHexValue, "enum hack")
												p.expect(lexer.ItemCloseParen, "enum hack")
												ds = append(ds, EnumConstExpr{
													Name: Ident{
														NamePos: ic.Pos,
														Name:    ic.Val,
													},
													Expr: BinaryExpr{
														X: Ident{
															NamePos: ic3.Pos,
															Name:    ic3.Val,
														},
														OpPos: ic4.Pos,
														Op:    BitwiseOrOp,
														Y: BasicLit{
															ValuePos: ic5.Pos,
															Kind:     Number,
															Value:    ic5.Val,
														},
													},
												})
											default:
												p.unexpected(p.peekNonSpace(), "enum")
											}
											//if p.peekNonSpace().Typ != lexer.ItemComma {
											//	break
											//}
											//p.nextNonSpace()
										} else {
											//if p.peekNonSpace().Typ == lexer.ItemComma {
											//	p.nextNonSpace()
											//}
											ds = append(ds, EnumValue{
												Name: &Ident{
													NamePos: ic.Pos,
													Name:    ic.Val,
												},
												Value: nil,
											})
										}
										if p.peekNonSpace().Typ == lexer.ItemComment {
											ic := p.nextNonSpace()
											c <- &Comment{
												Slash: ic.Pos,
												Text:  ic.Val,
											}
										}
										if p.peekNonSpace().Typ != lexer.ItemComma {
											break
										}
										p.nextNonSpace()
									}
									p.expect(lexer.ItemCloseCurly, "typedef")
									i2 := p.expect(lexer.ItemIdentifier, "typedef")
									p.expect(lexer.ItemSemicolon, "typedef")
									c <- TypeDecl{
										Name: Ident{
											NamePos: i2.Pos,
											Name:    i2.Val,
										},
										Expr: EnumDecl{
											Specs: ds,
										},
									}
									continue
							*/
						}
					}
					//fallthrough
					//case lexer.ItemOpenParen:
					//	p.backup()
					//	c <- p.parseParenExpr()
				case lexer.ItemSpace:
					if strings.Contains(i.Val, "\n") {
						continue
					}
					fallthrough
				default:
					c <- p.parseExpr()
				}
				continue
			}
			c <- f()
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

type Error struct {
}

func (Error) String() string { return "Error" }

type Token int

const (
	Number Token = iota
	BitwiseOrOp
)
