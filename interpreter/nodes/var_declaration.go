package nodes

import (
	"main/interpreter/environment"
)

type VarDeclaration struct {
	Identifier string
	Value      environment.Node
}

func (n *VarDeclaration) Eval(env *environment.Environment) any {
	env.Set(n.Identifier, n.Value.Eval(env))
	return nil
}

func (n *VarDeclaration) References() []string {
	return append(n.Value.References(), n.Identifier)
}
