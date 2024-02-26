package nodes

import "main/interpreter/environment"

type EqualityComparison struct {
	LeftSide  environment.Node
	RightSide environment.Node
}

func (c *EqualityComparison) Eval(env *environment.Environment) any {
	return c.LeftSide.Eval(env) == c.RightSide.Eval(env)
}

func (c *EqualityComparison) References() []string {
	return append(c.RightSide.References(), c.LeftSide.References()...)
}
