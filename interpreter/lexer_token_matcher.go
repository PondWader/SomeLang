package interpreter

import "errors"

func getCharTokenType(char string) (TokenType, error) {
	switch char {
	case ":":
		return TokenColon, nil
	case ";":
		return TokenSemiColon, nil
	case "\n":
		return TokenNewLine, nil
	case ",":
		return TokenComma, nil
	case "(":
		return TokenLeftBracket, nil
	case ")":
		return TokenRightBracket, nil
	case "{":
		return TokenLeftBrace, nil
	case "}":
		return TokenRightBrace, nil
	case "[":
		return TokenLeftSquareBracket, nil
	case "]":
		return TokenRightSquareBracket, nil
	case "*":
		return TokenAsterisk, nil
	case "+":
		return TokenPlus, nil
	case "-":
		return TokenDash, nil
	case "/":
		return TokenForwardSlash, nil
	case "&":
		return TokenAmpersand, nil
	case "|":
		return TokenBar, nil
	case "'":
		return TokenApostrophe, nil
	case "!":
		return TokenExclamationMark, nil
	case "=":
		return TokenEquals, nil
	case ">":
		return TokenGreaterThan, nil
	case "<":
		return TokenLessThan, nil
	case ".":
		return TokenPeriod, nil
	}
	return 0, errors.New("char provided is not a valid token")
}

func getLiteralTokenType(literal string) TokenType {
	if literal == "true" {
		return TokenTrue
	} else if literal == "false" {
		return TokenFalse
	}

	switch literal {
	// Statements
	case "if":
		return TokenIfStatement
	case "else":
		return TokenElseStatement
	case "fn":
		return TokenFunctionDeclaration
	case "class":
		return TokenClassDeclaration
	case "import":
		return TokenImportStatement
	case "export":
		return TokenExportStatement
	case "for":
		return TokenForStatement
	case "var":
		return TokenVarDeclaration
	case "return":
		return TokenReturnStatement

	// Types
	case "int8":
		return TokenTypeInt8
	case "int16":
		return TokenTypeInt16
	case "int32":
		return TokenTypeInt32
	case "int64":
		return TokenTypeInt64
	case "uint8":
		return TokenTypeUint8
	case "uint16":
		return TokenTypeUint16
	case "uint32":
		return TokenTypeUint32
	case "uint64":
		return TokenTypeUint64
	case "string":
		return TokenTypeString
	case "bool":
		return TokenTypeBool
	}

	for _, char := range literal {
		// 0-9 are between char code 48 and 57
		if char < 48 || char > 57 {
			return TokenIdentifier
		}
	}

	return TokenNumber
}
