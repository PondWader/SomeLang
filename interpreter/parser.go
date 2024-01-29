package interpreter

import (
	"fmt"
	"main/interpreter/environment"
	"main/interpreter/nodes"
	"os"
	"strconv"

	"github.com/logrusorgru/aurora/v4"
)

type Parser struct {
	lexer          *Lexer
	filePath       string
	currentTypeEnv *environment.Environment
}

// Creates abstract syntax tree
func Parse(content string, filePath string) []nodes.Node {
	p := &Parser{
		lexer:          NewLexer(content),
		filePath:       filePath,
		currentTypeEnv: environment.New(nil, environment.Call{}),
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
	defer func() {
		// Since this is a deferred function we must recover the panic
		// so that we the proper panic message can be displayed even if
		// ExpectToken would run in to an error.
		if err := recover(); err != nil {
			panic(err)
		}
		// Expect token ending the statement
		p.ExpectToken(TokenEOF, TokenNewLine, TokenSemiColon)
	}()

	switch token.Type {
	case TokenVarDeclaration:
		return p.ParseVarDeclaration()
	case TokenFunctionDeclarationStatement:
		return p.ParseFunctionDeclarationStatement()
	case TokenIdentifier:
		return p.ParseFullIdentifierExpression(&nodes.Identifier{Name: token.Literal})
	case TokenEOF:
		return nil
	default:
		p.ThrowSyntaxError("Unexpected token type \"", token.Literal, "\".")
	}
	return nil
}

// Parses everything that follows an identifier to parse function calls and key access
func (p *Parser) ParseFullIdentifierExpression(value nodes.Node) nodes.Node {
	switch p.lexer.PeekOrExit().Type {
	case TokenLeftBracket:
		p.lexer.Next()
		args := make([]nodes.Node, 0)
		for {
			if token := p.lexer.PeekOrExit(); token.Type == TokenRightBracket {
				p.lexer.Next()
				break
			}
			val, _ := p.ParseValue()
			args = append(args, val)
			token := p.ExpectToken(TokenRightBracket, TokenComma)
			if token.Type == TokenRightBracket {
				break
			}
		}
		return p.ParseFullIdentifierExpression(&nodes.FuncCall{
			Args:     args,
			Function: value,
		})
	case TokenPeriod:
		p.lexer.Next()
		identToken := p.ExpectToken(TokenIdentifier)
		return p.ParseFullIdentifierExpression(&nodes.KeyAccess{
			Object: value,
			Key:    identToken.Literal,
		})
	case TokenEquals:
		p.lexer.Next()
		nextToken := p.lexer.PeekOrExit()
		// Check if the expression is a comparison, otherwise it's an assignment expression
		if nextToken.Type == TokenEquals {

		}

		newVal, newType := p.ParseValue()
		if _, ok := value.(*nodes.Identifier); ok {

		} else if _, ok := value.(*nodes.KeyAccess); ok {

		} else {
			p.ThrowSyntaxError("Left hand side of assignment is not assignable.")
		}
	}
	return value
}

func (p *Parser) ParseValue() (nodes.Node, TypeDef) {
	token := p.ExpectToken(TokenString, TokenNumber, TokenIdentifier, TokenTrue, TokenFalse)
	switch token.Type {
	case TokenString:
		return &nodes.Value{Value: token.Literal}, GenericTypeDef{TypeString}
	case TokenTrue:
		return &nodes.Value{Value: true}, GenericTypeDef{TypeBool}
	case TokenFalse:
		return &nodes.Value{Value: false}, GenericTypeDef{TypeBool}
	case TokenNumber:
		// Check for decimal point
		if p.lexer.PeekOrExit().Type == TokenPeriod {
			p.lexer.Next()
			decimalNum := p.ExpectToken(TokenNumber)
			val, _ := strconv.ParseFloat(token.Literal+"."+decimalNum.Literal, 64)
			return &nodes.Value{Value: val}, GenericTypeDef{TypeFloat64}
		}
		val, _ := strconv.ParseInt(token.Literal, 10, 64)
		return &nodes.Value{Value: val}, GenericTypeDef{TypeInt64}
	case TokenIdentifier:
		identValueType, ok := p.currentTypeEnv.Get(token.Literal).(TypeDef)
		if !ok {
			p.ThrowTypeError(token.Literal, " is not defined in this scope")
		}
		return p.ParseFullIdentifierExpression(&nodes.Identifier{Name: token.Literal}), identValueType
	}
	return nil, nil
}

func (p *Parser) ParseVarDeclaration() nodes.Node {
	token := p.ExpectToken(TokenIdentifier)
	identifier := token.Literal
	p.ExpectToken(TokenEquals)
	valNode, valType := p.ParseValue()

	p.currentTypeEnv.Set(identifier, valType)

	return &nodes.VarDeclaration{
		Identifier: identifier,
		Value:      valNode,
	}
}

func (p *Parser) ParseFunctionDeclarationStatement() nodes.Node {
	funcName, args, returnType := p.ParseFunctionDef()

	inner := p.ParseBlock(args)

	argNames := make([]string, len(args))
	i := 0
	for name, _ := range args {
		argNames[i] = name
		i++
	}

	p.currentTypeEnv.Set(funcName, FuncDef{
		GenericTypeDef{TypeFunc},
		args,
		returnType,
	})

	return &nodes.FuncDeclaration{
		Name:     funcName,
		ArgNames: argNames,
		Inner:    inner,
		Line:     p.lexer.GetCurrentLine(),
	}
}

func (p *Parser) ParseBlock(scopedVariables map[string]TypeDef) []nodes.Node {
	ast := make([]nodes.Node, 0)
	p.ExpectToken(TokenLeftBrace)

	p.currentTypeEnv = p.currentTypeEnv.NewChild(environment.Call{})
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

	p.currentTypeEnv = p.currentTypeEnv.GetParent()
	return ast
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
