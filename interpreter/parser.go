package interpreter

import (
	"fmt"
	"main/interpreter/nodes"
	"os"

	"github.com/logrusorgru/aurora/v4"
)

type Parser struct {
	lexer          *Lexer
	filePath       string
	currentTypeEnv *TypeEnvironment
}

// Creates abstract syntax tree
func Parse(content string, filePath string) []nodes.Node {
	p := &Parser{
		lexer:          NewLexer(content),
		filePath:       filePath,
		currentTypeEnv: NewTypeEnvironment(nil, nil),
	}
	ast := make([]nodes.Node, 0)
	for {
		node := p.ParseNext(false)
		if node == nil {
			break
		}
		ast = append(ast, node)
	}
	return ast
}

func (p *Parser) ParseNext(inBlock bool) nodes.Node {
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
				if token, err := p.lexer.Next(); err == nil && token.Type == TokenNewLine {
					return p.ParseNext(inBlock)
				}
			}
		}
	}

	defer func() {
		// Since this is a deferred function it must recover the panic
		// so that the proper panic message can be displayed even if
		// ExpectToken would run in to an error.
		if err := recover(); err != nil {
			panic(err)
		}
		// Expect token ending the statement
		if token := p.ExpectToken(TokenEOF, TokenNewLine, TokenSemiColon, TokenForwardSlash); token.Type == TokenForwardSlash {
			p.lexer.Unread(token) // If the token is the start of a comment, it should be unread so the next call of ParseNext() reads the start of the comment
		}
	}()

	switch token.Type {
	case TokenVarDeclaration:
		return p.ParseVarDeclaration()
	case TokenFunctionDeclaration:
		return p.ParseFunctionDeclaration()
	case TokenIdentifier:
		typeDef, ok := p.currentTypeEnv.Get(token.Literal).(TypeDef)
		if !ok {
			p.ThrowTypeError(token.Literal, " is not defined in this scope.")
		}
		node, _ := p.ParseFullValueExpression(&nodes.Identifier{Name: token.Literal}, typeDef)
		return node
	case TokenReturnStatement:
		if p.currentTypeEnv.ReturnType == nil {
			p.ThrowSyntaxError("You cannot use a return statement outside of a function with a defined return type.")
		}

		returnValue, returnValueDef := p.ParseValue(p.currentTypeEnv.ReturnType)
		if !returnValueDef.Equals(p.currentTypeEnv.ReturnType) {
			p.ThrowTypeError("Incorrect type of value returned.")
		}
		p.currentTypeEnv.Returned = true
		return &nodes.Return{
			Value: returnValue,
		}
	case TokenEOF:
		return nil
	default:
		p.ThrowSyntaxError("Unexpected token \"", token.Literal, "\".")
	}
	return nil
}

func (p *Parser) ParseVarDeclaration() nodes.Node {
	token := p.ExpectToken(TokenIdentifier)
	identifier := token.Literal

	token = p.lexer.NextOrExit()
	var typeDef TypeDef = GenericTypeDef{TypeNil}
	if token.Type != TokenEquals {
		p.lexer.Unread(token) // Unread token so it can be parsed as the type
		typeDef = p.ParseTypeDef()
		p.ExpectToken(TokenEquals)
	}

	valNode, valType := p.ParseValue(typeDef)
	if typeDef.GetGenericType() != TypeNil && !valType.Equals(typeDef) {
		p.ThrowTypeError("Incorrect type of value on right hand side of variable declaration.")
	}

	p.currentTypeEnv.Set(identifier, valType)

	return &nodes.VarDeclaration{
		Identifier: identifier,
		Value:      valNode,
	}
}

func (p *Parser) ParseFunctionDeclaration() nodes.Node {
	funcName, argDefs, argNames, returnType := p.ParseFunctionDef()

	p.currentTypeEnv.Set(funcName, FuncDef{
		GenericTypeDef{TypeFunc},
		argDefs,
		returnType,
	})

	args := make(map[string]TypeDef, len(argDefs))
	for i, name := range argNames {
		args[name] = argDefs[i]
	}
	inner := p.ParseBlock(args, returnType)

	return &nodes.FuncDeclaration{
		Name:     funcName,
		ArgNames: argNames,
		Inner:    inner,
		Line:     p.lexer.GetCurrentLine(),
	}
}

func (p *Parser) ParseBlock(scopedVariables map[string]TypeDef, returnType TypeDef) *nodes.Block {
	ast := make([]nodes.Node, 0)
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
