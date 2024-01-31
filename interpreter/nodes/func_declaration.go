package nodes

import "main/interpreter/environment"

type FuncDeclaration struct {
	Name     string
  Line     int
	Inner    *Block
	ArgNames []string
}

func (fd *FuncDeclaration) Eval(env *environment.Environment) any {
  fn := func(args ...any) any {
    innerEnv := env.NewChild(environment.Call{
      FunctionName: fd.Name,
      File: env.Call.File,
      Line: fd.Line,
    })

    for i, arg := range args {
      innerEnv.Set(fd.ArgNames[i], arg)
    }

    var returnVal any
    innerEnv.Return = func(v any) {
      returnVal = v
    }

    fd.Inner.Eval(innerEnv)

    return returnVal
  }
  // Future proof for anonymous functions without a name
  if fd.Name != "" {
    env.Set(fd.Name, fn)
  }
  return fn
}