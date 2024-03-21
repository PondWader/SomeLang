package nodes

import (
	"main/interpreter/environment"
	"reflect"
)

type LoopRange struct {
	ValIdentifier string
	Start         environment.Node
	End           environment.Node
	Inner         *Block
}

func (n *LoopRange) Eval(env *environment.Environment) any {
	startVal := getLoopRangeVal(n.Start, env)
	endVal := getLoopRangeVal(n.End, env)
	for i := startVal; i < endVal; i++ {
		childEnv := env.NewChild(environment.Call{})
		childEnv.Set(n.ValIdentifier, i)
		n.Inner.Eval(childEnv)
	}
	return nil
}

func (n *LoopRange) References() []string {
	return append(n.Start.References(), append(n.End.References(), n.Inner.References()...)...)
}

// Required to transform the node into int64 since go requires strong typing so can't use the evaluated node value as any
func getLoopRangeVal(node environment.Node, env *environment.Environment) int64 {
	val := reflect.ValueOf(node.Eval(env))
	kind := val.Kind()
	// Check if index is signed integer or unsigned
	if kind == reflect.Uint || kind == reflect.Uint16 || kind == reflect.Uint32 || kind == reflect.Uint64 {
		return int64(val.Uint())
	} else {
		return val.Int()
	}
}
