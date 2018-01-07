package lexer

import (
	"fmt"
	"github.com/SHyx0rmZ/cgen/token"
)

type Item struct {
	Typ  ItemType
	Pos  token.Pos
	Val  string
	Line int
}

func (i Item) String() string {
	switch {
	case i.Typ == ItemEOF:
		return "EOF"
	case i.Typ == ItemError:
		return i.Val
	case len(i.Val) > 10:
		return fmt.Sprintf("%s(%.30q...)", i.Typ, i.Val)
	}
	return fmt.Sprintf("%s(%q)", i.Typ, i.Val)
}

type ItemType int

const (
	ItemError ItemType = iota
	ItemEOF
	ItemComment
	ItemIfNotDefined
	ItemIfDefined
	ItemDefine
	ItemElseDir
	ItemHexValue
	ItemInclude
	ItemExternC
	ItemOpenCurly
	ItemCloseCurly
	ItemOpenParen
	ItemCloseParen
	ItemIdentifier
	ItemEndIf
	ItemLogicalOr
	ItemBitwiseOr
	ItemExtern
	ItemSemicolon
	ItemStar
	ItemComma
	ItemEqualSign
	ItemLogicalAnd
	ItemBitwiseAnd
	ItemSpace
	ItemIncludePath
	ItemIncludePathSystem
	ItemDecrement
	ItemMinus
	ItemSlash
	ItemString
)

func (t ItemType) String() string {
	switch t {
	case ItemComment:
		return "ItemComment"
	case ItemIfNotDefined:
		return "ItemIfNotDefined"
	case ItemIfDefined:
		return "ItemIfDefined"
	case ItemDefine:
		return "ItemDefine"
	case ItemHexValue:
		return "ItemHexValue"
	case ItemInclude:
		return "ItemInclude"
	case ItemExternC:
		return "ItemExternC"
	case ItemOpenCurly:
		return "ItemOpenCurly"
	case ItemCloseCurly:
		return "ItemCloseCurly"
	case ItemEndIf:
		return "ItemEndIf"
	case ItemOpenParen:
		return "ItemOpenParen"
	case ItemCloseParen:
		return "ItemCloseParen"
	case ItemIdentifier:
		return "ItemIdentifier"
	case ItemLogicalOr:
		return "ItemLogicalOr"
	case ItemBitwiseOr:
		return "ItemBitwiseOr"
	case ItemExtern:
		return "ItemExtern"
	case ItemSemicolon:
		return "ItemSemicolon"
	case ItemStar:
		return "ItemStar"
	case ItemComma:
		return "ItemComma"
	case ItemEqualSign:
		return "ItemEqualSign"
	case ItemLogicalAnd:
		return "ItemLogicalAnd"
	case ItemBitwiseAnd:
		return "ItemBitwiseAnd"
	case ItemSpace:
		return "ItemSpace"
	case ItemIncludePath:
		return "ItemIncludePath"
	case ItemIncludePathSystem:
		return "ItemIncludePathSystem"
	case ItemMinus:
		return "ItemMinus"
	case ItemDecrement:
		return "ItemDecrement"
	case ItemSlash:
		return "ItemSlash"
	case ItemString:
		return "ItemString"
	case ItemElseDir:
		return "ItemElseDir"
	default:
		return "itemUnknown"
	}
}
