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
}

func NewParser(name, input string) *parser {
	return &parser{lex: NewLexer(name, input), name: name}
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
			//fmt.Printf("%s\n", i.typ)
			switch i.typ {
			case itemEOF, itemError:
				return
			case itemEndIf:
				c <- &EndIfDir{
					DirPos: i.pos,
				}
			case itemDefine:
				if p.peekNonSpace().typ != itemIdentifier {
					p.errorf("expected identifier but found %s", p.peekNonSpace().typ)
				}
				i1 := p.next()
				//fmt.Printf("def %s %s\n", i1.typ, p.peek().typ)
				switch p.peek().typ {
				case itemSpace, itemEOF:
					if p.peek().val == "" || strings.Contains(p.peek().val, "\n") {
						//p.next()
						c <- &MacroDir{
							DirPos: i.pos,
							Name: &Ident{
								NamePos: i1.pos,
								Name:    i1.val,
							},
							Value: nil,
						}
						continue
					}
					p.next()
					switch p.peek().typ {
					case itemHexValue:
						i2 := p.nextNonSpace()
						c <- &MacroDir{
							DirPos: i.pos,
							Name: &Ident{
								NamePos: i1.pos,
								Name:    i1.val,
							},
							Value: &BasicLit{
								ValuePos: i2.pos,
								Kind:     token.INT,
								Value:    i2.val,
							},
						}
					default:
						i2 := p.peek()
						fmt.Printf("%s\n", i1.typ)
						fmt.Printf("%s\n", i2.typ)
						for i2.typ != itemEOF && i2.typ != itemError && i2.typ != itemSpace {
							i1 = p.next()
							i2 = p.peek()
						}
						p.backup()
						c <- &BadDir{
							From: i.pos,
							To:   token.Pos(int(i1.pos) + len(i1.val)),
						}
						//p.errorf("unexpected %s", p.peek().typ)
					}
				}
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
				i1 := p.expect(itemHexValue, "unary expression")
				c <- &UnaryExpr{
					OpPos: i.pos,
					Op:    token.SUB,
					X: &BasicLit{
						ValuePos: i1.pos,
						Kind:     token.INT,
						Value:    i1.val,
					},
				}
			default:
				fmt.Printf("%s\n", i)
			}
		}
	}()
	return c
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
