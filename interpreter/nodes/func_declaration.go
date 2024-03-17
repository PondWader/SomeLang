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

func (n *FuncDeclaration) Eval(env *environment.Environment) any {
	fn := func(args ...any) any {
		innerEnv := env.NewChild(environment.Call{
			Name: n.Name + "()",
			File: env.Call.File,
			Line: n.Line,
		})

		for i, arg := range args {
			innerEnv.Set(n.ArgNames[i], arg)
		}

		var returnVal any
		innerEnv.SetReturnCallback(func(v any) {
			returnVal = v
		})

		n.Inner.Eval(innerEnv)

		env.GetCurrentExecutionEnv().ProfileFunctionCall(innerEnv.GetProfileResult())
		return returnVal
	}
	// Check that the function is not an anonymous functions without a name
	if n.Name != "" {
		env.Set(n.Name, fn)
		env.AttachReferences(n.Name, n.Inner.References())
	}
	return fn
}

func (n *FuncDeclaration) References() []string {
	return n.Inner.References()
}
