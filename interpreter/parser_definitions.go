package interpreter

func (p *Parser) ParseFunctionDef() (name string, args map[string]TypeDef, returnType TypeDef) {
	name = p.ExpectToken(TokenIdentifier).Literal
	p.ExpectToken(TokenLeftBracket)
	args = make(map[string]TypeDef, 0)
	for {
		token := p.ExpectToken(TokenIdentifier, TokenRightBracket)
		if token.Type == TokenRightBracket {
			break
		}

		p.ExpectToken(TokenColon)
		args[token.Literal] = p.ParseTypeDef()

		token = p.ExpectToken(TokenComma, TokenRightBracket)
		if token.Type == TokenRightBracket {
			break
		}
	}

	p.ExpectToken(TokenColon)
	returnType = p.ParseTypeDef()
	return
}

func (p *Parser) ParseTypeDef() TypeDef {
	// Expect a token of a type
	token := p.ExpectToken(TokenTypeInt8, TokenTypeInt16, TokenTypeInt32, TokenTypeInt48, TokenTypeInt64, TokenTypeUint8, TokenTypeUint16, TokenTypeUint32, TokenTypeUint48, TokenTypeUint64, TokenTypeFloat32, TokenTypeFloat64, TokenTypeString, TokenTypeBool, TokenTypeMap, TokenFunctionDeclarationStatement)

	var typeDef TypeDef
	switch token.Type {
	case TokenFunctionDeclarationStatement:
		_, args, returnType := p.ParseFunctionDef()
		typeDef = FuncDef{
			GenericTypeDef{TypeMap},
			args,
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
	}

	typeDef = GenericTypeDef{
		Type: TypeTokenToPrimitiveType(token),
	}

  // Check for array (type followed by [])
  token = p.lexer.PeekOrExit()
  if token.Type == TokenLeftSquareBracket {
    p.lexer.Next()
    p.ExpectToken(TokenRightSquareBracket)
    return ArrayDef{
      GenericTypeDef{TypeArray},
      typeDef,
    }
  }
  
  return typeDef
}

// A function that converts a token for a type (such as "int8") to it's corresponding type code
// There's not really any other way to do this
func TypeTokenToPrimitiveType(token Token) GenericType {
	switch token.Type {
	case TokenTypeInt8:
		return TypeInt8
	case TokenTypeInt16:
		return TypeInt16
	case TokenTypeInt32:
		return TypeInt32
	case TokenTypeInt48:
		return TypeInt48
	case TokenTypeInt64:
		return TypeUint64
	case TokenTypeUint8:
		return TypeUint8
	case TokenTypeUint16:
		return TypeInt16
	case TokenTypeUint32:
		return TypeInt32
	case TokenTypeUint48:
		return TypeInt48
	case TokenTypeUint64:
		return TypeInt64
	case TokenTypeFloat32:
		return TypeFloat32
	case TokenTypeFloat64:
		return TypeFloat64
	case TokenTypeBool:
		return TypeBool
	}

	return 0
}
