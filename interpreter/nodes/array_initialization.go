package nodes

import (
	"main/interpreter/environment"
)

type ArrayInitialization[T any] struct {
	Elements []environment.Node
}

func (n *ArrayInitialization[T]) Eval(env *environment.Environment) any {
	array := make([]T, len(n.Elements))
	for pos, el := range n.Elements {
		if el == nil {
			continue
		}
		array[pos] = el.Eval(env).(T)
	}
	return array
}

func (n *ArrayInitialization[T]) References() []string {
	refs := make([]string, 0)
	for _, el := range n.Elements {
		refs = append(refs, el.References()...)
	}
	return refs
}
