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

func (l *lexer) emit(t ItemType) {
	l.items <- Item{t, l.start, l.input[l.start:l.pos], l.line}
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
	l.items <- Item{ItemError, l.start, fmt.Sprintf(format, args...), l.line}
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
		l.emit(ItemSpace)
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
		l.emit(ItemElseDir)
		return lexLineStart
	case strings.HasPrefix(l.input[l.pos:], "#endif"):
		l.pos += token.Pos(len("#endif"))
		l.emit(ItemEndIf)
		return lexLineStart
	case l.accept(groupDigits):
		return lexHexValue
	case l.peek() == '{':
		l.next()
		l.emit(ItemOpenCurly)
		return lexLineStart
	case l.peek() == '}':
		l.next()
		l.emit(ItemCloseCurly)
		return lexLineStart
	case l.peek() == '(':
		l.next()
		l.emit(ItemOpenParen)
		return lexLineStart
	case l.peek() == ')':
		l.next()
		l.emit(ItemCloseParen)
		return lexLineStart
	case l.peek() == '\\':
		l.next()
		l.ignore()
		return lexLineStart
	case l.peek() == '|':
		l.next()
		if l.accept("|") {
			l.emit(ItemLogicalOr)
			return lexLineStart
		}
		l.emit(ItemBitwiseOr)
		return lexLineStart
	case l.peek() == ';':
		l.next()
		l.emit(ItemSemicolon)
		return lexLineStart
	case l.peek() == '*':
		l.next()
		l.emit(ItemStar)
		return lexLineStart
	case l.peek() == ',':
		l.next()
		l.emit(ItemComma)
		return lexLineStart
	case l.peek() == '=':
		l.next()
		l.emit(ItemEqualSign)
		return lexLineStart
	case l.peek() == '&':
		l.next()
		if l.accept("&") {
			l.emit(ItemLogicalAnd)
			return lexLineStart
		}
		l.emit(ItemBitwiseAnd)
		return lexLineStart
	case l.peek() == '-':
		l.next()
		if l.accept("-") {
			l.emit(ItemDecrement)
			return lexLineStart
		}
		l.emit(ItemMinus)
		return lexLineStart
	case l.peek() == '/':
		l.next()
		l.emit(ItemSlash)
		return lexLineStart
	case l.peek() == '"':
		return lexString
	default:
		if l.accept("_" + groupLower + groupUpper) {
			return lexIdentifier
		}
		if int(l.pos) == len(l.input) {
			l.emit(ItemEOF)
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
			l.emit(ItemString)
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
	l.emit(ItemExtern)
	l.acceptRun(" ")
	l.emit(ItemSpace)
	//if strings.HasPrefix(l.input[l.Pos:], `"C"`) {
	//	l.Pos += token.Pos(len(`"C"`))
	//	l.emit(ItemExternC)
	//	//return lexLineStart
	//}
	return lexLineStart
}

func lexInclude(l *lexer) stateFn {
	l.pos += token.Pos(len("#include"))
	l.emit(ItemInclude)
	l.acceptRun(" ")
	l.emit(ItemSpace)
	switch l.next() {
	case '"':
		l.acceptRun(groupLower + groupUpper + groupDigits + "_-/\\.")
		if l.accept(`"`) {
			l.emit(ItemIncludePath)
			return lexLineStart
		}
		return l.errorf("expected closing quotes")
	case '<':
		l.acceptRun(groupLower + groupUpper + groupDigits + "_-/\\.")
		if l.accept(">") {
			l.emit(ItemIncludePathSystem)
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
	l.emit(ItemHexValue)
	return lexLineStart
}

func lexIdentifier(l *lexer) stateFn {
	//if l.accept("_" + groupLower + groupUpper) {
	l.acceptRun("_" + groupLower + groupUpper + groupDigits)
	l.emit(ItemIdentifier)
	return lexLineStart
	//}
	//return l.errorf("expected identifier")
}

func lexDefine(l *lexer) stateFn {
	l.pos += token.Pos(len("#define"))
	l.emit(ItemDefine)
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
	l.emit(ItemIfDefined)
	return lexLineStart
}

func lexIfNotDefined(l *lexer) stateFn {
	l.pos += token.Pos(len("#ifndef"))
	l.emit(ItemIfNotDefined)
	return lexLineStart
}

func lexMultilineComment(l *lexer) stateFn {
	l.pos += 2
	for {
		n := l.next()
		for n != '*' {
			if n == eof {
				l.emit(ItemEOF)
				return nil
			}
			n = l.next()
		}
		if l.peek() == '/' {
			l.next()
			l.emit(ItemComment)
			break
		}
	}
	return lexLineStart
}
