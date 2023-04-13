package parser

type ConstUtf8 struct {
	length uint16
	bytes  []uint8
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
