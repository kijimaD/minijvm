package parser

type ConstUtf8 struct {
	Length uint16
	Bytes  []uint8
}

type ConstClass struct {
	NameIdx uint16
}

type ConstString struct {
	StringIdx uint16
}

type FieldRef struct {
	ClassIdx       uint16
	NameAndTypeIdx uint16
}

type MethodRef struct {
	ClassIdx       uint16
	NameAndTypeIdx uint16
}

type NameAndType struct {
	NameIdx       uint16
	DescriptorIdx uint16
}
