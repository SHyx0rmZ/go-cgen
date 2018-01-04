package cgen

import (
	"fmt"
	"strings"

	"github.com/SHyx0rmZ/cgen/token"
)

type parser struct {
	lex       *lexer
	token     [3]item
	peekCount int
	name      string
	indent    int
	trace     bool

	pos token.Pos
	tok token.Token
	lit string
}

func NewParser(name, input string) *parser {
	return &parser{lex: NewLexer(name, input), name: name, trace: true}
}

func (p *parser) Nodes() []Node {
	var nodes []Node
	for node := range p.Parse() {
		nodes = append(nodes, node)
	}
	return nodes
}

func (p *parser) Parse() chan Node {
	c := make(chan Node)
	go func() {
		defer close(c)
		for {
			i := p.next()
			switch i.typ {
			case itemEOF, itemError:
				return
			case itemEndIf:
				c <- &EndIfDir{
					DirPos: i.pos,
				}
			case itemDefine:
				p.backup()
				c <- p.parseMacroDir()
				//panic("sad")
			case itemInclude:
				if p.peekNonSpace().typ != itemIncludePath && p.peekNonSpace().typ != itemIncludePathSystem {
					p.errorf("expected include path but found %s", p.peekNonSpace().typ)
				}
				i1 := p.nextNonSpace()
				c <- &IncludeDir{
					DirPos:  i.pos,
					PathPos: i1.pos,
					Path:    i1.val,
				}
			case itemIfDefined:
				if p.peekNonSpace().typ != itemIdentifier {
					p.errorf("expected identifier but found %s", p.peekNonSpace().typ)
				}
				i1 := p.nextNonSpace()
				c <- &IfDefDir{
					DirPos: i.pos,
					Cond:   DEFINED,
					Name: &Ident{
						NamePos: i1.pos,
						Name:    i1.val,
					},
				}
			case itemIfNotDefined:
				if p.peekNonSpace().typ != itemIdentifier {
					p.errorf("expected identifier but found %s", p.peekNonSpace().typ)
				}
				i1 := p.nextNonSpace()
				c <- &IfDefDir{
					DirPos: i.pos,
					Cond:   NOT_DEFINED,
					Name: &Ident{
						NamePos: i1.pos,
						Name:    i1.val,
					},
				}
			case itemComment:
				c <- &Comment{
					Slash: i.pos,
					Text:  i.val,
				}
			case itemIdentifier:
				if i.val == "typedef" {
					i1 := p.peekNonSpace()
					switch {
					case i1.typ == itemIdentifier && i1.val == "struct":
						p.nextNonSpace()
						switch p.peekNonSpace().typ {
						/*
							case itemOpenCurly:
								var ns []Expr
								for {
									ic := p.nextNonSpace()
									if ic.typ == itemCloseCurly {
										break
									}
								}
								i2 := p.expect(itemIdentifier, "typedef")
								p.expect(itemSemicolon, "typedef")
								c <- &TypeDecl{
									Name: Ident{
										NamePos: i2.pos,
										Name:    i2.val,
									},
									Expr: StructDecl{
										Nodes: ns,
									},
								}
						*/
						/*
							case itemIdentifier:
								i2 := p.nextNonSpace()
								i3 := p.expect(itemIdentifier, "typedef")
								p.expect(itemSemicolon, "typedef")
								c <- TypeDecl{
									Name: Ident{
										NamePos: i3.pos,
										Name:    i3.val,
									},
									Expr: StructType{
										Name: Ident{
											NamePos: i2.pos,
											Name:    i2.val,
										},
									},
								}
						*/
						default:
							p.unexpected(p.peekNonSpace(), "typedef")
						}
						continue
						/*
							case i1.typ == itemIdentifier && i1.val == "enum":
								p.nextNonSpace()
								p.expect(itemOpenCurly, "typedef")
								var ds []EnumSpec
								for {
									ic := p.peekNonSpace()
									if ic.typ == itemCloseCurly {
										break
									}
									if ic.typ == itemComment {
										p.nextNonSpace()
										continue
									}
									ic = p.expect(itemIdentifier, "enum")
									ic1 := p.peekNonSpace()
									if ic1.typ != itemComma && ic1.typ != itemCloseCurly && ic1.typ != itemComment {
										p.expect(itemEqualSign, "enum")
										switch p.peekNonSpace().typ {
										case itemHexValue:
											ic3 := p.expect(itemHexValue, "enum")
											if p.peekNonSpace().typ == itemComment {
												p.nextNonSpace()
											}
											ds = append(ds, EnumValue{
												Name: Ident{
													NamePos: ic.pos,
													Name:    ic.val,
												},
												Value: &BasicLit{
													ValuePos: ic3.pos,
													Kind:     Number,
													Value:    ic3.val,
												},
											})
										case itemOpenParen:
											p.expect(itemOpenParen, "enum hack")
											ic3 := p.expect(itemIdentifier, "enum hack")
											ic4 := p.expect(itemBitwiseOr, "enum hack")
											ic5 := p.expect(itemHexValue, "enum hack")
											p.expect(itemCloseParen, "enum hack")
											ds = append(ds, EnumConstExpr{
												Name: Ident{
													NamePos: ic.pos,
													Name:    ic.val,
												},
												Expr: BinaryExpr{
													X: Ident{
														NamePos: ic3.pos,
														Name:    ic3.val,
													},
													OpPos: ic4.pos,
													Op:    BitwiseOrOp,
													Y: BasicLit{
														ValuePos: ic5.pos,
														Kind:     Number,
														Value:    ic5.val,
													},
												},
											})
										default:
											p.unexpected(p.peekNonSpace(), "enum")
										}
										//if p.peekNonSpace().typ != itemComma {
										//	break
										//}
										//p.nextNonSpace()
									} else {
										//if p.peekNonSpace().typ == itemComma {
										//	p.nextNonSpace()
										//}
										ds = append(ds, EnumValue{
											Name: &Ident{
												NamePos: ic.pos,
												Name:    ic.val,
											},
											Value: nil,
										})
									}
									if p.peekNonSpace().typ == itemComment {
										ic := p.nextNonSpace()
										c <- &Comment{
											Slash: ic.pos,
											Text:  ic.val,
										}
									}
									if p.peekNonSpace().typ != itemComma {
										break
									}
									p.nextNonSpace()
								}
								p.expect(itemCloseCurly, "typedef")
								i2 := p.expect(itemIdentifier, "typedef")
								p.expect(itemSemicolon, "typedef")
								c <- TypeDecl{
									Name: Ident{
										NamePos: i2.pos,
										Name:    i2.val,
									},
									Expr: EnumDecl{
										Specs: ds,
									},
								}
								continue
						*/
					}
				}
				fallthrough
			case itemHexValue:
				c <- &BasicLit{
					ValuePos: i.pos,
					Kind:     token.INT,
					Value:    i.val,
				}
			case itemMinus:
				p.backup()
				c <- p.parseUnaryExpr()
			//case itemOpenParen:
			//	p.backup()
			//	c <- p.parseParenExpr()
			default:
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

func (p *parser) parseOperand() Expr {
	/*switch p.tok {
	case token.INT:
		x := &BasicLit{
			ValuePos: p.pos,
			Kind:     p.tok,
			Value:    p.lit,
		}
		p.next()
		return x
	}*/
	return &BadExpr{
		From: p.pos,
		To:   p.pos,
	}
}

func (p *parser) parsePrimaryExpr() Expr {
	x := p.parseOperand()
	return x
}

func (p *parser) parseUnaryExpr() Expr {
	switch p.peekNonSpace().typ {
	case itemHexValue:
		number := p.next()
		return &BasicLit{
			ValuePos: number.pos,
			Kind:     token.INT,
			Value:    number.val,
		}
	case itemIdentifier:
		identifier := p.next()
		return &Ident{
			NamePos: identifier.pos,
			Name:    identifier.val,
		}
	case itemMinus:
		operator := p.next()
		expr := p.parseUnaryExpr()
		return &UnaryExpr{
			OpPos: operator.pos,
			Op:    token.SUB,
			X:     expr,
		}
	case itemOpenParen:
		opening := p.next()
		expr := p.parseExpr()
		closing := p.expect(itemCloseParen, "parentheses expression")
		return &ParenExpr{
			Opening: opening.pos,
			Expr:    expr,
			Closing: closing.pos,
		}
	}

	return p.parsePrimaryExpr()
}

func (p *parser) tokPrec() (item, token.Token, int) {
	tok := p.nextNonSpace()
	switch tok.typ {
	case itemMinus:
		return tok, token.SUB, 4
	case itemSlash:
		return tok, token.QUO, 5
	}
	p.backup()
	return tok, token.ILLEGAL_TOKEN, 0
}

func (p *parser) parseBinaryExpr(prec1 int) Expr {
	x := p.parseUnaryExpr()
	for {
		op, tok, oprec := p.tokPrec()
		if oprec < prec1 {
			return x
		}
		y := p.parseBinaryExpr(oprec + 1)
		x = &BinaryExpr{
			X:     x,
			OpPos: op.pos,
			Op:    tok,
			Y:     y,
		}
	}
}

func (p *parser) parseExpr() Expr {
	return p.parseBinaryExpr(1)
}

func (p *parser) parseArgList() *ArgList {
	if p.peek().typ != itemOpenParen {
		return nil
	}

	open := p.next()
	var list []*Ident
	if p.peek().typ == itemIdentifier {
		id := p.next()
		list = append(list, &Ident{
			NamePos: id.pos,
			Name:    id.val,
		})

		for p.peek().typ == itemComma {
			id = p.expect(itemIdentifier, "macro argument list")
			list = append(list, &Ident{
				NamePos: id.pos,
				Name:    id.val,
			})
		}
	}
	closing := p.expect(itemCloseParen, "macro argument list")

	return &ArgList{
		Opening: open.pos,
		List:    list,
		Closing: closing.pos,
	}
}

func (p *parser) parseMacroDir() Dir {
	keyword := p.expect(itemDefine, "macro definition")
	name := p.expect(itemIdentifier, "macro definition")
	args := p.parseArgList()
	switch p.peek().typ {
	case itemSpace, itemEOF:
		if p.peek().val == "" || strings.Contains(p.peek().val, "\n") {
			return &MacroDir{
				DirPos: keyword.pos,
				Name: &Ident{
					NamePos: name.pos,
					Name:    name.val,
				},
				Args:  args,
				Value: nil,
			}
		}
		fallthrough
	default:
		p.next()
		return &MacroDir{
			DirPos: keyword.pos,
			Name: &Ident{
				NamePos: name.pos,
				Name:    name.val,
			},
			Args:  args,
			Value: p.parseExpr(),
		}
	}
	/*
		p.next()
		switch p.peek().typ {
		case itemHexValue:
			//i2 := p.nextNonSpace()
			return &MacroDir{
				DirPos: keyword.pos,
				Name: &Ident{
					NamePos: name.pos,
					Name:    name.val,
				},
				Args:  args,
				Value: p.parseExpr(),
			}
		default:
			i1 := p.next()
			i2 := p.peek()
			//fmt.Printf("%s\n", i1.typ)
			//fmt.Printf("%s\n", i2.typ)
			for i2.typ != itemEOF && i2.typ != itemError && i2.typ != itemSpace {
				i1 = p.next()
				i2 = p.peek()
			}
			p.backup()
			return &BadDir{
				From: keyword.pos,
				To:   token.Pos(int(i1.pos) + len(i1.val)),
			}
			//p.errorf("unexpected %s", p.peek().typ)
		}
	}*/
}

func (p *parser) parseParenExpr() Expr {
	opening := p.expect(itemOpenParen, "parentheses expression")
	var expr Expr
	if p.peek().typ != itemCloseParen {
		expr = p.parseExpr()
	}
	closing := p.expect(itemCloseParen, "parentheses expression")
	return &ParenExpr{
		Opening: opening.pos,
		Expr:    expr,
		Closing: closing.pos,
	}
}

func (p *parser) next() item {
	if p.peekCount > 0 {
		p.peekCount--
	} else {
		p.token[0] = p.lex.nextItem()
	}
	p.pos, p.tok, p.lit = p.token[p.peekCount].pos, token.INT, p.token[p.peekCount].val
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

func (p *parser) nextNonSpace() (token item) {
	for {
		token = p.next()
		if token.typ != itemSpace {
			break
		}
	}
	return token
}

func (p *parser) peekNonSpace() (token item) {
	for {
		token = p.next()
		if token.typ != itemSpace {
			break
		}
	}
	p.backup()
	return token
}

func (p *parser) errorf(format string, args ...interface{}) {
	format = fmt.Sprintf("cgen: %s:%d: %s", p.name, p.token[0].line, format)
	panic(fmt.Errorf(format, args...))
}

func (p *parser) expect(expected itemType, context string) item {
	token := p.nextNonSpace()
	if token.typ != expected {
		p.unexpected(token, context)
	}
	return token
}

func (p *parser) expectOneOf(expected1, expected2 itemType, context string) item {
	token := p.nextNonSpace()
	if token.typ != expected1 && token.typ != expected2 {
		p.unexpected(token, context)
	}
	return token
}

func (p *parser) unexpected(token item, context string) {
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
