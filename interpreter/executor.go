package interpreter

import (
	"main/interpreter/environment"
)

func Execute(ast []environment.Node, fileName string, globals map[string]any) {
	env := environment.New(nil, ast, environment.Call{
		File:         fileName,
		Line:         0,
		FunctionName: "main",
	})

	for name, val := range globals {
		env.Set(name, val)
	}

	env.Execute()
}
