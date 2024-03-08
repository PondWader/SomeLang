package interpreter

import (
	"fmt"
	"main/interpreter/environment"
	"main/interpreter/nodes"
	"os"

	"github.com/logrusorgru/aurora/v4"
)

type Parser struct {
	lexer          *Lexer
	filePath       string
	currentTypeEnv *TypeEnvironment
	modules        map[string]map[string]TypeDef
}

// Creates abstract syntax tree
func Parse(content string, filePath string, globals map[string]TypeDef, modules map[string]map[string]TypeDef) []environment.Node {
	p := &Parser{
		lexer:          NewLexer(content),
		filePath:       filePath,
		currentTypeEnv: NewTypeEnvironment(nil, nil, 0),
		modules:        modules,
	}

	for name, def := range globals {
		p.currentTypeEnv.Set(name, def)
	}

	ast := make([]environment.Node, 0)
	for {
		node := p.ParseNext(false)
		if node == nil {
			break
		}
		ast = append(ast, node)
	}
	return ast
}

func (p *Parser) ParseNext(inBlock bool) environment.Node {
	token := p.lexer.NextOrExit()
	if token.Type == TokenNewLine || token.Type == TokenSemiColon {
		return p.ParseNext(inBlock)
	}
	if token.Type == TokenRightBrace && inBlock {
		return nil
	}

	// Check for comment
	if token.Type == TokenForwardSlash {
		if p.lexer.PeekOrExit().Type == TokenForwardSlash {
			for {
				if token, err := p.lexer.Next(); err == nil && token.Type == TokenNewLine || token.Type == TokenEOF {
					return p.ParseNext(inBlock)
				}
			}
		}
	}

	// If the current block has already returned, we don't want to read anymore statements
	// and instead can recursively read through all tokens until the closing curly right brace
	if p.currentTypeEnv.Returned {
		return p.ParseNext(inBlock)
	}

	defer func() {
		// Since this is a deferred function it must recover the panic
		// so that the proper panic message can be displayed even if
		// ExpectToken would run in to an error.
		if err := recover(); err != nil {
			panic(err)
		}
		// Expect token ending the statement
		if token := p.ExpectToken(TokenEOF, TokenNewLine, TokenSemiColon, TokenForwardSlash, TokenRightBrace); token.Type == TokenForwardSlash || token.Type == TokenRightBrace {
			p.lexer.Unread(token) // If the token is the start of a comment, it should be unread so the next call of ParseNext() reads the start of the comment
		}
	}()

	switch token.Type {
	case TokenVarDeclaration:
		return p.ParseVarDeclaration()
	case TokenFunctionDeclaration:
		return p.ParseFunctionDeclaration()
	case TokenIfStatement:
		return p.ParseIfStatement()
	case TokenIdentifier:
		typeDef, _ := p.currentTypeEnv.Get(token.Literal)
		if typeDef == nil {
			p.ThrowTypeError(token.Literal, " is not defined in this scope.")
		}
		node, _ := p.ParseOperator(p.ParseValueExpression(&nodes.Identifier{Name: token.Literal}, typeDef))
		return node
	case TokenReturnStatement:
		returnType := p.currentTypeEnv.GetReturnType()
		if returnType == nil {
			p.ThrowSyntaxError("You cannot use a return statement outside of a function with a defined return type.")
		}

		returnValue, returnValueDef := p.ParseValue(returnType)
		if !returnValueDef.Equals(returnType) {
			p.ThrowTypeError("Incorrect type of value returned.")
		}
		p.currentTypeEnv.SetReturned()
		return &nodes.Return{
			Value: returnValue,
		}
	case TokenForStatement:
		return p.ParseForStatement()
	case TokenStructDeclaration:
		return p.ParseStructDeclaration()
	case TokenImportStatement:
		return p.ParseImportStatement()
	case TokenWhileStatement:
		return p.ParseWhileStatement()
	case TokenEOF:
		return nil
	default:
		p.ThrowSyntaxError("Unexpected token \"", token.Literal, "\".")
	}
	return nil
}

func (p *Parser) ParseBlock(scopedVariables map[string]TypeDef, returnType TypeDef) *nodes.Block {
	ast := make([]environment.Node, 0)
	p.ExpectToken(TokenLeftBrace)

	p.currentTypeEnv = p.currentTypeEnv.NewChild(returnType)
	for name, valType := range scopedVariables {
		p.currentTypeEnv.Set(name, valType)
	}

	for {
		token := p.ParseNext(true)
		if token == nil {
			break
		}
		ast = append(ast, token)
	}

	if returnType != nil && !p.currentTypeEnv.Returned {
		p.ThrowTypeError("The function is missing a return statement.")
	}

	p.currentTypeEnv = p.currentTypeEnv.GetParent()
	return &nodes.Block{Nodes: ast}
}

func (p *Parser) ExpectToken(tokenType ...TokenType) Token {
	token := p.lexer.NextOrExit()
	for _, allowedType := range tokenType {
		if token.Type == allowedType {
			return token
		}
	}
	p.ThrowSyntaxError("Unexpected token \"", token.Literal, "\".")
	return Token{}
}

func (p *Parser) ThrowSyntaxError(msg ...any) {
	fmt.Println(aurora.Red("[ERROR]"), aurora.Gray(5, "Syntax error at line "+fmt.Sprint(p.lexer.GetCurrentLine())+":"))
	fmt.Println(" ", aurora.Gray(3, ">"), aurora.Red(fmt.Sprint(msg...)))
	fmt.Println(" ", aurora.Gray(18, "in "+p.filePath))
	os.Exit(1)
}

func (p *Parser) ThrowTypeError(msg ...any) {
	fmt.Println(aurora.Red("[ERROR]"), aurora.Gray(5, "Type error at line "+fmt.Sprint(p.lexer.GetCurrentLine())+":"))
	fmt.Println(" ", aurora.Gray(3, ">"), aurora.Red(fmt.Sprint(msg...)))
	fmt.Println(" ", aurora.Gray(18, "in "+p.filePath))
	os.Exit(1)
}
