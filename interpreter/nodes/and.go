package nodes

import "main/interpreter/environment"

type And struct {
	LeftSide  Node
	RightSide Node
}

func (a *And) Eval(env *environment.Environment) any {
	return a.LeftSide.Eval(env).(bool) && a.RightSide.Eval(env).(bool)
}
