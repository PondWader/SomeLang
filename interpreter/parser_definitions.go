package interpreter

import "strconv"

func (p *Parser) ParseFunctionDef() (name string, argDefs []TypeDef, argNames []string, returnType TypeDef) {
	name = p.ExpectToken(TokenIdentifier).Literal
	p.ExpectToken(TokenLeftBracket)

	argDefs = make([]TypeDef, 0)
	argNames = make([]string, 0)
	for i := 0; ; i++ {
		token := p.ExpectToken(TokenIdentifier, TokenRightBracket)
		if token.Type == TokenRightBracket {
			break
		}

		p.ExpectToken(TokenColon)
		argDefs = append(argDefs, p.ParseTypeDef())
		argNames = append(argNames, token.Literal)

		token = p.ExpectToken(TokenComma, TokenRightBracket)
		if token.Type == TokenRightBracket {
			break
		}
	}

	token := p.lexer.NextOrExit()
	if token.Type == TokenColon {
		returnType = p.ParseTypeDef()
	} else {
		p.lexer.Unread(token)
	}
	return
}

func (p *Parser) ParseTypeDef() TypeDef {
	// Expect a token of a type
	token := p.ExpectToken(TokenTypeInt8, TokenTypeInt16, TokenTypeInt32, TokenTypeInt64, TokenTypeUint8, TokenTypeUint16, TokenTypeUint32, TokenTypeUint64, TokenTypeFloat32, TokenTypeFloat64, TokenTypeString, TokenTypeBool, TokenTypeMap, TokenLeftSquareBracket, TokenFunctionDeclaration)

	var typeDef TypeDef
	switch token.Type {
	case TokenFunctionDeclaration:
		_, argDefs, _, returnType := p.ParseFunctionDef()
		typeDef = FuncDef{
			GenericTypeDef{TypeFunc},
			argDefs,
			false,
			returnType,
		}

	case TokenTypeMap:
		p.ExpectToken(TokenLeftSquareBracket)
		keyType := p.ParseTypeDef()
		p.ExpectToken(TokenRightSquareBracket)
		valueType := p.ParseTypeDef()

		typeDef = MapDef{
			GenericTypeDef{TypeMap},
			keyType,
			valueType,
		}

	case TokenLeftSquareBracket:
		token = p.ExpectToken(TokenRightSquareBracket, TokenNumber)
		size := -1
		if token.Type == TokenNumber {
			size, _ = strconv.Atoi(token.Literal)
			if size < 0 {
				p.ThrowSyntaxError("Size of array must be greater than or equal to 0")
			}
			p.ExpectToken(TokenRightSquareBracket)
		}
		typeDef = ArrayDef{
			GenericTypeDef: GenericTypeDef{TypeArray},
			ElementType:    p.ParseTypeDef(),
			Size:           size,
		}

	default:
		typeDef = GenericTypeDef{
			Type: TypeTokenToPrimitiveType(token),
		}
	}

	return typeDef
}

// A function that converts a token for a type (such as "int8") to it's corresponding type code
func TypeTokenToPrimitiveType(token Token) GenericType {
	// Uses the offset of the first token type, based on the assumption they are in the same order as the types
	return GenericType(token.Type - TokenTypeInt8)
}
