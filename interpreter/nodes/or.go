package nodes

import "main/interpreter/environment"

// Node that returns true if either LeftSide or RightSide is true or both
type Or struct {
	LeftSide  environment.Node
	RightSide environment.Node
}

func (n *Or) Eval(env *environment.Environment) any {
	return n.LeftSide.Eval(env).(bool) || n.RightSide.Eval(env).(bool)
}

func (n *Or) References() []string {
	return append(n.LeftSide.References(), n.RightSide.References()...)
}
