package interpreter

import (
	"main/interpreter/environment"
	"main/interpreter/nodes"
)

// Can generate nodes that are type specific, such as maths operations or comparisons where the go type must be known
// ahead of time to perform the operation
type TypeNodeGenerator interface {
	GetMathsOperation(operation nodes.MathsOperationType, leftSide environment.Node, rightSide environment.Node) environment.Node
	GetInequalityComparison(comparison nodes.ComparisonType, leftSide environment.Node, rightSide environment.Node) environment.Node
}

func GetGenericTypeNode(def TypeDef) TypeNodeGenerator {
	genericType := def.GetGenericType()
	switch genericType {
	case TypeInt8:
		return TypeNodeGeneratorImpl[int8]{}
	case TypeInt16:
		return TypeNodeGeneratorImpl[int16]{}
	case TypeInt32:
		return TypeNodeGeneratorImpl[int32]{}
	case TypeInt64:
		return TypeNodeGeneratorImpl[int64]{}
	case TypeUint8:
		return TypeNodeGeneratorImpl[uint8]{}
	case TypeUint16:
		return TypeNodeGeneratorImpl[uint16]{}
	case TypeUint32:
		return TypeNodeGeneratorImpl[uint32]{}
	case TypeUint64:
		return TypeNodeGeneratorImpl[uint64]{}
	case TypeFloat32:
		return TypeNodeGeneratorImpl[float32]{}
	case TypeFloat64:
		return TypeNodeGeneratorImpl[float64]{}
	}
	return nil
}

// Implementation of TypeNodeGenerator with generic types
type TypeNodeGeneratorImpl[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float32 | float64] struct{}

func (tn TypeNodeGeneratorImpl[T]) GetMathsOperation(operation nodes.MathsOperationType, leftSide environment.Node, rightSide environment.Node) environment.Node {
	return &nodes.MathsOperation[T]{
		Operation: operation,
		LeftSide:  leftSide,
		RightSide: rightSide,
	}
}

func (tn TypeNodeGeneratorImpl[T]) GetInequalityComparison(comparison nodes.ComparisonType, leftSide environment.Node, rightSide environment.Node) environment.Node {
	return &nodes.InequalityComparison[T]{
		Type:      comparison,
		LeftSide:  leftSide,
		RightSide: rightSide,
	}
}
