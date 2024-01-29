package environment

import "strconv"

type Environment struct {
  identifiers map[string]any
  parent *Environment
  // call stores which function call initialized the environment
  Call Call
  // return can store a callback that returns a value in a function call
  Return func(any)
}

type Call struct {
  File string 
  Line int 
  FunctionName string
}

func New(parent *Environment, call Call) *Environment {
  return &Environment{
    identifiers: make(map[string]any),
    parent: parent,
    Call: call,
  }
}

func (e *Environment) Get(name string) any {
  if value, ok := e.identifiers[name]; ok {
    return value
  }
  if e.parent != nil {
    return e.parent.Get(name)
  }
  return nil
}

func (e *Environment) Set(name string, value any) {
  e.identifiers[name] = value
}

func (e *Environment) NewChild(call Call) *Environment {
  child := New(e, call)
  child.Return = e.Return
  return child
}

func (e *Environment) GetCallStackOutput() string {
  output := "File, " + e.Call.File + ", Line, " + strconv.Itoa(e.Call.Line) + ", In " + e.Call.FunctionName
  if e.parent != nil {
    output += "\n" + e.parent.GetCallStackOutput()
  }
  return output
}