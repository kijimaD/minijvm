package parser

import "os"

type ClassFile struct {
	File              *os.File
	Magic             uint32 // マジックナンバー
	MinorVersion      uint16
	MajorVersion      uint16
	ConstantPoolCount uint16        // constantPoolの長さに1足した数
	ConstantPool      []interface{} // 定数プール。クラス名やメソッド名、文字列などを定義。
	AccessFlags       uint16        // クラスあるいはインターフェースの情報、アクセス制御に関するフラグ
	ThisClass         uint16        // このクラスあるいはインターフェースが何なのか。constant_poolで定義されているはずのこのクラス情報のインデックスが入る
	SuperClass        uint16        // 親クラスを示すconstant_poolのインデックスが入る
	InterfacesCount   uint16
	Interfaces        []uint16 // このクラスが実装しているインターフェース情報。定数プールに定義されているインターフェースのインデックスが入る
	FieldsCount       uint16
	// Fields             []FieldInfo // 各フィールドの定義情報
	MethodsCount    uint16
	Methods         []MethodInfo
	AttributesCount uint16
	Attributes      []interface{} // その他の付加情報
}

///////////////////////////
// Constant Pool structs //
///////////////////////////

type ConstUtf8 struct {
	Length uint16
	Bytes  []uint8
}

type MethodInfo struct {
	AccessFlags     uint16 // メソッドへのアクセス権
	NameIdx         uint16
	DescriptorIdx   uint16 // 引数の銃砲
	AttributesCount uint16
	Attributes      []interface{}
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

///////////////////////
// Attribute structs //
///////////////////////

type CodeAttr struct {
	AttrNameIdx          uint16
	AttrLength           uint32
	MaxStack             uint16
	MaxLocals            uint16
	CodeLength           uint32
	Code                 []uint8 // オペコードが格納されている
	ExceptionTableLength uint16
	// ExceptionTable []exceptionTable
	AttrCount uint16
	// AttrInfoAttributes []interface{}
}

type LineNumberTableAttr struct {
	AttrNameIdx     uint16
	AttrLength      uint16
	LineNumberTable []LineNumberTable
}

type LineNumberTable struct {
	StartPC    uint16
	LineNumber uint16
}

type SourceFileAttr struct {
	AttrNameIdx   uint16
	AttrLength    uint32
	SourceFileIdx uint16
}
