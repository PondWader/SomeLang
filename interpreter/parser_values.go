package interpreter

import (
	"main/interpreter/environment"
	"main/interpreter/nodes"
	"strconv"
)

// Parses everything that follows a value to parse the full value expression.
// This includes things such as function calls, key access, comparisons, operations etc
func (p *Parser) ParseValueExpression(value environment.Node, def TypeDef) (environment.Node, TypeDef) {
	token := p.lexer.NextOrExit()
	switch token.Type {
	case TokenLeftBracket:
		funcDef, ok := def.(FuncDef)
		if !ok {
			p.ThrowTypeError("Cannot call a non-function value")
		}

		args := make([]environment.Node, len(funcDef.Args))
		for i := 0; ; i++ {
			if token := p.lexer.PeekOrExit(); token.Type == TokenRightBracket {
				p.lexer.Next()
				break
			} else if token.Type == TokenNewLine { // Allow new lines between arguments
				p.lexer.Next()
				i--
				continue
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
			token := p.ExpectToken(TokenRightBracket, TokenComma, TokenNewLine)
			if token.Type == TokenRightBracket {
				break
			}
		}

		return p.ParseValueExpression(&nodes.FuncCall{
			Args:     args,
			Function: value,
		}, funcDef.ReturnType)

	case TokenLeftSquareBracket:
		index, indexDef := p.ParseValue(nil)
		if !indexDef.IsInteger() {
			p.ThrowTypeError("Arrays must be indexed with an integer value")
		}
		arrayDef, ok := def.(ArrayDef)
		if !ok {
			p.ThrowTypeError("Cannot access index on non-array value")
		}
		p.ExpectToken(TokenRightSquareBracket)
		return GetGenericTypeNode(arrayDef.ElementType).GetArrayIndex(value, index), arrayDef.ElementType

	case TokenLeftBrace:

	}
	p.lexer.Unread(token)
	return value, def
}

// Parses maths operations, respecting the correct order of operations
func (p *Parser) ParseMathsOperations(value environment.Node, def TypeDef, onlyMultiplication bool) (environment.Node, TypeDef) {
	token := p.lexer.NextOrExit()
	operationType := token.Type
	if operationType != TokenPlus && operationType != TokenDash && operationType != TokenAsterisk && operationType != TokenForwardSlash {
		p.lexer.Unread(token)
		return value, def
	}
	if !def.IsNumber() {
		p.ThrowTypeError("Mathematical operations cannot be performed on values that don't represent a number.")
	}

	for {
		rhsVal, rhsDef := p.ParsePartialValue(def)
		if !rhsDef.Equals(def) {
			p.ThrowTypeError("Mathematical operations must be performed on values of the same type.")
		}
		if onlyMultiplication && operationType != TokenAsterisk && operationType != TokenForwardSlash {
			return value, def
		}

		var operation nodes.MathsOperationType
		switch operationType {
		case TokenAsterisk:
			operation = nodes.MathsMultiplication
		case TokenForwardSlash:
			operation = nodes.MathsDivision
		case TokenPlus:
			operation = nodes.MathsAddition
		case TokenDash:
			operation = nodes.MathsSubtraction
		default:
			panic("Non-operation token passed as operationType token")
		}

		// Read the next operation
		token := p.lexer.NextOrExit()
		if token.Type == TokenAsterisk || token.Type == TokenForwardSlash {
			p.lexer.Unread(token)
			rhsVal, _ = p.ParseMathsOperations(rhsVal, def, true)
			token = p.lexer.NextOrExit()
		}
		value = GetGenericTypeNode(def).GetMathsOperation(operation, value, rhsVal)
		if token.Type != TokenAsterisk && token.Type != TokenForwardSlash && token.Type != TokenPlus && token.Type != TokenDash {
			p.lexer.Unread(token)
			return value, def
		}
		operationType = token.Type
	}
}

func (p *Parser) ParseOperator(value environment.Node, def TypeDef) (environment.Node, TypeDef) {
	token := p.lexer.NextOrExit()

	switch token.Type {
	case TokenAmpersand:
		p.ExpectToken(TokenAmpersand)

		if !def.Equals(GenericTypeDef{TypeBool}) {
			p.ThrowTypeError("Operation && can only be used on boolean value.")
		}

		rhsVal, rhsValDef := p.ParseValue(nil)
		if !rhsValDef.Equals(GenericTypeDef{TypeBool}) {
			p.ThrowTypeError("Right hand side of && operation must be a boolean value.")
		}

		return &nodes.And{
			LeftSide:  value,
			RightSide: rhsVal,
		}, GenericTypeDef{TypeBool}

	case TokenBar:
		p.ExpectToken(TokenBar)

		if !def.Equals(GenericTypeDef{TypeBool}) {
			p.ThrowTypeError("Operation || can only be used on boolean value.")
		}

		rhsVal, rhsValDef := p.ParseValue(nil)
		if !rhsValDef.Equals(GenericTypeDef{TypeBool}) {
			p.ThrowTypeError("Right hand side of || operation must be a boolean value.")
		}

		return &nodes.Or{
			LeftSide:  value,
			RightSide: rhsVal,
		}, GenericTypeDef{TypeBool}

	case TokenEquals:
		// Check for comparison
		if p.lexer.PeekOrExit().Type == TokenEquals {
			p.lexer.Next()
			rhsVal, rhsValDef := p.ParseCalculatedValue(def)
			if !rhsValDef.Equals(def) {
				p.ThrowTypeError("Right hand side of comparison must be the same type as the left hand side.")
			}
			return p.ParseOperator(&nodes.EqualityComparison{LeftSide: value, RightSide: rhsVal}, GenericTypeDef{TypeBool})
		}

		// Check for assignment
		if ident, ok := value.(*nodes.Identifier); ok {
			newVal, newValDef := p.ParseValue(def)
			if !def.Equals(newValDef) {
				p.ThrowTypeError("Cannot assign new type to variable \"", ident.Name, "\".")
			}

			_, depth := p.currentTypeEnv.Get(ident.Name)

			return &nodes.Assignment{
				Identifier: ident.Name,
				NewValue:   newVal,
				Depth:      depth,
			}, def
		} else if arrayNode, indexNode, ok := GetGenericTypeNode(def).ArrayIndexDetails(value); ok {
			// Assignment to element of array
			newVal, newValDef := p.ParseValue(def)
			if !newValDef.Equals(def) {
				p.ThrowTypeError("Incorrect type in array element assignment.")
			}
			return GetGenericTypeNode(def).GetArrayAssignment(arrayNode, indexNode, newVal), def
		} else {
			p.ThrowSyntaxError("Left hand side of assignment is not assignable.")
		}

	case TokenGreaterThan, TokenLessThan:
		if !def.IsNumber() {
			p.ThrowTypeError("Cannot perform comparison on non-number.")
		}
		nextToken := p.lexer.PeekOrExit()
		var comparison nodes.ComparisonType
		if token.Type == TokenGreaterThan && nextToken.Type == TokenEquals {
			p.lexer.Next()
			comparison = nodes.ComparisonGreaterThanOrEquals
		} else if token.Type == TokenGreaterThan {
			comparison = nodes.ComparisonGreaterThan
		} else if nextToken.Type == TokenEquals {
			p.lexer.Next()
			comparison = nodes.ComparisonLessThanOrEquals
		} else {
			comparison = nodes.ComparisonLessThan
		}
		rhsVal, rhsValDef := p.ParseCalculatedValue(def)
		if !rhsValDef.Equals(def) {
			p.ThrowTypeError("Right hand side of comparison must be the same type as the left hand side.")
		}
		return p.ParseOperator(GetGenericTypeNode(def).GetInequalityComparison(comparison, value, rhsVal), GenericTypeDef{TypeBool})

	case TokenExclamationMark:
		p.ExpectToken(TokenEquals)
		rhsVal, rhsValDef := p.ParseValue(def)
		if !rhsValDef.Equals(def) {
			p.ThrowTypeError("Right hand side of comparison must be the same type as the left hand side.")
		}
		return p.ParseOperator(&nodes.Not{
			Value: &nodes.EqualityComparison{LeftSide: value, RightSide: rhsVal},
		}, GenericTypeDef{TypeBool})
	}

	p.lexer.Unread(token)
	return value, def
}

// Parses a value of any type, without accounting for logical operations that follow it.
func (p *Parser) ParsePartialValue(implicitType TypeDef) (environment.Node, TypeDef) {
	token := p.ExpectToken(TokenString, TokenNumber, TokenIdentifier, TokenTrue, TokenFalse, TokenDash, TokenLeftBracket, TokenLeftSquareBracket, TokenNewLine, TokenExclamationMark)
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
			if implicitType != nil {
				if implicitVal := ConvertFloat64ToTypeDef(val, implicitType.GetGenericType()); implicitVal != nil {
					return &nodes.Value{Value: implicitVal}, implicitType
				}
			}
			return p.ParseValueExpression(&nodes.Value{Value: val}, GenericTypeDef{TypeFloat64})
		}
		val, _ := strconv.ParseInt(token.Literal, 10, 64)

		if implicitType != nil {
			if implicitVal := ConvertInt64ToTypeDef(val, implicitType.GetGenericType()); implicitVal != nil {
				return p.ParseValueExpression(&nodes.Value{Value: implicitVal}, implicitType)
			}
		}
		return p.ParseValueExpression(&nodes.Value{Value: val}, GenericTypeDef{TypeInt64})

	case TokenIdentifier:
		typeDef, _ := p.currentTypeEnv.Get(token.Literal)
		if typeDef == nil {
			p.ThrowTypeError(token.Literal, " is not defined in this scope.")
		}
		return p.ParseValueExpression(&nodes.Identifier{Name: token.Literal}, typeDef)

	case TokenLeftBracket:
		defer p.ExpectToken(TokenRightBracket)
		return p.ParseValue(implicitType)

	case TokenExclamationMark:
		val, def := p.ParseValue(nil)
		if def.GetGenericType() != TypeBool {
			p.ThrowTypeError("Not operator must be used on a boolean value.")
		}
		return &nodes.Not{Value: val}, GenericTypeDef{TypeBool}

	case TokenLeftSquareBracket:
		var elements []environment.Node
		var elementType TypeDef
		size := -1
		if implicitType != nil && implicitType.GetGenericType() == TypeArray {
			size = implicitType.(ArrayDef).Size
			if size == -1 {
				elements = make([]environment.Node, 0)
			} else {
				elements = make([]environment.Node, size)
			}
			elementType = implicitType.(ArrayDef).ElementType
		} else {
			elements = make([]environment.Node, 0)
		}
		for position := 0; ; position++ {
			if p.lexer.PeekOrExit().Type == TokenRightBracket {
				p.lexer.Next()
				break
			}
			element, def := p.ParseValue(elementType)
			if elementType == nil {
				elementType = def
			} else if !def.Equals(elementType) {
				p.ThrowTypeError("Incorrect type for element of array.")
			}
			if size != -1 && position == size {
				p.ThrowTypeError("Maximum number of elements in array reached.")
			} else if size == -1 {
				// If the array is an unkown size we need to append to it
				elements = append(elements, element)
			} else {
				// If the array is a fixed size, all of the elements should already have been created
				elements[position] = element
			}

			if token = p.ExpectToken(TokenComma, TokenRightSquareBracket); token.Type == TokenRightSquareBracket {
				break
			}
		}
		if elementType == nil {
			// Since if the program has no specified type for the array and there are no elements to implicitly
			// get the type from, an array with no elements and no explicit type definition is not allowed.
			p.ThrowTypeError("An array of an unkown type cannot have 0 elements")
		}
		return GetGenericTypeNode(elementType).GetArrayInitialization(elements), ArrayDef{GenericTypeDef: GenericTypeDef{TypeArray}, ElementType: elementType, Size: size}
	}
	return p.ParseValue(implicitType)
}

func (p *Parser) ParseCalculatedValue(implicitType TypeDef) (environment.Node, TypeDef) {
	val, def := p.ParsePartialValue(implicitType)
	return p.ParseMathsOperations(val, def, false)
}

// Parses a value of any type, accounting for operations that follow it.
//
// If implicitType is passed, the value will be coerced to the implicit type if possible.
func (p *Parser) ParseValue(implicitType TypeDef) (environment.Node, TypeDef) {
	return p.ParseOperator(p.ParseCalculatedValue(implicitType))
}
