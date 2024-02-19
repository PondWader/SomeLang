package nodes

import "main/interpreter/environment"

type Or struct {
	LeftSide  Node
	RightSide Node
}

func (o *Or) Eval(env *environment.Environment) any {
	return o.LeftSide.Eval(env).(bool) || o.RightSide.Eval(env).(bool)
}
