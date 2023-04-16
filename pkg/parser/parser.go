package parser

import (
	"encoding/binary"
	"fmt"
	"os"
)

// u4 -> uint32
// u2 -> uint16
// u1 -> uint8

type ClassFile struct {
	File              *os.File
	Magic             uint32 // マジックナンバー
	MinorVersion      uint16
	MajorVersion      uint16
	ConstantPoolCount uint16        // constantPoolの長さに1足した数
	ConstantPool      []interface{} // 定数プール。クラス名やメソッド名、文字列などを定義
	// AccessFlags        uint     // クラスあるいはインターフェースの情報、アクセス制御に関するフラグ
	ThisClass       uint16 // このクラスあるいはインターフェースが何なのか。constant_poolで定義されているはずのこのクラス情報のインデックスが入る
	SuperClass      uint16 // 親クラスを示すconstant_poolのインデックスが入る
	InterfacesCount uint16
	// Interfaces         []uint // このクラスが実装しているインターフェース情報。定数プールに定義されているインターフェースのインデックスが入る
	FieldsCount uint16
	// Fields             []FieldInfo // 各フィールドの定義情報
	MethodsCount uint16
	// Methods            MethodInfo
	AttributesCount uint16
	// Attributes         AttributeInfo // その他の付加情報
}

func (cl *ClassFile) Run() {
	f, err := os.Open("../Main.class")
	if err != nil {
		panic(err)
	}
	cl.File = f

	data := struct {
		Magic             uint32
		MinorVersion      uint16
		MajorVersion      uint16
		ConstantPoolCount uint16
	}{}

	errb := binary.Read(cl.File, binary.BigEndian, &data)
	if errb != nil {
		panic(errb)
	}

	cl.Magic = data.Magic
	cl.MinorVersion = data.MinorVersion
	cl.MajorVersion = data.MajorVersion
	cl.ConstantPoolCount = data.ConstantPoolCount

	cl.ReadConstantPool()

	data2 := struct {
		AccessFlags     uint16
		ThisClass       uint16
		SuperClass      uint16
		InterfacesCount uint16
		// Interfaces      uint16 // skip
		FieldsCount uint16
		// Fields      FieldInfo // skip
		MethodsCount uint16
	}{}

	errb = binary.Read(f, binary.BigEndian, &data2)
	if errb != nil {
		panic(errb)
	}

	cl.ThisClass = data2.ThisClass
	cl.SuperClass = data2.SuperClass
	cl.InterfacesCount = data2.InterfacesCount
	cl.FieldsCount = data2.FieldsCount
	cl.MethodsCount = data2.MethodsCount

	cl.ReadMethods(cl.ConstantPool)

	var attrCount uint16
	errb = binary.Read(f, binary.BigEndian, &attrCount)
	if errb != nil {
		panic(errb)
	}

	for i := 0; i < int(attrCount); i++ {
		cl.ReadAttr(cl.ConstantPool)
	}

	cl.AttributesCount = attrCount
}

func (cl *ClassFile) ReadConstantPool() {
	var constPoolItems []interface{}

	for i := 0; i < int(cl.ConstantPoolCount-1); i++ {
		var tag uint8
		errb := binary.Read(cl.File, binary.BigEndian, &tag)
		if errb != nil {
			panic(errb)
		}
		switch tag {
		case 1:
			var length uint16
			errb = binary.Read(cl.File, binary.BigEndian, &length)
			if errb != nil {
				panic(errb)
			}
			var bs []uint8
			for j := 0; j < int(length); j++ {
				var b uint8
				errb = binary.Read(cl.File, binary.BigEndian, &b)
				if errb != nil {
					panic(errb)
				}
				bs = append(bs, b)
			}
			utf8 := ConstUtf8{
				Length: length,
				Bytes:  bs,
			}
			constPoolItems = append(constPoolItems, utf8)
		case 7:
			constClass := ConstClass{}
			errb = binary.Read(cl.File, binary.BigEndian, &constClass)
			if errb != nil {
				panic(errb)
			}
			constPoolItems = append(constPoolItems, constClass)
		case 8:
			constString := ConstString{}
			errb = binary.Read(cl.File, binary.BigEndian, &constString)
			if errb != nil {
				panic(errb)
			}
			constPoolItems = append(constPoolItems, constString)
		case 9:
			fieldRef := FieldRef{}
			errb = binary.Read(cl.File, binary.BigEndian, &fieldRef)
			if errb != nil {
				panic(errb)
			}
			constPoolItems = append(constPoolItems, fieldRef)
		case 10:
			methodRef := MethodRef{}
			errb = binary.Read(cl.File, binary.BigEndian, &methodRef)
			if errb != nil {
				panic(errb)
			}
			constPoolItems = append(constPoolItems, methodRef)
		case 12:
			nameandtype := NameAndType{}
			errb = binary.Read(cl.File, binary.BigEndian, &nameandtype)
			if errb != nil {
				panic(errb)
			}
			constPoolItems = append(constPoolItems, nameandtype)
		default:
			panic(fmt.Sprintf("%d is not support!", tag))
		}
	}

	cl.ConstantPool = constPoolItems
}

