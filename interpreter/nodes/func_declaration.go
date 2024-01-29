package nodes

import "main/interpreter/environment"

type FuncDeclaration struct {
	Name     string
  Line     int
	Inner    []Node
	ArgNames []string
}

func (fd *FuncDeclaration) Eval(env *environment.Environment) any {
  fn := func(args ...any) any {
    innerEnv := env.NewChild(environment.Call{
      FunctionName: fd.Name,
      File: env.Call.File,
      Line: fd.Line,
    })

    var returnVal any
    innerEnv.Return = func(v any) {
      returnVal = v
    }

    return returnVal
  }
  // Future proof for anonymous functions without a name
  if fd.Name != "" {
    env.Set(fd.Name, fn)
  }
  return fn
}