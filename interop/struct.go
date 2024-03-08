package interop

import (
	"reflect"
)

// Converts go structs to structs for use in the interpreter

func CreateRuntimeStruct(structVal any) []any {
	val := reflect.ValueOf(structVal)
	runtimeStruct := []any{}

	for i := 0; i < val.NumMethod(); i++ {
		method := val.Method(i)
		runtimeStruct = append(runtimeStruct, func(args ...any) {
			argVals := make([]reflect.Value, len(args))
			for i, arg := range args {
				argVals[i] = reflect.ValueOf(arg)
			}
			method.Call(argVals)
		})
	}
	return runtimeStruct
}
