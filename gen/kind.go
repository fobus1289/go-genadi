package gen

type GenKind int8

const (
	Controller       GenKind = 1
	ServiceStruct    GenKind = 2
	ServiceInterface GenKind = 3
	Configure        GenKind = 4
)

type Kind int8

const (
	Invalid Kind = iota - 1
	Bool
	String
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Float32
	Float64
	Complex32
	Complex64
	Struct
	Interface
	Map

	SliceBool
	SliceString
	SliceInt
	SliceInt8
	SliceInt16
	SliceInt32
	SliceInt64
	SliceUint
	SliceUint8
	SliceUint16
	SliceUint32
	SliceUint64
	SliceFloat32
	SliceFloat64
	SliceComplex32
	SliceComplex64
	SliceStruct
	SliceInterface
	SliceMap
	SliceByte
)

func (k Kind) IsValid() bool {
	return k != -1
}

func (k Kind) Name() string {
	switch k {
	case Bool:
		return "bool"
	case String:
		return "string"
	case Int:
		return "string"
	case Int8:
		return "string"
	case Int16:
		return "string"
	case Int32:
		return "string"
	case Int64:
		return "string"
	case Uint:
		return "string"
	case Uint8:
		return "string"
	case Uint16:
		return "string"
	case Uint32:
		return "string"
	case Uint64:
	case Float32:
	case Float64:
	case Complex32:
	case Complex64:
	case Struct:
	case Interface:

	}

	return ""
}
