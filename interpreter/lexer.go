package interpreter

import (
	"errors"
	"fmt"
	"os"
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
}

type TokenType uint8

const (
	// Statements
	TokenIfStatement TokenType = iota
	TokenElseStatement
	TokenFunctionDeclaration
	TokenClassDeclaration
	TokenImportStatement
	TokenExportStatement
	TokenForStatement
	TokenVarDeclaration
	TokenReturnStatement

	// Values
	TokenTrue
	TokenFalse
	TokenString
	TokenNumber
	TokenIdentifier

	// Types
	TokenTypeInt8
	TokenTypeInt16
	TokenTypeInt32
	TokenTypeInt64
	TokenTypeUint8
	TokenTypeUint16
	TokenTypeUint32
	TokenTypeUint64
	TokenTypeFloat32
	TokenTypeFloat64
	TokenTypeString
	TokenTypeBool
	TokenTypeMap

	// Symbols
	TokenColon
	TokenSemiColon
	TokenNewLine
	TokenComma
	TokenLeftBracket
	TokenRightBracket
	TokenLeftBrace
	TokenRightBrace
	TokenLeftSquareBracket
	TokenRightSquareBracket
	TokenAsterisk
	TokenPlus
	TokenDash
	TokenForwardSlash
	TokenAmpersand
	TokenBar
	TokenApostrophe
	TokenExclamationMark
	TokenEquals
	TokenGreaterThan
	TokenLessThan
	TokenPeriod

	TokenEOF
)

type Lexer struct {
	content     string
	cursor      int
	currentLine int
}

func NewLexer(content string) *Lexer {
	return &Lexer{content, 0, 1}
}

func (l *Lexer) Next() (Token, error) {
	currentStr := ""
	for l.cursor < len(l.content) {
		ch := l.content[l.cursor : l.cursor+1]
		l.cursor++

		if ch == " " || ch == "	" {
			continue
		}

		if token, err := getCharTokenType(ch); err == nil {
			if token == TokenNewLine {
				l.currentLine++
			}
			return Token{
				Type:    token,
				Literal: ch,
				Line:    l.currentLine,
			}, nil
		}

		if ch == "\"" {
			strContent, err := l.readString()
			if err != nil {
				return Token{}, err
			}
			return Token{
				Type:    TokenString,
				Literal: strContent,
				Line:    l.currentLine,
			}, nil
		}

		currentStr += ch
		var endOfToken bool
		if l.cursor >= len(l.content) {
			endOfToken = true
		} else {
			nextCh := l.content[l.cursor : l.cursor+1]
			if nextCh == "" || nextCh == " " || nextCh == "\n" || nextCh == "	" || nextCh == "\"" {
				endOfToken = true
			} else if _, err := getCharTokenType(nextCh); err == nil {
				endOfToken = true
			}
		}

		if endOfToken {
			return Token{
				Type:    getLiteralTokenType(currentStr),
				Literal: currentStr,
				Line:    l.currentLine,
			}, nil
		}
	}

	return Token{
		Type:    TokenEOF,
		Literal: "EOF",
		Line:    l.currentLine,
	}, nil
}

func (l *Lexer) Peek() (Token, error) {
	originalPos := l.cursor
	originalLine := l.currentLine
	token, err := l.Next()
	l.cursor = originalPos
	l.currentLine = originalLine
	return token, err
}

func (l *Lexer) PeekOrExit() Token {
	token, err := l.Peek()
	if err != nil {
		fmt.Println(err, fmt.Sprint(l.currentLine)+":"+fmt.Sprint(l.cursor))
		os.Exit(1)
	}
	return token
}

func (l *Lexer) NextOrExit() Token {
	token, err := l.Next()
	if err != nil {
		fmt.Println(err, fmt.Sprint(l.currentLine)+":"+fmt.Sprint(l.cursor))
		os.Exit(1)
	}
	return token
}

func (l *Lexer) readString() (string, error) {
	currentStr := ""
	escapedChar := false
	for l.cursor < len(l.content) {
		char := l.content[l.cursor : l.cursor+1]
		l.cursor++

		if char == "\\" && !escapedChar {
			escapedChar = true
			continue
		}
		if escapedChar {
			switch char {
			case "\"":
				return currentStr, nil
			case "n":
				currentStr += "\n"
			case "t":
				currentStr += "\t"
			case "r":
				currentStr += "\r"
			default:
				currentStr += char
			}
			escapedChar = false
		} else {
			if char == "\"" {
				return currentStr, nil
			}
			currentStr += char
		}
		if char == "\n" {
			return "", errors.New("unexpected newline while reading string literal")
		}
	}

	return "", errors.New("reached EOF without string finishing")
}

func (l *Lexer) GetCurrentLine() int {
	return l.currentLine
}

// Moves the cursor back to the start of the previously read token so it will be read at the next call of Next().
// Only the last read token should be passed to Unread.
func (l *Lexer) Unread(token Token) {
	if token.Type == TokenEOF {
		return
	}
	l.cursor -= len(token.Literal)
	if token.Type == TokenString {
		l.cursor -= 2 // Account for quotation marks on either side
	} else if token.Type == TokenNewLine {
		l.currentLine--
	}
}
