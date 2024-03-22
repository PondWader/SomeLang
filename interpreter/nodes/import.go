package nodes

import "main/interpreter/environment"

// Node that imports a module in to the current environment
type Import struct {
	Module     string
	Identifier string
}

func (n *Import) Eval(env *environment.Environment) any {
	env.Set(n.Identifier, env.GetBuiltInModule(n.Module))
	return nil
}

func (n *Import) References() []string {
	return []string{n.Identifier}
}
