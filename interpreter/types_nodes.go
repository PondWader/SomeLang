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
	GetArrayInitialization(elements []environment.Node) environment.Node
	GetArrayIndex(array environment.Node, index environment.Node) environment.Node
	GetArrayAssignment(array environment.Node, index environment.Node, value environment.Node) environment.Node
	GetLoopArray(valIdentifier string, indexIdentifier string, array environment.Node, inner *nodes.Block) environment.Node
	ArrayIndexDetails(node environment.Node) (array environment.Node, index environment.Node, ok bool)
}

func GetGenericTypeNode(def TypeDef) TypeNodeGenerator {
	genericType := def.GetGenericType()
	switch genericType {
	case TypeString:
		return TypeNodeGeneratorAny[string]{}
	case TypeBool:
		return TypeNodeGeneratorAny[bool]{}
	case TypeInt8:
		return TypeNodeGeneratorNumber[int8]{}
	case TypeInt16:
		return TypeNodeGeneratorNumber[int16]{}
	case TypeInt32:
		return TypeNodeGeneratorNumber[int32]{}
	case TypeInt64:
		return TypeNodeGeneratorNumber[int64]{}
	case TypeUint8:
		return TypeNodeGeneratorNumber[uint8]{}
	case TypeUint16:
		return TypeNodeGeneratorNumber[uint16]{}
	case TypeUint32:
		return TypeNodeGeneratorNumber[uint32]{}
	case TypeUint64:
		return TypeNodeGeneratorNumber[uint64]{}
	case TypeFloat32:
		return TypeNodeGeneratorNumber[float32]{}
	case TypeFloat64:
		return TypeNodeGeneratorNumber[float64]{}
	}
	return nil
}

type TypeNodeGeneratorAny[T any] struct{}

func (tn TypeNodeGeneratorAny[T]) GetArrayInitialization(elements []environment.Node) environment.Node {
	return &nodes.ArrayInitialization[T]{
		Elements: elements,
	}
}

func (tn TypeNodeGeneratorAny[T]) GetArrayIndex(array environment.Node, index environment.Node) environment.Node {
	return &nodes.ArrayIndex[T]{
		Array: array,
		Index: index,
	}
}

func (tn TypeNodeGeneratorAny[T]) GetArrayAssignment(array environment.Node, index environment.Node, value environment.Node) environment.Node {
	return &nodes.ArrayAssignment[T]{
		ArrayIndex: &nodes.ArrayIndex[T]{
			Array: array,
			Index: index,
		},
		Value: value,
	}
}

func (tn TypeNodeGeneratorAny[T]) GetLoopArray(valIdentifier string, indexIdentifier string, array environment.Node, inner *nodes.Block) environment.Node {
	return &nodes.LoopArray[T]{
		ValIdentifier:   valIdentifier,
		IndexIdentifier: indexIdentifier,
		Array:           array,
		Inner:           inner,
	}
}

func (tn TypeNodeGeneratorAny[T]) ArrayIndexDetails(node environment.Node) (array environment.Node, index environment.Node, ok bool) {
	val, ok := node.(*nodes.ArrayIndex[T])
	return val.Array, val.Index, ok
}

func (tn TypeNodeGeneratorAny[T]) GetMathsOperation(operation nodes.MathsOperationType, leftSide environment.Node, rightSide environment.Node) environment.Node {
	panic("Cannot get maths operation on a non-number type")
}

func (tn TypeNodeGeneratorAny[T]) GetInequalityComparison(comparison nodes.ComparisonType, leftSide environment.Node, rightSide environment.Node) environment.Node {
	panic("Cannot get inequalty comparison on a non-number type")
}

// Implementation of TypeNodeGenerator with generic types
type TypeNodeGeneratorNumber[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float32 | float64] struct {
	TypeNodeGeneratorAny[T]
}

func (tn TypeNodeGeneratorNumber[T]) GetMathsOperation(operation nodes.MathsOperationType, leftSide environment.Node, rightSide environment.Node) environment.Node {
	return &nodes.MathsOperation[T]{
		Operation: operation,
		LeftSide:  leftSide,
		RightSide: rightSide,
	}
}

func (tn TypeNodeGeneratorNumber[T]) GetInequalityComparison(comparison nodes.ComparisonType, leftSide environment.Node, rightSide environment.Node) environment.Node {
	return &nodes.InequalityComparison[T]{
		Type:      comparison,
		LeftSide:  leftSide,
		RightSide: rightSide,
	}
}
