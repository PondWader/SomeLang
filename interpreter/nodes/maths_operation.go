package nodes

import (
	"fmt"
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
	LeftSide  Node
	RightSide Node
}

func (mo *MathsOperation[T]) Eval(env *environment.Environment) any {
	lhs, rhs := mo.LeftSide.Eval(env).(T), mo.RightSide.Eval(env).(T)
	switch mo.Operation {
	case MathsAddition:
		fmt.Println(lhs, "+", rhs)
		return lhs + rhs
	case MathsSubtraction:
		fmt.Println(lhs, "-", rhs)
		return lhs - rhs
	case MathsMultiplication:
		fmt.Println(lhs, "*", rhs)
		return lhs * rhs
	case MathsDivision:
		fmt.Println(lhs, "/", rhs)
		return lhs / rhs
	}
	return 0
}