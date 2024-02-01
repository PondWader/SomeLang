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
		node, _ := p.ParseFullIdentifierExpression(&nodes.Identifier{Name: token.Literal}, typeDef)
		return node
	case TokenEOF:
		return nil
	default:
		p.ThrowSyntaxError("Unexpected token \"", token.Literal, "\".")
	}
	return nil
}

// Parses everything that follows an identifier to parse function calls and key access
func (p *Parser) ParseFullIdentifierExpression(value nodes.Node, def TypeDef) (nodes.Node, TypeDef) {
	switch p.lexer.PeekOrExit().Type {
	case TokenLeftBracket:
		p.lexer.Next()

		funcDef, ok := def.(FuncDef)
		if !ok {
			p.ThrowTypeError("Cannot call a non-function value")
		}

		args := make([]nodes.Node, len(funcDef.Args))
		for i := 0; ; i++ {
			if token := p.lexer.PeekOrExit(); token.Type == TokenRightBracket {
				p.lexer.Next()
				break
			}
			if i >= len(args) {
				p.ThrowTypeError("Too many arguments passed to function.")
			}

			val, valDef := p.ParseValue()
			if !valDef.Equals(funcDef.Args[i]) {
				p.ThrowTypeError("Incorrect type passed for argument ", i+1, " of function call.")
			}

			args[i] = val
			token := p.ExpectToken(TokenRightBracket, TokenComma)
			if token.Type == TokenRightBracket {
				break
			}
		}

		return p.ParseFullIdentifierExpression(&nodes.FuncCall{
			Args:     args,
			Function: value,
		}, funcDef.ReturnType)
	/*case TokenPeriod:
	p.lexer.Next()
	identToken := p.ExpectToken(TokenIdentifier)
	return p.ParseFullIdentifierExpression(&nodes.KeyAccess{
		Object: value,
		Key:    identToken.Literal,
	})*/
	case TokenEquals:
		p.lexer.Next()
		nextToken := p.lexer.PeekOrExit()
		// Check if the expression is a comparison, otherwise it's an assignment expression
		if nextToken.Type == TokenEquals {

		}

		newVal, newType := p.ParseValue()
		if ident, ok := value.(*nodes.Identifier); ok {
			identType := p.currentTypeEnv.Get(ident.Name).(TypeDef)
			if !identType.Equals(newType) {
				p.ThrowTypeError("Cannot assign new type to variable \"", ident.Name, "\".")
			}

			return &nodes.Assignment{
				Identifier: ident.Name,
				NewValue:   newVal,
			}, newType
		} else if _, ok := value.(*nodes.KeyAccess); ok {

		} else {
			p.ThrowSyntaxError("Left hand side of assignment is not assignable.")
		}
	}
	return value, def
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
		typeDef, ok := p.currentTypeEnv.Get(token.Literal).(TypeDef)
		if !ok {
			p.ThrowTypeError(token.Literal, " is not defined in this scope.")
		}
		return p.ParseFullIdentifierExpression(&nodes.Identifier{Name: token.Literal}, typeDef)
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
	inner := p.ParseBlock(args)

	return &nodes.FuncDeclaration{
		Name:     funcName,
		ArgNames: argNames,
		Inner:    inner,
		Line:     p.lexer.GetCurrentLine(),
	}
}

func (p *Parser) ParseBlock(scopedVariables map[string]TypeDef) *nodes.Block {
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
