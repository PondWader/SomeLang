package nodes

import (
	"main/interpreter/environment"
	"reflect"
)

type StructProperty struct {
	Struct   environment.Node
	Index    int
	IsMethod bool
	Name     string
}

func (n *StructProperty) Eval(env *environment.Environment) any {
	instance := n.Struct.Eval(env).([]any)
	val := instance[n.Index]
	// If the value is a method, a proxy function is used to set the instance as the first argument (self arg)
	if n.IsMethod {
		return func(argVals ...any) any {
			function := reflect.ValueOf(val)
			args := make([]reflect.Value, len(argVals)+1)
			args[0] = reflect.ValueOf(instance)
			for i, arg := range argVals {
				args[i+1] = reflect.ValueOf(arg)
			}
			out := function.Call(args)
			if len(out) > 0 {
				return out[0].Interface()
			}
			return nil
		}
	}
	return val
}

func (n *StructProperty) References() []string {
	return n.Struct.References()
}
