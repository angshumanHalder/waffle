package ast

import "monkeyinterpreter/token"

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type LetStatement struct {
	Value Expression
	Name  *Indentifier
	Token token.Token // the token.LET token
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

type Indentifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Indentifier) expressionNode() {}
func (i *Indentifier) TokenLiteral() string {
	return i.Token.Literal
}

type ReturnStatement struct {
	ReturnValue Expression
	Token       token.Token // the token.RETURN token
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}
