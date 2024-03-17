package interop

import (
	"reflect"
)

// Converts go structs to structs for use in the interpreter

func CreateRuntimeStruct(structVal any, methodOrder []string) []any {
	val := reflect.ValueOf(structVal)
	runtimeStruct := make([]any, val.NumMethod())

	for i, methodName := range methodOrder {
		method := val.MethodByName(methodName)
		if !reflect.Value.IsValid(method) {
			panic(methodName + " is not a valid method on the struct")
		}

		runtimeStruct[i] = func(args ...any) any {
			argVals := make([]reflect.Value, len(args)-1)
			for i, arg := range args {
				if i == 0 {
					continue
				}
				argVals[i-1] = reflect.ValueOf(arg)
			}
			out := method.Call(argVals)
			if len(out) > 0 {
				return out[0].Interface()
			}
			return nil
		}
	}
	return runtimeStruct
}
