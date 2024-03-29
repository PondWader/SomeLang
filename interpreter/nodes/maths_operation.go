package nodes

import (
	"main/interpreter/environment"
)

type MathsOperationType uint8

const (
	MathsAddition MathsOperationType = iota
	MathsSubtraction
	MathsMultiplication
	MathsDivision
)

// Node that performs a maths operation on a value
type MathsOperation[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float32 | float64] struct {
	Operation MathsOperationType
	LeftSide  environment.Node
	RightSide environment.Node
}

func (n *MathsOperation[T]) Eval(env *environment.Environment) any {
	lhs, rhs := n.LeftSide.Eval(env).(T), n.RightSide.Eval(env).(T)
	switch n.Operation {
	case MathsAddition:
		return lhs + rhs
	case MathsSubtraction:
		return lhs - rhs
	case MathsMultiplication:
		return lhs * rhs
	case MathsDivision:
		return lhs / rhs
	}
	return 0
}

func (n *MathsOperation[T]) References() []string {
	return append(n.LeftSide.References(), n.RightSide.References()...)
}
