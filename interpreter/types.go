package interpreter

type GenericType uint8

const (
	TypeInt8 GenericType = iota
	TypeInt16
	TypeInt32
	TypeInt64
	TypeUint8
	TypeUint16
	TypeUint32
	TypeUint64
	TypeFloat32
	TypeFloat64
	TypeString
	TypeBool
	TypeMap
	TypeFunc
	TypeArray
	TypeAny

	TypeNil
)

type TypeDef interface {
	GetGenericType() GenericType
	Equals(TypeDef) bool
	IsNumber() bool
}

type GenericTypeDef struct {
	Type GenericType
}

func (def GenericTypeDef) GetGenericType() GenericType {
	return def.Type
}

func (def GenericTypeDef) Equals(other TypeDef) bool {
	if def.Type == TypeAny || other.GetGenericType() == TypeAny {
		return true
	}
	return def == other.(GenericTypeDef)
}

func (def GenericTypeDef) IsNumber() bool {
	genericType := def.GetGenericType()
	return genericType == TypeInt8 || genericType == TypeInt16 || genericType == TypeInt32 || genericType == TypeInt64 || genericType == TypeUint8 || genericType == TypeUint16 || genericType == TypeUint32 || genericType == TypeUint64 || genericType == TypeFloat32 || genericType == TypeFloat64
}

type FuncDef struct {
	GenericTypeDef
	Args []TypeDef
	// Whether or not the function has a variable number of arguments
	Variadic   bool
	ReturnType TypeDef
}

func (def FuncDef) Equals(other TypeDef) bool {
	if other.GetGenericType() == TypeAny {
		return true
	}
	otherDef, ok := other.(FuncDef)
	if !ok || len(def.Args) != len(otherDef.Args) {
		return false
	}
	for i, argDef := range def.Args {
		if otherDef.Args[i].Equals(argDef) {
			return false
		}
	}
	return def.ReturnType.Equals(otherDef.ReturnType)
}

type MapDef struct {
	GenericTypeDef
	KeyType   TypeDef
	ValueType TypeDef
}

func (def MapDef) Equals(other TypeDef) bool {
	if other.GetGenericType() == TypeAny {
		return true
	}
	otherDef, ok := other.(MapDef)
	return ok && def.KeyType.Equals(otherDef.KeyType) && def.ValueType.Equals(otherDef.ValueType)
}

type ArrayDef struct {
	GenericTypeDef
	ElementType TypeDef
}

func (def ArrayDef) Equals(other TypeDef) bool {
	if other.GetGenericType() == TypeAny {
		return true
	}
	otherDef, ok := other.(ArrayDef)
	return ok && def.ElementType.Equals(otherDef.ElementType)
}
