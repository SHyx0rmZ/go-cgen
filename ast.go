package cgen

import "github.com/SHyx0rmZ/cgen/token"

type Node interface {
	Pos() token.Pos
	End() token.Pos
}

type Expr interface {
	Node
	exprNode()
}

type Stmt interface {
	Node
	stmtNode()
}

type Decl interface {
	Node
	declNode()
}

// ----------------------------------------------------------------------------
// Comments

type Comment struct {
	Slash token.Pos
	Text  string
}

func (c *Comment) Pos() token.Pos { return c.Slash }
func (c *Comment) End() token.Pos { return token.Pos(int(c.Slash) + len(c.Text)) }

// An expression is represented by a tree consisting of one
// or more of the following concrete expression nodes.
//
type (
	// A BadExpr node is a placeholder for expressions containing
	// syntax errors for which no correct expression nodes can be
	// created.
	//
	BadExpr struct {
		From, To token.Pos // position range of bad expression
	}

	// An Ident node represents an identifier.
	Ident struct {
		NamePos token.Pos // identifier position
		Name    string    // identifier name
		// Obj *Object
	}

	// A BasicLit node represents a literal of basic type.
	BasicLit struct {
		ValuePos token.Pos   // literal position
		Kind     token.Token // token.INT
		Value    string      // literal string
	}
)

func (x *BadExpr) Pos() token.Pos  { return x.From }
func (x *Ident) Pos() token.Pos    { return x.NamePos }
func (x *BasicLit) Pos() token.Pos { return x.ValuePos }

func (x *BadExpr) End() token.Pos  { return x.To }
func (x *Ident) End() token.Pos    { return token.Pos(int(x.NamePos) + len(x.Name)) }
func (x *BasicLit) End() token.Pos { return token.Pos(int(x.ValuePos) + len(x.Value)) }

// exprNode() ensures that only expression/type nodes can be
// assigned to an Expr.
//
func (*BadExpr) exprNode()  {}
func (*Ident) exprNode()    {}
func (*BasicLit) exprNode() {}

func (id *Ident) String() string {
	if id != nil {
		return id.Name
	}
	return "<nil>"
}

// ----------------------------------------------------------------------------
// Pre-processor directives

type Dir interface {
	Node
	dirNode()
}

type (
	DefineDir struct {
		DirPos token.Pos
		Name   *Ident
		Value  Expr //todo
	}

	IncludeDir struct {
		DirPos  token.Pos
		PathPos token.Pos
		Path    string
	}

	IfDefDir struct {
		DirPos token.Pos
		Dir    string
		Name   *Ident
	}
)

func (d *DefineDir) Pos() token.Pos  { return d.Name.Pos() }
func (d *IncludeDir) Pos() token.Pos { return d.DirPos }
func (d *IfDefDir) Pos() token.Pos   { return d.DirPos }

func (d *DefineDir) End() token.Pos {
	if d.Value != nil {
		return d.Value.End()
	}
	return d.Name.End()
}
func (d *IncludeDir) End() token.Pos { return token.Pos(int(d.PathPos) + len(d.Path)) }
func (d *IfDefDir) End() token.Pos   { return d.Name.End() }

// ----------------------------------------------------------------------------
// Statements

// A statement is represented by a tree consisting of one
// or more of the following concrete statement nodes.
//
type (
	// A BadStmt node is a placeholder for statements containing
	// syntax errors for which no correct statement nodes can be
	// created.
	//
	BadStmt struct {
		From, To token.Pos
	}
)

func (StructDecl) exprNode() {}

type DefineStmt struct {
	Name  Ident
	Value Expr
}

type IncludeStmt struct {
	PathPos Pos
	Path    string
}

type IfDefinedStmt struct {
	Hash  Pos
	Value Ident
}

type IfNotDefinedStmt struct {
	Hash  Pos
	Value Ident
}

type TypeDecl struct {
	Name Ident
	Expr Expr
}

type StructDecl struct {
	Nodes []Expr
}

type StructType struct {
	Name Ident
}

func (StructType) exprNode() {}

type EnumDecl struct {
	Specs []EnumSpec
}

func (EnumDecl) exprNode() {}

type EnumSpec interface {
	Node
	enumSpecNode()
}

type EnumValue struct {
	Name  Ident
	Value *BasicLit
}

func (EnumValue) exprNode()     {}
func (EnumValue) enumSpecNode() {}

type BinaryExpr struct {
	X     Expr
	OpPos Pos
	Op    Token
	Y     Expr
}

type EnumConstExpr struct {
	Name Ident
	Expr Expr
}

func (EnumConstExpr) exprNode()     {}
func (EnumConstExpr) enumSpecNode() {}

func (BinaryExpr) exprNode()      {}
func (BinaryExpr) enumSpecValue() {} // TODO: revert hack
