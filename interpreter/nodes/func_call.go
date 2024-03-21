package nodes

import (
	"main/interpreter/environment"
	"reflect"
)

type FuncCall struct {
	Args     []environment.Node
	Function environment.Node
}

func (n *FuncCall) Eval(env *environment.Environment) any {
	funcVal := reflect.ValueOf(n.Function.Eval(env))
	args := make([]reflect.Value, len(n.Args))
	for i, arg := range n.Args {
		args[i] = reflect.ValueOf(arg.Eval(env))
	}
	out := funcVal.Call(args)
	if len(out) > 0 {
		return out[0].Interface()
	}
	return nil
}

func (n *FuncCall) References() []string {
	refs := n.Function.References()
	for _, arg := range n.Args {
		refs = append(refs, arg.References()...)
	}
	return refs
}
