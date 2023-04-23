package inst

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/kijimaD/minijvm/pkg/parser"
)

type Inst struct {
	Code  *bytes.Buffer
	Pools []interface{}
}

type Operand struct {
	callableValue string // たぶん、リフレクションで取得した値を入れる。callで実行できるような感じの
	returnValue   string
}

func (i *Inst) Run() {
	var stack []interface{}
MainLoop:
	for {
		var opcode uint8
		errb := binary.Read(i.Code, binary.BigEndian, &opcode)
		if errb != nil {
			panic(errb)
		}

		fmt.Printf("%#v\n", opcode)

		switch opcode {
		case 0xb2: // getstatic
			var opr uint16
			errb = binary.Read(i.Code, binary.BigEndian, &opr)
			if errb != nil {
				panic(errb)
			}

			// java/lang/System
			fieldRef := i.Pools[opr-1].(parser.FieldRef) // なんかインデックスが1つずれてる?
			constClass := i.Pools[fieldRef.ClassIdx-1].(parser.ConstClass)
			classname := i.Pools[constClass.NameIdx-1].(parser.ConstUtf8)

			// out
			nameAndType := i.Pools[fieldRef.NameAndTypeIdx-1].(parser.NameAndType)
			fieldname := i.Pools[nameAndType.NameIdx-1].(parser.ConstUtf8)

			// フィールドの型情報
			descriptor := i.Pools[nameAndType.DescriptorIdx].(parser.ConstUtf8)

			o := Operand{
				callableValue: "a", // class, field から取得
				returnValue:   fmt.Sprintf("%s", descriptor.Bytes),
			}
			stack = append(stack, o)
			fmt.Printf("getstatic: %s %s %s\n", classname.Bytes, fieldname.Bytes, descriptor.Bytes)
		case 0x12: // ldc
			var opr uint8
			errb = binary.Read(i.Code, binary.BigEndian, &opr)
			if errb != nil {
				panic(errb)
			}

			constString := i.Pools[opr-1].(parser.ConstString)
			symNameIdx := constString.StringIdx
			str := i.Pools[symNameIdx-1].(parser.ConstUtf8)
			stack = append(stack, str.Bytes)
			fmt.Printf("ldc: %s\n", str.Bytes)
		case 0xb6: // invoke virtual
			var opr uint16
			errb = binary.Read(i.Code, binary.BigEndian, &opr)
			if errb != nil {
				panic(errb)
			}
			break
		case 0xb1: // return
			break MainLoop
		default:
			panic(fmt.Sprintf("[%d] is not support!", opcode))
		}
	}
}
