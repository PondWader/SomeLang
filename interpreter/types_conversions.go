package interpreter

// Maps go int64 values to other integer sizes using a language generic type definition
func ConvertInt64ToTypeDef(n int64, genericDef GenericType) any {
	switch genericDef {
	case TypeInt8:
		return int8(n)
	case TypeInt16:
		return int16(n)
	case TypeInt32:
		return int32(n)
	case TypeInt64:
		return int64(n)
	case TypeUint8:
		return uint8(n)
	case TypeUint16:
		return uint16(n)
	case TypeUint32:
		return uint32(n)
	case TypeUint64:
		return uint64(n)
	case TypeFloat32:
		return float32(n)
	case TypeFloat64:
		return float64(n)
	}
	return nil
}

// Maps go int64 values to other integer sizes using a language generic type definition
func ConvertFloat64ToTypeDef(n float64, genericDef GenericType) any {
	switch genericDef {
	case TypeFloat32:
		return float32(n)
	case TypeFloat64:
		return float64(n)
	}
	return nil
}
