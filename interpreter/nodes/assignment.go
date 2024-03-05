package nodes

import (
	"main/interpreter/environment"
)

type Assignment struct {
	Identifier string
	NewValue   environment.Node
	Depth      int
}

func (a *Assignment) Eval(env *environment.Environment) any {
	newVal := a.NewValue.Eval(env)
	env.SetWithDepth(a.Identifier, newVal, a.Depth)
	return newVal
}

func (a *Assignment) References() []string {
	return append(a.NewValue.References(), a.Identifier)
}
