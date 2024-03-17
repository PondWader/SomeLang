package interpreter

import (
	"main/interpreter/environment"
	"main/profiler"
)

func Execute(ast []environment.Node, fileName string, runProfiler bool, globals map[string]any, modules map[string]map[string]any) *profiler.ProfileResult {
	env := environment.New(nil, environment.Call{
		File: fileName,
		Line: 0,
		Name: "main",
	}, modules, runProfiler)

	for name, val := range globals {
		env.Set(name, val)
	}

	env.Execute(ast)
	return env.GetProfileResult()
}
