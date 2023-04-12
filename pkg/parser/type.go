package parser

type ConstClass struct {
	NameIdx uint16
}

type ConstUtf8 struct {
	length uint8
	bytes  []uint8
}
