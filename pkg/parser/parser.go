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
	// Magic              uint // マジックナンバー
	// MinorVersion       uint
	// MajorVersion       uint
	// ConstantPool_count uint     // constantPoolの長さに1足した数
	// ConstantPool       []CpInfo // 定数プール。クラス名やメソッド名、文字列などを定義
	// AccessFlags        uint     // クラスあるいはインターフェースの情報、アクセス制御に関するフラグ
	// ThisClass          uint     // このクラスあるいはインターフェースが何なのか。constant_poolで定義されているはずのこのクラス情報のインデックスが入る
	// SuperClass         uint     // 親クラスを示すconstant_poolのインデックスが入る
	// InterfacesCount    uint
	// Interfaces         []uint // このクラスが実装しているインターフェース情報。定数プールに定義されているインターフェースのインデックスが入る
	// FieldsCount        uint
	// Fields             []FieldInfo // 各フィールドの定義情報
	// MethodsCount       uint
	// Methods            MethodInfo
	// AttributesCount    uint
	// Attributes         AttributeInfo // その他の付加情報
}

type Data1 struct {
	Magic             uint32
	MinorVersion      uint16
	MajorVersion      uint16
	ConstantPoolCount uint16
}

func (cl *ClassFile) Run() {
	f, err := os.Open("../Main.class")
	if err != nil {
		panic(err)
	}

	data1 := Data1{}
	errb := binary.Read(f, binary.BigEndian, &data1)
	if errb != nil {
		panic(errb)
	}

	var constPoolItems []interface{}

	for i := 0; i < int(data1.ConstantPoolCount-1); i++ {
		var tag uint8
		errb = binary.Read(f, binary.BigEndian, &tag)
		if errb != nil {
			panic(errb)
		}
		switch tag {
		case 1:
			var length uint16
			errb = binary.Read(f, binary.BigEndian, &length)
			if errb != nil {
				panic(errb)
			}
			var bs []uint8
			for j := 0; j < int(length); j++ {
				var b uint8
				errb = binary.Read(f, binary.BigEndian, &b)
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
			errb = binary.Read(f, binary.BigEndian, &constClass)
			if errb != nil {
				panic(errb)
			}
			constPoolItems = append(constPoolItems, constClass)
		case 8:
			constString := ConstString{}
			errb = binary.Read(f, binary.BigEndian, &constString)
			if errb != nil {
				panic(errb)
			}
			constPoolItems = append(constPoolItems, constString)
		case 9:
			fieldRef := FieldRef{}
			errb = binary.Read(f, binary.BigEndian, &fieldRef)
			if errb != nil {
				panic(errb)
			}
			constPoolItems = append(constPoolItems, fieldRef)
		case 10:
			methodRef := MethodRef{}
			errb = binary.Read(f, binary.BigEndian, &methodRef)
			if errb != nil {
				panic(errb)
			}
			constPoolItems = append(constPoolItems, methodRef)
		case 12:
			nameandtype := NameAndType{}
			errb = binary.Read(f, binary.BigEndian, &nameandtype)
			if errb != nil {
				panic(errb)
			}
			constPoolItems = append(constPoolItems, nameandtype)
		default:
			panic(fmt.Sprintf("%d is not support!", tag))
		}
	}

	// check
	// tc, _ := constPoolItems[data2.ThisClass-1].(ConstClass)
	// sc, _ := constPoolItems[data2.SuperClass-1].(ConstClass)
	// fmt.Printf("%s\n", constPoolItems[tc.NameIdx])
	// fmt.Printf("%s\n", constPoolItems[sc.NameIdx])

	ReadMethod(f, constPoolItems)
}

type Method struct {
	AccessFlags     uint16
	NameIdx         uint16
	DescriptorIdx   uint16
	AttributesCount uint16
}

func ReadMethod(f *os.File, cpe []interface{}) {
	var methods_count uint16
	errb := binary.Read(f, binary.BigEndian, &methods_count)
	if errb != nil {
		panic(errb)
	}

	var method Method
	errb = binary.Read(f, binary.BigEndian, &method)
	if errb != nil {
		panic(errb)
	}

	// attributes
	for i := 0; i < int(method.AttributesCount); i++ {
		for j := 0; j < int(method.AttributesCount-1); j++ {
			ReadAttr(f, cpe)
		}
	}
}

func ReadAttr(f *os.File, cpe []interface{}) {
	var attributeNameIdx uint16
	errb := binary.Read(f, binary.BigEndian, &attributeNameIdx)
	if errb != nil {
		panic(errb)
	}

	var attributeLen uint32
	errb = binary.Read(f, binary.BigEndian, &attributeLen)
	if errb != nil {
		panic(errb)
	}

	attridx := cpe[attributeNameIdx]
	attridx_utf8, ok := attridx.(ConstUtf8)

	if ok {
		switch fmt.Sprintf("%s", attridx_utf8.Bytes) {
		case "Code":
			ReadCodeAttr(f, cpe, attributeLen)
		case "LineNumberTable":
			fmt.Println("linenum")
		case "SourceFile":
			fmt.Println("sourcefile")
		default:
			panic(fmt.Sprintf("%s is not implemented", attridx_utf8.Bytes))
		}
	}
}

func ReadCodeAttr(f *os.File, cpe []interface{}, attrLen uint32) {
	var maxStack uint16
	errb := binary.Read(f, binary.BigEndian, &maxStack)
	if errb != nil {
		panic(errb)
	}

	var maxLocals uint16
	errb = binary.Read(f, binary.BigEndian, &maxLocals)
	if errb != nil {
		panic(errb)
	}

	var codeLen uint32
	errb = binary.Read(f, binary.BigEndian, &codeLen)
	if errb != nil {
		panic(errb)
	}

	var code []uint8
	for i := 0; i < int(codeLen); i++ {
		var b uint8
		errb = binary.Read(f, binary.BigEndian, &b)
		if errb != nil {
			panic(errb)
		}
		code = append(code, b)
	}

	var exceptionTableLen uint16
	errb = binary.Read(f, binary.BigEndian, &exceptionTableLen)
	if errb != nil {
		panic(errb)
	}

	for i := 0; i < int(exceptionTableLen); i++ {
		// skip
	}

	var attrCount uint16
	errb = binary.Read(f, binary.BigEndian, &attrCount)
	if errb != nil {
		panic(errb)
	}

	for i := 0; i < int(attrCount); i++ {
		ReadAttr(f, cpe) // attr =
	}
}
