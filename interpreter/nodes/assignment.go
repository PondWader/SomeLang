package nodes

import "main/interpreter/environment"

type Assignment struct {
	Identifier string
	NewValue   Node
}

func (a *Assignment) Eval(env *environment.Environment) any {
  newVal := a.NewValue.Eval(env)
	env.Set(a.Identifier, newVal)
  return newVal 
}
