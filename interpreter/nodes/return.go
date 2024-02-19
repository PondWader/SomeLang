package nodes

import "main/interpreter/environment"

type Return struct {
	Value environment.Node
}

func (r *Return) Eval(env *environment.Environment) any {
	env.Return(r.Value.Eval(env))
	return nil
}

func (r *Return) References() []string {
	return r.Value.References()
}
