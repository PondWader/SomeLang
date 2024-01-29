package interpreter

type ScriptGenericType uint8

const (
	TypeFunc ScriptGenericType = iota
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

type ScriptType struct {
	GenericType ScriptGenericType
	Args        []ScriptType
	ReturnType  []ScriptType
}
