package nodes

import (
	"main/interpreter/environment"
)

// Node that iterates through an array, and runs inner for each iteration
type LoopArray[Element any] struct {
	ValIdentifier   string
	IndexIdentifier string
	Array           environment.Node
	Inner           *Block
}

func (n *LoopArray[Element]) Eval(env *environment.Environment) any {
	for index, value := range n.Array.Eval(env).([]Element) {
		childEnv := env.NewChild(environment.Call{})
		childEnv.Set(n.ValIdentifier, value)
		childEnv.Set(n.IndexIdentifier, int64(index))
		n.Inner.Eval(childEnv)
	}
	return nil
}

func (n *LoopArray[Element]) References() []string {
	return append(n.Array.References(), n.Inner.References()...)
}
