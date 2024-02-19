package nodes

import "main/interpreter/environment"

type IfStatement struct {
	Condition environment.Node
	Inner     *Block
	Else      environment.Node
}

func (is *IfStatement) Eval(env *environment.Environment) any {
	if is.Condition.Eval(env).(bool) {
		env.NewChild(environment.Call{})
		is.Inner.Eval(env)
	} else if is.Else != nil {
		is.Else.Eval(env)
	}
	return nil
}

func (is *IfStatement) References() []string {
	refs := append(is.Condition.References(), is.Inner.References()...)
	if is.Else != nil {
		return append(refs, is.Else.References()...)
	}
	return refs
}