func (cl *ClassFile) ReadMethods(cpe []interface{}) {
	for i := 0; i < int(cl.MethodsCount); i++ {
		m := struct {
			AccessFlags     uint16
			NameIdx         uint16
			DescriptorIdx   uint16
			AttributesCount uint16
		}{}
		errb := binary.Read(cl.File, binary.BigEndian, &m)
		if errb != nil {
			panic(errb)
		}

		for i := 0; i < int(m.AttributesCount); i++ {
			cl.ReadAttr(cpe)
		}
	}
}

func (cl *ClassFile) ReadAttr(cpe []interface{}) {
	var attributeNameIdx uint16
	errb := binary.Read(cl.File, binary.BigEndian, &attributeNameIdx)
	if errb != nil {
		panic(errb)
	}

	var attributeLen uint32
	errb = binary.Read(cl.File, binary.BigEndian, &attributeLen)
	if errb != nil {
		panic(errb)
	}

	attridx := cpe[attributeNameIdx-1]
	attridx_utf8, ok := attridx.(ConstUtf8)

	if ok {
		switch fmt.Sprintf("%s", attridx_utf8.Bytes) {
		case "Code":
			cl.ReadCodeAttr(cpe, attributeLen)
		case "LineNumberTable":
			cl.ReadLineNumTableAttr(cpe)
		case "SourceFile":
			cl.ReadSourceFileAttr(cpe)
		default:
			panic(fmt.Sprintf("%s is not implemented", attridx_utf8.Bytes))
		}
	} else {
		panic("attr is not ConstUtf8")
	}
}

func (cl *ClassFile) ReadCodeAttr(cpe []interface{}, attrLen uint32) {
	var maxStack uint16
	errb := binary.Read(cl.File, binary.BigEndian, &maxStack)
	if errb != nil {
		panic(errb)
	}

	var maxLocals uint16
	errb = binary.Read(cl.File, binary.BigEndian, &maxLocals)
	if errb != nil {
		panic(errb)
	}

	var codeLen uint32
	errb = binary.Read(cl.File, binary.BigEndian, &codeLen)
	if errb != nil {
		panic(errb)
	}

	var code []uint8
	for i := 0; i < int(codeLen); i++ {
		var b uint8
		errb = binary.Read(cl.File, binary.BigEndian, &b)
		if errb != nil {
			panic(errb)
		}
		code = append(code, b)
	}

	var exceptionTableLen uint16
	errb = binary.Read(cl.File, binary.BigEndian, &exceptionTableLen)
	if errb != nil {
		panic(errb)
	}

	for i := 0; i < int(exceptionTableLen); i++ {
		// skip
	}

	var attrCount uint16
	errb = binary.Read(cl.File, binary.BigEndian, &attrCount)
	if errb != nil {
		panic(errb)
	}

	for i := 0; i < int(attrCount); i++ {
		cl.ReadAttr(cpe) // attr =
	}
}

func (cl *ClassFile) ReadLineNumTableAttr(cpe []interface{}) {
	var lineNumberTableLength uint16
	errb := binary.Read(cl.File, binary.BigEndian, &lineNumberTableLength)
	if errb != nil {
		panic(errb)
	}

	for i := 0; i < int(lineNumberTableLength); i++ {
		linumtable := struct {
			StartPC uint16
			LineNum uint16
		}{}
		errb = binary.Read(cl.File, binary.BigEndian, &linumtable)
		if errb != nil {
			panic(errb)
		}
	}
}

func (cl *ClassFile) ReadSourceFileAttr(cpe []interface{}) {
	var sourcefileIdx uint16
	errb := binary.Read(cl.File, binary.BigEndian, &sourcefileIdx)
	if errb != nil {
		panic(errb)
	}
}
