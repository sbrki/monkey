package ast

import "github.com/sbrki/monkey/pkg/token"

type Node interface {
	TokenLiteral() string // used only for testing
}

type Statement interface {
	Node
	isStatementNode() // dummy for catching errors at compile time
}

type Expression interface {
	Node
	isExpressionNode() // dummy for catching errors at compile time
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) isStatementNode()     {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) isExpressionNode()    {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) isStatementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}
