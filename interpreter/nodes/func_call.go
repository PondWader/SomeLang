package nodes

import (
  "main/interpreter/environment"
  "reflect"
)

type FuncCall struct {
  Args []Node
  Function Node
}

func (fc FuncCall) Eval(env *environment.Environment) any {
  v := reflect.ValueOf(fc.Function.Eval(env))
  args := make([]reflect.Value, len(fc.Args))
  for i, arg := range fc.Args {
    args[i] = reflect.ValueOf(arg.Eval(env))
  }
  out := v.Call(args)
  return out[0].Interface()
}

func (fc *FuncCall) Type(env *environment.Environment) string {
  // TODO: Store return type of functions in some sort of metadata some where
  // or... could actually just not have return type in execution and have 
  // everything stored in seperate state by the interpreter to make sure types are correct
  // I think this aproach makes a lot more sense :)
  return ""
}