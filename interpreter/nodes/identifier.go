package nodes

import (
  "main/interpreter/environment"
  "reflect"
)

type Identifier struct {
  Name string
}

func (i *Identifier) Eval(env *environment.Environment) any {
  return env.Get(i.Name)
}

func (i *Identifier) Type(env *environment.Environment) string {
  return reflect.TypeOf(env.Get(i.Name)).Name()
}