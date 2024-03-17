package nodes

import "main/interpreter/environment"

type MapAccess[KeyType comparable, ValueType any] struct {
	Identifier string
	Key        KeyType
}

func (n *MapAccess[KeyType, ValueType]) Eval(env *environment.Environment) any {
	return env.Get(n.Identifier).(map[KeyType]ValueType)[n.Key]
}

func (n *MapAccess[KeyType, ValueType]) References() []string {
	return []string{n.Identifier}
}
