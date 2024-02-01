package nodes

import "main/interpreter/environment"

type Return struct {
	Value Node
}

func (r *Return) Eval(env *environment.Environment) any {
	env.Return(r.Value.Eval(env))
	return nil
}
