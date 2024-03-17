package interpreter

import (
	"main/interpreter/environment"
)

func Execute(ast []environment.Node, fileName string, runProfiler bool, globals map[string]any, modules map[string]map[string]any) {
	env := environment.New(nil, environment.Call{
		File:         fileName,
		Line:         0,
		FunctionName: "main",
	}, modules)

	for name, val := range globals {
		env.Set(name, val)
	}

	env.Execute(ast)
}
