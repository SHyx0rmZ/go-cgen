package cgen

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/SHyx0rmZ/cgen/token"
)

const eof = -1

type stateFn func(*lexer) stateFn

type lexer struct {
	name    string
	input   string
	state   stateFn
	pos     token.Pos
	start   token.Pos
	width   token.Pos
	lastPos token.Pos
	items   chan item
	line    int
}

func NewLexer(name, input string) *lexer {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan item),
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

func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.start, l.input[l.start:l.pos], l.line}
	switch t {
	}
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
	l.items <- item{itemError, l.start, fmt.Sprintf(format, args...), l.line}
	return nil
}

func (l *lexer) nextItem() item {
	item := <-l.items
	l.lastPos = item.pos
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
		l.emit(itemSpace)
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
	case strings.HasPrefix(l.input[l.pos:], "#endif"):
		l.pos += token.Pos(len("#endif"))
		l.emit(itemEndIf)
		return lexLineStart
	case l.accept(groupDigits):
		return lexHexValue
	case l.peek() == '{':
		l.next()
		l.emit(itemOpenCurly)
		return lexLineStart
	case l.peek() == '}':
		l.next()
		l.emit(itemCloseCurly)
		return lexLineStart
	case l.peek() == '(':
		l.next()
		l.emit(itemOpenParen)
		return lexLineStart
	case l.peek() == ')':
		l.next()
		l.emit(itemCloseParen)
		return lexLineStart
	case l.peek() == '\\':
		l.next()
		l.ignore()
		return lexLineStart
	case l.peek() == '|':
		l.next()
		if l.accept("|") {
			l.emit(itemLogicalOr)
			return lexLineStart
		}
		l.emit(itemBitwiseOr)
		return lexLineStart
	case l.peek() == ';':
		l.next()
		l.emit(itemSemicolon)
		return lexLineStart
	case l.peek() == '*':
		l.next()
		l.emit(itemStar)
		return lexLineStart
	case l.peek() == ',':
		l.next()
		l.emit(itemComma)
		return lexLineStart
	case l.peek() == '=':
		l.next()
		l.emit(itemEqualSign)
		return lexLineStart
	case l.peek() == '&':
		l.next()
		if l.accept("&") {
			l.emit(itemLogicalAnd)
			return lexLineStart
		}
		l.emit(itemBitwiseAnd)
		return lexLineStart
	case l.peek() == '-':
		l.next()
		if l.accept("-") {
			l.emit(itemDecrement)
			return lexLineStart
		}
		l.emit(itemMinus)
		return lexLineStart
	case l.peek() == '/':
		l.next()
		l.emit(itemSlash)
		return lexLineStart
	case l.peek() == '"':
		return lexString
	default:
		if l.accept("_" + groupLower + groupUpper) {
			return lexIdentifier
		}
		if int(l.pos) == len(l.input) {
			l.emit(itemEOF)
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
			l.emit(itemString)
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
	l.emit(itemExtern)
	l.acceptRun(" ")
	l.emit(itemSpace)
	//if strings.HasPrefix(l.input[l.pos:], `"C"`) {
	//	l.pos += token.Pos(len(`"C"`))
	//	l.emit(itemExternC)
	//	//return lexLineStart
	//}
	return lexLineStart
}

func lexInclude(l *lexer) stateFn {
	l.pos += token.Pos(len("#include"))
	l.emit(itemInclude)
	l.acceptRun(" ")
	l.emit(itemSpace)
	switch l.next() {
	case '"':
		l.acceptRun(groupLower + groupUpper + groupDigits + "_-/\\.")
		if l.accept(`"`) {
			l.emit(itemIncludePath)
			return lexLineStart
		}
		return l.errorf("expected closing quotes")
	case '<':
		l.acceptRun(groupLower + groupUpper + groupDigits + "_-/\\.")
		if l.accept(">") {
			l.emit(itemIncludePathSystem)
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
	l.emit(itemHexValue)
	return lexLineStart
}

func lexIdentifier(l *lexer) stateFn {
	//if l.accept("_" + groupLower + groupUpper) {
	l.acceptRun("_" + groupLower + groupUpper + groupDigits)
	l.emit(itemIdentifier)
	return lexLineStart
	//}
	//return l.errorf("expected identifier")
}

func lexDefine(l *lexer) stateFn {
	l.pos += token.Pos(len("#define"))
	l.emit(itemDefine)
	return lexLineStart
}

//func lexIfDefined(l *lexer) stateFn {
//	l.pos += Pos(len("#ifdef"))
//	l.acceptRun(" ")
//	l.ignore()
//	if l.accept("_" + groupUpper + groupLower) {
//		l.acceptRun("_" + groupUpper + groupLower + groupDigits)
//		if l.peek() != '\n' {
//			return l.errorf("expected line break")
//		}
//		l.emit(itemIfDefined)
//		return lexLineStart
//	}
//	return l.errorf("expected identifier")
//}

func lexIfDefined(l *lexer) stateFn {
	l.pos += token.Pos(len("#ifdef"))
	l.emit(itemIfDefined)
	return lexLineStart
}

func lexIfNotDefined(l *lexer) stateFn {
	l.pos += token.Pos(len("#ifndef"))
	l.emit(itemIfNotDefined)
	return lexLineStart
}

func lexMultilineComment(l *lexer) stateFn {
	l.pos += 2
	for {
		n := l.next()
		for n != '*' {
			if n == eof {
				l.emit(itemEOF)
				return nil
			}
			n = l.next()
		}
		if l.peek() == '/' {
			l.next()
			l.emit(itemComment)
			break
		}
	}
	return lexLineStart
}
