package nodes

import (
  "main/interpreter/environment"
  "reflect"
)

type FuncCall struct {
  Args []Node
  Function Node
}

func (fc *FuncCall) Eval(env *environment.Environment) any {
  v := reflect.ValueOf(fc.Function.Eval(env))
  args := make([]reflect.Value, len(fc.Args))
  for i, arg := range fc.Args {
    args[i] = reflect.ValueOf(arg.Eval(env))
  }
  out := v.Call(args)
  return out[0].Interface()
}
