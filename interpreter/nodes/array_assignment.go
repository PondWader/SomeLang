package nodes

import (
	"fmt"
	"main/interpreter/environment"
)

type ArrayAssignment[T any] struct {
	Array environment.Node
	Index environment.Node
	Value environment.Node
}

func (n *ArrayAssignment[T]) Eval(env *environment.Environment) any {
	newVal := n.Value.Eval(env)
	fmt.Println(n.Array.Eval(env))
	n.Array.Eval(env).([]T)[n.Index.Eval(env).(int64)] = n.Value.Eval(env).(T)
	return newVal
}

func (n *ArrayAssignment[T]) References() []string {
	return append(n.Array.References(), append(n.Value.References(), n.Index.References()...)...)
}
