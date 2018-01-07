package lexer

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/SHyx0rmZ/cgen/token"
)

const eof = -1

type stateFn func(*lexer) stateFn

type Lexer interface {
	NextItem() Item
}

type lexer struct {
	name    string
	input   string
	state   stateFn
	pos     token.Pos
	start   token.Pos
	width   token.Pos
	lastPos token.Pos
	items   chan Item
	line    int
}

func NewLexer(name, input string) *lexer {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan Item),
		line:  1,
	}
	go l.run()
	return l
}

func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = token.Pos(w)
	l.pos += l.width
	if r == '\n' {
		l.line++
	}
	return r
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) backup() {
	l.pos -= l.width
	if l.width == 1 && l.input[l.pos] == '\n' {
		l.line--
	}
}

func (l *lexer) emit(t token.Token) {
	l.items <- Item{l.start, l.input[l.start:l.pos], t, l.line}
	l.start = l.pos
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- Item{l.start, fmt.Sprintf(format, args...), token.ILLEGAL, l.line}
	return nil
}

func (l *lexer) NextItem() Item {
	item := <-l.items
	l.lastPos = item.Pos
	return item
}

func (l *lexer) drain() {
	for range l.items {
	}
}

func (l *lexer) run() {
	for l.state = lexLineStart; l.state != nil; {
		l.state = l.state(l)
	}
	close(l.items)
}

func lexLineStart(l *lexer) stateFn {
	switch {
	case strings.HasPrefix(l.input[l.pos:], "/*"):
		return lexMultilineComment
	case l.accept(" \n\t"):
		l.acceptRun(" \n\t")
		l.emit(token.WHITESPACE)
		return lexLineStart
	case strings.HasPrefix(l.input[l.pos:], "#ifndef"):
		return lexIfNotDefined
	case strings.HasPrefix(l.input[l.pos:], "#ifdef"):
		return lexIfDefined
	case strings.HasPrefix(l.input[l.pos:], "#define"):
		return lexDefine
	case strings.HasPrefix(l.input[l.pos:], "#include"):
		return lexInclude
	case strings.HasPrefix(l.input[l.pos:], "extern"):
		return lexExtern
	case strings.HasPrefix(l.input[l.pos:], "#else"):
		l.pos += token.Pos(len("#else"))
		l.emit(token.ELSE)
		return lexLineStart
	case strings.HasPrefix(l.input[l.pos:], "#endif"):
		l.pos += token.Pos(len("#endif"))
		l.emit(token.ENDIF)
		return lexLineStart
	case l.accept(groupDigits):
		return lexHexValue
	case l.peek() == '{':
		l.next()
		l.emit(token.LBRACE)
		return lexLineStart
	case l.peek() == '}':
		l.next()
		l.emit(token.RBRACE)
		return lexLineStart
	case l.peek() == '(':
		l.next()
		l.emit(token.LPAREN)
		return lexLineStart
	case l.peek() == ')':
		l.next()
		l.emit(token.RPAREN)
		return lexLineStart
	case l.peek() == '\\':
		l.next()
		l.ignore()
		return lexLineStart
	case l.peek() == '|':
		l.next()
		if l.accept("|") {
			l.emit(token.LOR)
			return lexLineStart
		}
		l.emit(token.OR)
		return lexLineStart
	case l.peek() == ';':
		l.next()
		l.emit(token.SEMICOLON)
		return lexLineStart
	case l.peek() == '*':
		l.next()
		l.emit(token.MUL)
		return lexLineStart
	case l.peek() == ',':
		l.next()
		l.emit(token.COMMA)
		return lexLineStart
	case l.peek() == '=':
		l.next()
		l.emit(token.ASSIGN)
		return lexLineStart
	case l.peek() == '&':
		l.next()
		if l.accept("&") {
			l.emit(token.LAND)
			return lexLineStart
		}
		l.emit(token.AND)
		return lexLineStart
	case l.peek() == '-':
		l.next()
		if l.accept("-") {
			l.emit(token.DEC)
			return lexLineStart
		}
		l.emit(token.SUB)
		return lexLineStart
	case l.peek() == '/':
		l.next()
		l.emit(token.QUO)
		return lexLineStart
	case l.peek() == '"':
		return lexString
	default:
		if l.accept("_" + groupLower + groupUpper) {
			return lexIdentifier
		}
		if int(l.pos) == len(l.input) {
			l.emit(token.EOF)
			return nil
		}
		return l.errorf("unknown: %.30q...", l.input[l.pos:])
	}
}

