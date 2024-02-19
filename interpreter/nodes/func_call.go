package nodes

import (
	"main/interpreter/environment"
	"reflect"
)

type FuncCall struct {
	Args     []environment.Node
	Function environment.Node
}

func (fc *FuncCall) Eval(env *environment.Environment) any {
	v := reflect.ValueOf(fc.Function.Eval(env))
	args := make([]reflect.Value, len(fc.Args))
	for i, arg := range fc.Args {
		args[i] = reflect.ValueOf(arg.Eval(env))
	}
	out := v.Call(args)
	if len(out) > 0 {
		return out[0].Interface()
	} else {
		return nil
	}
}

func (fc *FuncCall) References() []string {
	refs := fc.Function.References()
	for _, arg := range fc.Args {
		refs = append(refs, arg.References()...)
	}
	return refs
}
