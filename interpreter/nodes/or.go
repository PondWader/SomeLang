package nodes

import "main/interpreter/environment"

type Or struct {
	LeftSide  environment.Node
	RightSide environment.Node
}

func (o *Or) Eval(env *environment.Environment) any {
	return o.LeftSide.Eval(env).(bool) || o.RightSide.Eval(env).(bool)
}

func (o *Or) References() []string {
	return append(o.LeftSide.References(), o.RightSide.References()...)
}
