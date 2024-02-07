package nodes

import "main/interpreter/environment"

type IfStatement struct {
	Condition Node
	Inner     *Block
	Else      Node
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
