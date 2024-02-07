package interpreter

import (
	"main/interpreter/nodes"
	"strconv"
)

// Parses everything that follows a value to parse the full value expression.
// This includes things such as function calls, key access, comparisons, operations etc
func (p *Parser) ParseValueExpression(value nodes.Node, def TypeDef) (nodes.Node, TypeDef) {
	token := p.lexer.NextOrExit()
	switch token.Type {
	case TokenLeftBracket:
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

			var argDef TypeDef
			// If the function is variadic, there can be an infinite number of args of the last arg type
			if i >= len(args) && funcDef.Variadic {
				args = append(args, nil)
				argDef = funcDef.Args[len(funcDef.Args)-1]
			} else if i >= len(args) {
				p.ThrowTypeError("Too many arguments passed to function.")
			} else {
				argDef = funcDef.Args[i]
			}

			val, valDef := p.ParseValue(argDef)
			if !valDef.Equals(argDef) {
				p.ThrowTypeError("Incorrect type passed for argument ", i+1, " of function call.")
			}

			args[i] = val
			token := p.ExpectToken(TokenRightBracket, TokenComma)
			if token.Type == TokenRightBracket {
				break
			}
		}

		return p.ParseValueExpression(&nodes.FuncCall{
			Args:     args,
			Function: value,
		}, funcDef.ReturnType)
	/*case TokenPeriod:
	p.lexer.Next()
	identToken := p.ExpectToken(TokenIdentifier)
	return p.ParseValueExpression(&nodes.KeyAccess{
		Object: value,
		Key:    identToken.Literal,
	})*/
	case TokenLeftSquareBracket:
	case TokenEquals:
		if p.lexer.PeekOrExit().Type == TokenEquals { // Comparison
			p.lexer.Next()
			rhsVal, rhsValDef := p.ParsePartialValue(def)
			if !rhsValDef.Equals(def) {
				p.ThrowTypeError("Right hand side of comparison must be the same type as the left hand side.")
			}
			return &nodes.Comparison{
				Type:      nodes.ComparisonEquals,
				LeftSide:  value,
				RightSide: rhsVal,
			}, GenericTypeDef{TypeBool}
		}

		if ident, ok := value.(*nodes.Identifier); ok {
			newVal, newValDef := p.ParsePartialValue(def)
			if !def.Equals(newValDef) {
				p.ThrowTypeError("Cannot assign new type to variable \"", ident.Name, "\".")
			}

			return &nodes.Assignment{
				Identifier: ident.Name,
				NewValue:   newVal,
			}, def
		} else if _, ok := value.(*nodes.KeyAccess); ok {

		} else {
			p.ThrowSyntaxError("Left hand side of assignment is not assignable.")
		}
	}
	p.lexer.Unread(token)
	return value, def
}

func (p *Parser) ParseLogicalOperator(value nodes.Node, def TypeDef) (nodes.Node, TypeDef) {
	token := p.lexer.NextOrExit()

	switch token.Type {
	case TokenAmpersand:
		p.ExpectToken(TokenAmpersand)

		if !def.Equals(GenericTypeDef{TypeBool}) {
			p.ThrowTypeError("Operation && can only be used on boolean value.")
		}

		rhsVal, rhsValDef := p.ParsePartialValue(nil)
		if !rhsValDef.Equals(GenericTypeDef{TypeBool}) {
			p.ThrowTypeError("Right hand side of && operation must be a boolean value.")
		}

		return &nodes.And{
			LeftSide:  value,
			RightSide: rhsVal,
		}, GenericTypeDef{TypeBool}
	}

	p.lexer.Unread(token)
	return value, def
}

// Parses a value of any type, without accounting for logical operations that follow it.
func (p *Parser) ParsePartialValue(implicitType TypeDef) (nodes.Node, TypeDef) {
	token := p.ExpectToken(TokenString, TokenNumber, TokenIdentifier, TokenTrue, TokenFalse)
	switch token.Type {
	case TokenString:
		return p.ParseValueExpression(&nodes.Value{Value: token.Literal}, GenericTypeDef{TypeString})
	case TokenTrue:
		return p.ParseValueExpression(&nodes.Value{Value: true}, GenericTypeDef{TypeBool})
	case TokenFalse:
		return p.ParseValueExpression(&nodes.Value{Value: false}, GenericTypeDef{TypeBool})
	case TokenNumber:
		// Check for decimal point, in which case it's a float
		if p.lexer.PeekOrExit().Type == TokenPeriod {
			p.lexer.Next()
			decimalNum := p.ExpectToken(TokenNumber)
			val, _ := strconv.ParseFloat(token.Literal+"."+decimalNum.Literal, 64)
			if implicitVal := ConvertFloat64ToTypeDef(val, implicitType.GetGenericType()); implicitVal != nil {
				return &nodes.Value{Value: implicitVal}, implicitType
			}
			return p.ParseValueExpression(&nodes.Value{Value: val}, GenericTypeDef{TypeFloat64})
		}
		val, _ := strconv.ParseInt(token.Literal, 10, 64)

		if implicitVal := ConvertInt64ToTypeDef(val, implicitType.GetGenericType()); implicitVal != nil {
			return p.ParseValueExpression(&nodes.Value{Value: implicitVal}, implicitType)
		}
		return p.ParseValueExpression(&nodes.Value{Value: val}, GenericTypeDef{TypeInt64})
	case TokenIdentifier:
		typeDef, ok := p.currentTypeEnv.Get(token.Literal).(TypeDef)
		if !ok {
			p.ThrowTypeError(token.Literal, " is not defined in this scope.")
		}
		return p.ParseValueExpression(&nodes.Identifier{Name: token.Literal}, typeDef)
	}
	return nil, nil
}

// Parses a value of any type, without accounting for operations that follow it.
//
// If implicitType is passed, the value will be coerced to the implicit type if possible.
func (p *Parser) ParseValue(implicitType TypeDef) (nodes.Node, TypeDef) {
	return p.ParseLogicalOperator(p.ParsePartialValue(implicitType))
}
