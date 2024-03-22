package nodes

import "main/interpreter/environment"

// Node that returns true if LeftSide and RightSide are both true
type And struct {
	LeftSide  environment.Node
	RightSide environment.Node
}

func (a *And) Eval(env *environment.Environment) any {
	return a.LeftSide.Eval(env).(bool) && a.RightSide.Eval(env).(bool)
}

func (a *And) References() []string {
	return append(a.LeftSide.References(), a.RightSide.References()...)
}
