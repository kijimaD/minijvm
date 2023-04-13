package parser

import (
	"encoding/binary"
	"fmt"
	"os"
)

type ClassFile struct {
	Magic              uint // マジックナンバー
	MinorVersion       uint
	MajorVersion       uint
	ConstantPool_count uint     // constantPoolの長さに1足した数
	ConstantPool       []CpInfo // 定数プール。クラス名やメソッド名、文字列などを定義
	AccessFlags        uint     // クラスあるいはインターフェースの情報、アクセス制御に関するフラグ
	ThisClass          uint     // このクラスあるいはインターフェースが何なのか。constant_poolで定義されているはずのこのクラス情報のインデックスが入る
	SuperClass         uint     // 親クラスを示すconstant_poolのインデックスが入る
	InterfacesCount    uint
	Interfaces         []uint // このクラスが実装しているインターフェース情報。定数プールに定義されているインターフェースのインデックスが入る
	FieldsCount        uint
	Fields             []FieldInfo // 各フィールドの定義情報
	MethodsCount       uint
	Methods            MethodInfo
	AttributesCount    uint
	Attributes         AttributeInfo // その他の付加情報
}

type FieldInfo struct {
	Dummy uint
}

type MethodInfo struct {
	Dummy uint
}

type AttributeInfo struct {
	Dummy uint
}

type CpInfo struct {
	Tag  uint8
	Info []uint
}

type Data1 struct {
	Magic             uint32
	MinorVersion      uint16
	MajorVersion      uint16
	ConstantPoolCount uint16
}

type Data2 struct {
	AccessFlags     uint16
	ThisClass       uint16
	SuperClass      uint16
	InterfacesCount uint16
	// Interfaces      uint16 // skip
	FieldsCount uint16
	// Fields      FieldInfo // skip
	MethodsCount uint16
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
				length: length,
				bytes:  bs,
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

	data2 := Data2{}
	errb = binary.Read(f, binary.BigEndian, &data2)
	if errb != nil {
		panic(errb)
	}
	// check
	// tc, _ := constPoolItems[data2.ThisClass-1].(ConstClass)
	// sc, _ := constPoolItems[data2.SuperClass-1].(ConstClass)
	// fmt.Printf("%s\n", constPoolItems[tc.NameIdx])
	// fmt.Printf("%s\n", constPoolItems[sc.NameIdx])
}
