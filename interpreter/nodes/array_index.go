package nodes

import (
	"main/interpreter/environment"
	"reflect"
)

// Node that accesses a value at an array index
type ArrayIndex[Element any] struct {
	Array environment.Node
	Index environment.Node
}

func (n *ArrayIndex[E]) Eval(env *environment.Environment) any {
	array, index := n.GetArrayAndValidatedIndex(env)
	return array[index]
}

func (n *ArrayIndex[E]) References() []string {
	return append(n.Array.References(), n.Index.References()...)
}

func (n *ArrayIndex[E]) GetArrayAndValidatedIndex(env *environment.Environment) ([]E, uint64) {
	index := n.GetIndexVal(env)
	array := n.Array.Eval(env).([]E)
	if index > uint64(len(array)) {
		env.Panic("Index out of array bounds")
	}
	return array, index
}

// Required since Go generics are being used to ensure a valid index is returned
func (n *ArrayIndex[T]) GetIndexVal(env *environment.Environment) uint64 {
	indexVal := reflect.ValueOf(n.Index.Eval(env))
	indexKind := indexVal.Kind()
	// Check if index is signed integer or unsigned
	if indexKind == reflect.Int || indexKind == reflect.Int16 || indexKind == reflect.Int32 || indexKind == reflect.Int64 {
		index := indexVal.Int()
		if index < 0 {
			env.Panic("Array index cannot be less than 0")
		}
		return uint64(index)
	} else {
		return indexVal.Uint()
	}
}
