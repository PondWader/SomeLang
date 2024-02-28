package nodes

import "main/interpreter/environment"

type ArrayIndex[T any] struct {
	Array environment.Node
	Index environment.Node
}

func (n *ArrayIndex[T]) Eval(env *environment.Environment) any {
	return n.Array.Eval(env).([]T)[n.Index.Eval(env).(int64)]
}

func (n *ArrayIndex[T]) References() []string {
	return append(n.Array.References(), n.Index.References()...)
}
