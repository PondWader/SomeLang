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

type MathsOperation[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float32 | float64] struct {
	Operation MathsOperationType
	LeftSide  environment.Node
	RightSide environment.Node
}

func (mo *MathsOperation[T]) Eval(env *environment.Environment) any {
	lhs, rhs := mo.LeftSide.Eval(env).(T), mo.RightSide.Eval(env).(T)
	switch mo.Operation {
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

func (mo *MathsOperation[T]) References() []string {
	return append(mo.LeftSide.References(), mo.RightSide.References()...)
}
