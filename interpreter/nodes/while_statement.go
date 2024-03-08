package nodes

import (
	"main/interpreter/environment"
)

type WhileStatement struct {
	Condition environment.Node
	Inner     *Block
}

func (n *WhileStatement) Eval(env *environment.Environment) any {
	for n.Condition.Eval(env).(bool) && !env.IsBroken {
		childEnv := env.NewChild(environment.Call{})
		n.Inner.Eval(childEnv)
	}
	return nil
}

func (n *WhileStatement) References() []string {
	return append(n.Condition.References(), n.Inner.References()...)
}
