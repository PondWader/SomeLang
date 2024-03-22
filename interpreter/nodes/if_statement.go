package nodes

import "main/interpreter/environment"

// Node that runs an inner block if Condition is true, if Else is set that will be run if Condition is not true
type IfStatement struct {
	Condition environment.Node
	Inner     *Block
	Else      environment.Node
}

func (n *IfStatement) Eval(env *environment.Environment) any {
	if n.Condition.Eval(env).(bool) {
		childEnv := env.NewChild(environment.Call{})
		n.Inner.Eval(childEnv)
	} else if n.Else != nil {
		n.Else.Eval(env)
	}
	return nil
}

func (n *IfStatement) References() []string {
	refs := append(n.Condition.References(), n.Inner.References()...)
	if n.Else != nil {
		return append(refs, n.Else.References()...)
	}
	return refs
}
