package interpreter

type GenericType uint8

const (
	TypeFunc GenericType = iota
	TypeMap
	TypeArray
	TypeInt8
	TypeInt16
	TypeInt32
	TypeInt48
	TypeInt64
	TypeUint8
	TypeUint16
	TypeUint32
	TypeUint48
	TypeUint64
	TypeFloat32
	TypeFloat64
	TypeString
	TypeBool
)

type TypeDef interface {
	GetGenericType() GenericType
	Equals(TypeDef) bool
}

type GenericTypeDef struct {
	Type GenericType
}

func (def GenericTypeDef) GetGenericType() GenericType {
	return def.Type
}

func (def GenericTypeDef) Equals(other TypeDef) bool {
	return def == other.(GenericTypeDef)
}

type FuncDef struct {
	GenericTypeDef
	Args       map[string]TypeDef
	ReturnType TypeDef
}

func (def FuncDef) Equals(other TypeDef) bool {
	otherDef, ok := other.(FuncDef)
	if !ok || len(def.Args) != len(otherDef.Args) {
		return false
	}
	for key, value := range def.Args {
		if otherDef.Args[key] != value {
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
	otherDef, ok := other.(MapDef)
	return ok && def.KeyType.Equals(otherDef.KeyType) && def.ValueType.Equals(otherDef.ValueType)
}

type ArrayDef struct {
	GenericTypeDef
	ElementType TypeDef
}

func (def ArrayDef) Equals(other TypeDef) bool {
	otherDef, ok := other.(ArrayDef)
	return ok && def.ElementType.Equals(otherDef.ElementType)
}
