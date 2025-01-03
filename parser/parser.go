package parser

import (
	"fmt"

	"github.com/MohamTahaB/interpreter-go/ast"
	"github.com/MohamTahaB/interpreter-go/lexer"
	"github.com/MohamTahaB/interpreter-go/token"
)

type Parser struct {
	l *lexer.Lexer

	currToken token.Token
	peekToken token.Token

	errors []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// nextToken twice to populate both current and peek tokens
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.currToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{
		Token: p.currToken,
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.currToken,
		Value: p.currToken.Literal,
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: the expression block to be collected afterwards, basically the same thing until we get to the semicolon;
	for !p.currTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) expectPeek(tokType token.TokenType) bool {
	if p.peekTokenIs(tokType) {
		p.nextToken()
		return true
	}
	p.peekError(tokType)
	return false
}

func (p *Parser) peekTokenIs(tokType token.TokenType) bool {
	return p.peekToken.Type == tokType
}

func (p *Parser) currTokenIs(tokType token.TokenType) bool {
	return p.currToken.Type == tokType
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(tokType token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", tokType, p.peekToken.Type)

	p.errors = append(p.errors, msg)
}
