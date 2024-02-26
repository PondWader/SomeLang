package nodes

import "main/interpreter/environment"

type Not struct {
	Value environment.Node
}

func (n *Not) Eval(env *environment.Environment) any {
	return !n.Value.Eval(env).(bool)
}

func (n *Not) References() []string {
	return n.Value.References()
}
