package cgen

import "fmt"

type Pos int

type item struct {
	typ  itemType
	pos  Pos
	val  string
	line int
}

func (i item) String() string {
	switch {
	case i.typ == itemEOF:
		return "EOF"
	case i.typ == itemError:
		return i.val
	case len(i.val) > 10:
		return fmt.Sprintf("%s(%.30q...)", i.typ, i.val)
	}
	return fmt.Sprintf("%s(%q)", i.typ, i.val)
}

type itemType int

const (
	itemError itemType = iota
	itemEOF
	itemComment
	itemIfNotDefined
	itemIfDefined
	itemDefine
	itemHexValue
	itemInclude
	itemExternC
	itemOpenCurly
	itemCloseCurly
	itemOpenParen
	itemCloseParen
	itemIdentifier
	itemEndIf
	itemLogicalOr
	itemBitwiseOr
	itemExtern
	itemSemicolon
	itemStar
	itemComma
	itemEqualSign
	itemLogicalAnd
	itemBitwiseAnd
	itemSpace
	itemIncludePath
	itemIncludePathSystem
)

func (t itemType) String() string {
	switch t {
	case itemComment:
		return "itemComment"
	case itemIfNotDefined:
		return "itemIfNotDefined"
	case itemIfDefined:
		return "itemIfDefined"
	case itemDefine:
		return "itemDefine"
	case itemHexValue:
		return "itemHexValue"
	case itemInclude:
		return "itemInclude"
	case itemExternC:
		return "itemExternC"
	case itemOpenCurly:
		return "itemOpenCurly"
	case itemCloseCurly:
		return "itemCloseCurly"
	case itemEndIf:
		return "itemEndIf"
	case itemOpenParen:
		return "itemOpenParen"
	case itemCloseParen:
		return "itemCloseParen"
	case itemIdentifier:
		return "itemIdentifier"
	case itemLogicalOr:
		return "itemLogicalOr"
	case itemBitwiseOr:
		return "itemBitwiseOr"
	case itemExtern:
		return "itemExtern"
	case itemSemicolon:
		return "itemSemicolon"
	case itemStar:
		return "itemStar"
	case itemComma:
		return "itemComma"
	case itemEqualSign:
		return "itemEqualSign"
	case itemLogicalAnd:
		return "itemLogicalAnd"
	case itemBitwiseAnd:
		return "itemBitwiseAnd"
	case itemSpace:
		return "itemSpace"
	case itemIncludePath:
		return "itemIncludePath"
	case itemIncludePathSystem:
		return "itemIncludePathSystem"
	default:
		return "itemUnknown"
	}
}