const (
	groupLower  = "abcdefghijklmnopqrstuvwxyz"
	groupUpper  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	groupDigits = "0123456789"
)

func lexString(l *lexer) stateFn {
	l.next()
	for {
		n := l.peek()
		if n == '"' {
			l.next()
			l.emit(token.STRING)
			return lexLineStart
		}
		if n == '\\' {
			l.next()
			n = l.peek()
		}
		l.next()
	}
}

func lexExtern(l *lexer) stateFn {
	l.pos += token.Pos(len("extern"))
	l.emit(token.EXTERN)
	l.acceptRun(" ")
	l.emit(token.WHITESPACE)
	//if strings.HasPrefix(l.input[l.Pos:], `"C"`) {
	//	l.Pos += token.Pos(len(`"C"`))
	//	l.emit(ItemExternC)
	//	//return lexLineStart
	//}
	return lexLineStart
}

func lexInclude(l *lexer) stateFn {
	l.pos += token.Pos(len("#include"))
	l.emit(token.INCLUDE)
	l.acceptRun(" ")
	l.emit(token.WHITESPACE)
	switch l.next() {
	case '"':
		l.acceptRun(groupLower + groupUpper + groupDigits + "_-/\\.")
		if l.accept(`"`) {
			l.emit(token.INCLUDE_PATH)
			return lexLineStart
		}
		return l.errorf("expected closing quotes")
	case '<':
		l.acceptRun(groupLower + groupUpper + groupDigits + "_-/\\.")
		if l.accept(">") {
			l.emit(token.INCLUDE_PATH)
			return lexLineStart
		}
		return l.errorf("expected closing angle bracket")
	}
	l.backup()
	return l.errorf("expected include path")
}

func lexHexValue(l *lexer) stateFn {
	l.acceptRun("x0123456789abcdefABCDEF")
	l.accept("u")
	l.emit(token.INT)
	return lexLineStart
}

func lexIdentifier(l *lexer) stateFn {
	//if l.accept("_" + groupLower + groupUpper) {
	l.acceptRun("_" + groupLower + groupUpper + groupDigits)
	l.emit(token.IDENT)
	return lexLineStart
	//}
	//return l.errorf("expected identifier")
}

func lexDefine(l *lexer) stateFn {
	l.pos += token.Pos(len("#define"))
	l.emit(token.DEFINE)
	return lexLineStart
}

//func lexIfDefined(l *lexer) stateFn {
//	l.Pos += Pos(len("#ifdef"))
//	l.acceptRun(" ")
//	l.ignore()
//	if l.accept("_" + groupUpper + groupLower) {
//		l.acceptRun("_" + groupUpper + groupLower + groupDigits)
//		if l.peek() != '\n' {
//			return l.errorf("expected line break")
//		}
//		l.emit(ItemIfDefined)
//		return lexLineStart
//	}
//	return l.errorf("expected identifier")
//}

func lexIfDefined(l *lexer) stateFn {
	l.pos += token.Pos(len("#ifdef"))
	l.emit(token.IFDEF)
	return lexLineStart
}

func lexIfNotDefined(l *lexer) stateFn {
	l.pos += token.Pos(len("#ifndef"))
	l.emit(token.IFNDEF)
	return lexLineStart
}

func lexMultilineComment(l *lexer) stateFn {
	l.pos += 2
	for {
		n := l.next()
		for n != '*' {
			if n == eof {
				l.emit(token.EOF)
				return nil
			}
			n = l.next()
		}
		if l.peek() == '/' {
			l.next()
			l.emit(token.COMMENT)
			break
		}
	}
	return lexLineStart
}
