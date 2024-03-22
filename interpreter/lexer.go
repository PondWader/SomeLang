package interpreter

import (
	"errors"
	"fmt"
	"os"
)

type TokenType uint8

const (
	// Statements
	TokenIfStatement TokenType = iota
	TokenElseStatement
	TokenFunctionDeclaration
	TokenImportStatement
	TokenExportStatement
	TokenForStatement
	TokenVarDeclaration
	TokenReturnStatement
	TokenStructDeclaration
	TokenAsStatement
	TokenRangeStatement
	TokenWhileStatement

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

type Token struct {
	Type    TokenType
	Literal string
	Line    int
}

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
		// Need to use a substring instead of an index to get a value of type string
		char := l.content[l.cursor : l.cursor+1]
		l.cursor++

		if char == " " || char == "\t" || char == "\r" {
			continue
		}

		// Check character is a valid token
		if token, err := getCharTokenType(char); err == nil {
			if token == TokenNewLine {
				l.currentLine++
			}
			return Token{
				Type:    token,
				Literal: char,
				Line:    l.currentLine,
			}, nil
		}

		// If the character is a quotation mark, it's the beginning of a string
		if char == "\"" {
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

		currentStr += char
		var endOfToken bool
		if l.cursor >= len(l.content) {
			endOfToken = true
		} else {
			nextChar := l.content[l.cursor : l.cursor+1]
			// Check if the next character terminates a token
			if nextChar == "" || nextChar == " " || nextChar == "\n" || nextChar == "\r" || nextChar == "\t" || nextChar == "\"" {
				endOfToken = true
			} else if _, err := getCharTokenType(nextChar); err == nil {
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

// Returns the contents of the next token without progressing the cursor
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

func (l *Lexer) SetCurrentLine(line int) {
	l.currentLine = line
}

func (l *Lexer) GetCursor() int {
	return l.cursor
}

func (l *Lexer) SetCursor(cursor int) {
	l.cursor = cursor
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

func (l *Lexer) SavePos() LexerPos {
	return LexerPos{l.cursor, l.currentLine, l}
}

// Stores the position of a lexer
type LexerPos struct {
	Cursor int
	Line   int
	lexer  *Lexer
}

func (pos LexerPos) GoTo() (undo func()) {
	originalLine := pos.lexer.GetCurrentLine()
	originalCursor := pos.lexer.GetCursor()

	pos.lexer.SetCurrentLine(pos.Line)
	pos.lexer.SetCurrentLine(pos.Cursor)

	return func() {
		pos.lexer.SetCurrentLine(originalLine)
		pos.lexer.SetCursor(originalCursor)
	}
}
