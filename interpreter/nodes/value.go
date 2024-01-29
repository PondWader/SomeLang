package nodes

import (
	"main/interpreter/environment"
	"reflect"
)

type Value struct {
	Value any
}

func (v *Value) Eval(env *environment.Environment) any {
  return v.Value
}

func (v *Value) Type(env *environment.Environment) string {
  return reflect.TypeOf(v.Value).Name()
}