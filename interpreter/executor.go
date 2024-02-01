package interpreter

import (
	"main/interpreter/environment"
	"main/interpreter/nodes"
)

func Execute(ast []nodes.Node, fileName string, globals map[string]any) {
	env := environment.New(nil, environment.Call{
		File:         fileName,
		Line:         0,
		FunctionName: "main",
	})

	for name, val := range globals {
		env.Set(name, val)
	}

	for _, node := range ast {
		node.Eval(env)
	}
}
