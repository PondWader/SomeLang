package nodes

import "main/interpreter/environment"

type Return struct {
	Value environment.Node
}

func (n *Return) Eval(env *environment.Environment) any {
	env.Return(n.Value.Eval(env))
	return nil
}

func (n *Return) References() []string {
	return n.Value.References()
}
