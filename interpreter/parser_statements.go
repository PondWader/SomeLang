package interpreter

import (
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

	p.currentTypeEnv.Set(funcName, NewFuncDef(
		argDefs,
		false,
		returnType,
	))

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
	valIdent := p.ExpectToken(TokenIdentifier).Literal

	indexIdent := ""
	if token := p.ExpectToken(TokenComma, TokenRangeStatement); token.Type == TokenComma {
		indexIdent = p.ExpectToken(TokenIdentifier).Literal
		p.ExpectToken(TokenRangeStatement)
	}

	iterableValue, def := p.ParseValue(nil)

	if def.IsInteger() {
		if indexIdent != "" {
			p.ThrowSyntaxError("Two values cannot be specified on the left hand side of an integer range loop.")
		}

		var startVal environment.Node = &nodes.Value{Value: 0}
		endVal := iterableValue
		if token := p.ExpectToken(TokenLeftBrace, TokenComma); token.Type == TokenComma {
			endVal, def = p.ParseValue(nil)
			if !def.IsInteger() {
				p.ThrowTypeError("Integers values must be used for an integer range loop.")
			}
			startVal = iterableValue
		} else {
			p.lexer.Unread(token)
		}
		return &nodes.LoopRange{
			ValIdentifier: valIdent,
			Start:         startVal,
			End:           endVal,
			Inner:         p.ParseBlock(map[string]TypeDef{valIdent: GenericTypeDef{TypeInt64}}, nil),
		}
	}

	arrayDef, ok := def.(ArrayDef)
	if !ok {
		p.ThrowTypeError("Right hand side of range loop must either be an integer or array.")
	}
	return GetGenericTypeNode(arrayDef.ElementType).GetLoopArray(valIdent, indexIdent, iterableValue, p.ParseBlock(
		map[string]TypeDef{valIdent: arrayDef.ElementType, indexIdent: GenericTypeDef{Type: TypeInt64}},
		nil,
	))
}

func (p *Parser) ParseWhileStatement() environment.Node {
	val, def := p.ParseValue(GenericTypeDef{TypeBool})
	if !def.Equals(GenericTypeDef{TypeBool}) {
		p.ThrowTypeError("Value in while statement must be of type boolean")
	}
	return &nodes.LoopWhile{
		Condition: val,
		Inner:     p.ParseBlock(map[string]TypeDef{}, nil),
	}
}
