package nodes

import (
	"main/interpreter/environment"
)

// Node that assigns a value to an identifier in the current environment or in the parent environment with a certain depth
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
