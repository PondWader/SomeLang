package nodes

import "main/interpreter/environment"

type Assignment struct {
	Object Node
	Key    string
}

func (a *Assignment) Eval(env *environment.Environment) any {
	return ka.Object.Eval(env).(map[string]any)[ka.Key]
}
