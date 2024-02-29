package nodes

import (
	"main/interpreter/environment"
)

type ArrayAssignment[Element any] struct {
	ArrayIndex *ArrayIndex[Element]
	Value      environment.Node
}

func (n *ArrayAssignment[E]) Eval(env *environment.Environment) any {
	newVal := n.Value.Eval(env).(E)
	array, index := n.ArrayIndex.GetArrayAndValidatedIndex(env)
	array[index] = newVal
	return newVal
}

func (n *ArrayAssignment[E]) References() []string {
	return append(n.ArrayIndex.References(), n.Value.References()...)
}
