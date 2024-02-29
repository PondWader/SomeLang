package nodes

import (
	"main/interpreter/environment"
	"reflect"
)

type ArrayAssignment[Element any] struct {
	Array environment.Node
	Index environment.Node
	Value environment.Node
}

func (n *ArrayAssignment[E]) Eval(env *environment.Environment) any {
	newVal := n.Value.Eval(env).(E)
	indexVal := reflect.ValueOf(n.Index.Eval(env))
	indexKind := indexVal.Kind()
	var index uint64
	// Check if index is signed integer or unsigned
	if indexKind == reflect.Int || indexKind == reflect.Int16 || indexKind == reflect.Int32 || indexKind == reflect.Int64 {
		index = uint64(indexVal.Int())
	} else {
		index = indexVal.Uint()
	}
	n.Array.Eval(env).([]E)[index] = newVal
	return newVal
}

func (n *ArrayAssignment[E]) References() []string {
	return append(n.Array.References(), append(n.Value.References(), n.Index.References()...)...)
}
