package nodes

import (
	"main/interpreter/environment"
)

// Node that represents a literal value
type Value struct {
	Value any
}

func (v *Value) Eval(env *environment.Environment) any {
	return v.Value
}

func (v *Value) References() []string {
	return []string{}
}
