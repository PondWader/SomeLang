package nodes

import (
	"main/interpreter/environment"
)

type LoopWhile struct {
	Condition environment.Node
	Inner     *Block
}

func (n *LoopWhile) Eval(env *environment.Environment) any {
	for !env.IsBroken && n.Condition.Eval(env).(bool) {
		childEnv := env.NewChild(environment.Call{})
		n.Inner.Eval(childEnv)
	}
	return nil
}

func (n *LoopWhile) References() []string {
	return append(n.Condition.References(), n.Inner.References()...)
}
