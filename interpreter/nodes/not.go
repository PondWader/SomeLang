package nodes

import "main/interpreter/environment"

// Node that returns opposite of a boolean value
type Not struct {
	Value environment.Node
}

func (n *Not) Eval(env *environment.Environment) any {
	return !n.Value.Eval(env).(bool)
}

func (n *Not) References() []string {
	return n.Value.References()
}
