package interpreter

import (
	"fmt"
	"main/interpreter/environment"
	"main/interpreter/nodes"
)

func (p *Parser) ParseVarDeclaration() environment.Node {
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

	return &nodes.Assignment{
		Identifier: identifier,
		NewValue:   valNode,
		Depth:      0,
	}
}

func (p *Parser) ParseFunctionDeclaration() environment.Node {
	funcName, argDefs, argNames, returnType := p.ParseFunctionDef()

	p.currentTypeEnv.Set(funcName, FuncDef{
		GenericTypeDef{TypeFunc},
		argDefs,
		false,
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

func (p *Parser) ParseIfStatement() environment.Node {
	val, valDef := p.ParseValue(nil)
	if !valDef.Equals(GenericTypeDef{TypeBool}) {
		p.ThrowTypeError("If statement must be followed by a bool value.")
	}
	inner := p.ParseBlock(make(map[string]TypeDef), nil)

	var elseNode environment.Node
	// Check for else statement
	if token := p.lexer.NextOrExit(); token.Type == TokenElseStatement {
		token = p.ExpectToken(TokenIfStatement, TokenLeftBrace)
		if token.Type == TokenIfStatement {
			elseNode = p.ParseIfStatement()
		} else {
			p.lexer.Unread(token)
			elseNode = p.ParseBlock(make(map[string]TypeDef), nil)
		}
	} else {
		p.lexer.Unread(token)
	}

	return &nodes.IfStatement{
		Condition: val,
		Inner:     inner,
		Else:      elseNode,
	}
}

func (p *Parser) ParseImportStatement() environment.Node {
	module := p.ExpectToken(TokenString).Literal

	moduleDef := p.modules[module]
	if moduleDef == nil {
		p.ThrowSyntaxError("Module \"", module, "\" does not exist")
	}
	identifier := module
	if token := p.lexer.NextOrExit(); token.Type == TokenAsStatement {
		identifier = p.ExpectToken(TokenIdentifier).Literal
	} else {
		p.lexer.Unread(token)
	}

	p.currentTypeEnv.Set(identifier, ModuleDef{
		GenericTypeDef: GenericTypeDef{TypeModule},
		Properties:     moduleDef,
	})
	return &nodes.Import{
		Module:     module,
		Identifier: identifier,
	}
}

func (p *Parser) ParseForStatement() environment.Node {
	// Should support commas i.e. v, i range ["a", "b"] (left side value, right side index)
	token := p.ExpectToken(TokenIdentifier)
	ident := token.Literal
	p.ExpectToken(TokenRangeStatement)
	iterableValue, def := p.ParseValue(nil)
	fmt.Println(ident, iterableValue, def)
	return nil
}

func (p *Parser) ParseWhileStatement() environment.Node {
	val, def := p.ParseValue(GenericTypeDef{TypeBool})
	if !def.Equals(GenericTypeDef{TypeBool}) {
		p.ThrowTypeError("Value in while statement must be of type boolean")
	}
	return &nodes.WhileStatement{
		Condition: val,
		Inner:     p.ParseBlock(map[string]TypeDef{}, nil),
	}
}
