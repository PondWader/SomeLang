package interpreter

import (
	"fmt"
	"main/interpreter/environment"
	"main/interpreter/nodes"
)

// Structs allow for custom data types to be defined
// Under the hood the interpreter converts these custom data types in to arrays of data for maximum performance
// This means all validation is done at the parsing step ahead of time
// When a struct is declared, at runtime a function is declare that accepts the slices properties as arguments.
// This then creates an array with the number of elements of the number of properties and methods.
// The parser has to match property and method names to the index that they are expected to be at during runtime

type structMethodDeclaration struct {
	name         string
	def          FuncDef
	argNames     []string
	args         map[string]TypeDef
	codeBlockPos LexerPos
}

func (p *Parser) ParseStructDeclaration() environment.Node {
	fmt.Println("Parsing struct")
	name := p.ExpectToken(TokenIdentifier).Literal
	fmt.Println("struct name:", name)

	p.ExpectToken(TokenLeftBrace)
	def := StructDef{
		GenericTypeDef: GenericTypeDef{TypeStruct},
		Properties:     make(map[string]int),
		PropertyDefs:   make([]TypeDef, 0),
	}

	methodDeclarations := make([]structMethodDeclaration, 0)

	i := 0
	for {
		// Expect comma between struct values
		if i != 0 {
			p.ExpectToken()
		}

		token := p.ExpectToken(TokenRightBrace, TokenIdentifier, TokenFunctionDeclaration, TokenNewLine)
		if token.Type == TokenNewLine {
			continue
		} else if token.Type == TokenRightBrace {
			break
		} else if token.Type == TokenIdentifier {
			p.ExpectToken(TokenColon)
			propertyDef := p.ParseTypeDef()
			def.Properties[token.Literal] = i
			def.PropertyDefs = append(def.PropertyDefs, propertyDef)
		} else {
			methodName, argDefs, argNames, returnType := p.ParseFunctionDef()

			funcDef := FuncDef{
				GenericTypeDef: GenericTypeDef{TypeFunc},
				Args:           argDefs,
				ReturnType:     returnType,
			}
			args := make(map[string]TypeDef, len(argDefs))
			for i, name := range argNames {
				args[name] = argDefs[i]
			}
			def.Properties[methodName] = i
			def.PropertyDefs = append(def.PropertyDefs, funcDef)
			methodDeclarations = append(methodDeclarations, structMethodDeclaration{
				name:         methodName,
				def:          funcDef,
				args:         args,
				argNames:     argNames,
				codeBlockPos: p.lexer.SavePos(),
			})
		}
		i++
	}

	methods := make([]environment.Node, len(methodDeclarations))
	for i, methodDeclaration := range methodDeclarations {
		revertPos := methodDeclaration.codeBlockPos.GoTo()

		methodDeclaration.args["self"] = def
		innerBlock := p.ParseBlock(methodDeclaration.args, methodDeclaration.def.ReturnType)
		methods[i] = &nodes.FuncDeclaration{
			Name:     methodDeclaration.name,
			Line:     methodDeclaration.codeBlockPos.Line,
			Inner:    innerBlock,
			ArgNames: methodDeclaration.argNames,
		}

		revertPos()
	}

	p.currentTypeEnv.Set(name, def)

	return &nodes.StructDeclaration{
		Name:    name,
		Methods: methods,
	}
}

func (p *Parser) ParseStructInitialization(name string, def StructDef) environment.Node {
	namedProperties := false
	unnamedProperties := false

	values := make([]environment.Node, len(def.PropertyDefs))

	for i := 0; ; i++ {
		val, _ := p.ParseValue(def.PropertyDefs[i])
		token := p.ExpectToken(TokenComma, TokenColon, TokenNewLine, TokenRightBrace)
		if token.Type == TokenColon {
			if unnamedProperties {
				p.ThrowSyntaxError("Cannot use mix of named and unnamed parameters")
			}

			namedProperties = true
			propertyNameNode, ok := val.(*nodes.Identifier)
			if !ok {
				p.ThrowTypeError("Left side of struct property must be an identifier")
			}
			val, valDef := p.ParseValue(def.PropertyDefs[i])
			propertyIndex := def.Properties[propertyNameNode.Name]
			if !valDef.Equals(def.PropertyDefs[propertyIndex]) {
				p.ThrowTypeError("Incorrect struct type for ", propertyNameNode.Name)
			}
			values[propertyIndex] = val
		} else if namedProperties {
			p.ThrowSyntaxError("Expected colon after identifier for named property struct initialization")
		} else if token.Type == TokenRightBrace {
			break
		} else {

		}
	}

	return &nodes.FuncCall{
		Args: values,
		Function: &nodes.Identifier{
			Name: name,
		},
	}
}
