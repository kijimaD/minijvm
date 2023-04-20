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

func (i *Inst) Run() {
	fmt.Printf("%#v\n", i.Code)
	for {
		var opcode uint8
		errb := binary.Read(i.Code, binary.BigEndian, &opcode)
		if errb != nil {
			panic(errb)
		}

		switch opcode {
		case 0xb2: // getstatic
			var operand uint16
			errb = binary.Read(i.Code, binary.BigEndian, &operand)
			if errb != nil {
				panic(errb)
			}

			// java/lang/System
			fieldRef := i.Pools[operand-1].(parser.FieldRef) // なんかインデックスが1つずれてる?
			constClass := i.Pools[fieldRef.ClassIdx-1].(parser.ConstClass)
			classname := i.Pools[constClass.NameIdx-1].(parser.ConstUtf8)

			// out
			nameAndType := i.Pools[fieldRef.NameAndTypeIdx-1].(parser.NameAndType)
			fieldname := i.Pools[nameAndType.NameIdx-1].(parser.ConstUtf8)

			// フィールドの型情報
			descriptor := i.Pools[nameAndType.DescriptorIdx].(parser.ConstUtf8)

			fmt.Printf("%s %s $s\n", classname.Bytes, fieldname.Bytes, descriptor.Bytes)
		case 0x12:
		case 0xb6:
		case 0xb1:
		default:
			panic(fmt.Sprintf("[%d] is not support!", opcode))
		}
	}
}
