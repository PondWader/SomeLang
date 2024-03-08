package nodes

import (
	"main/interpreter/environment"
)

type FuncDeclaration struct {
	Name     string
	Line     int
	Inner    *Block
	ArgNames []string
}

func (fd *FuncDeclaration) Eval(env *environment.Environment) any {
	fn := func(args ...any) any {
		innerEnv := env.NewChild(environment.Call{
			FunctionName: fd.Name + "()",
			File:         env.Call.File,
			Line:         fd.Line,
		})

		for i, arg := range args {
			innerEnv.Set(fd.ArgNames[i], arg)
		}

		var returnVal any
		innerEnv.SetReturnCallback(func(v any) {
			returnVal = v
		})

		fd.Inner.Eval(innerEnv)

		return returnVal
	}
	// Check that the function is not an anonymous functions without a name
	if fd.Name != "" {
		env.Set(fd.Name, fn)
		env.AttachReferences(fd.Name, fd.Inner.References())
	}
	return fn
}

func (fd *FuncDeclaration) References() []string {
	return fd.Inner.References()
}
