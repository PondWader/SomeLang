package nodes

import "main/interpreter/environment"

type MapValue[KeyType comparable, ValueType any] struct {
	Map environment.Node
	Key environment.Node
}

func (n *MapValue[KeyType, ValueType]) Eval(env *environment.Environment) any {
	return n.Map.Eval(env).(map[KeyType]ValueType)[n.Key.Eval(env).(KeyType)]
}

func (n *MapValue[KeyType, ValueType]) References() []string {
	return n.Map.References()
}
