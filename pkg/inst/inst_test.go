package inst

import (
	"bytes"
	"testing"

	"github.com/kijimaD/minijvm/pkg/parser"
)

func TestRun(t *testing.T) {
	// assert := assert.New(t)

	cl := parser.ClassFile{}
	cl.Run()

	// for _, m := range cl.Methods {
	// 	for _, a := range m.Attributes {
	// 		codeattr, _ := a.(parser.CodeAttr)

	// 		code := bytes.NewBuffer(codeattr.Code)
	// 		inst := Inst{
	// 			Code: code,
	// 		}
	// 		inst.Run()
	// 	}
	// }

	codeattr, _ := cl.Methods[1].Attributes[0].(parser.CodeAttr)
	code := bytes.NewBuffer(codeattr.Code)
	inst := Inst{
		Code:  code,
		Pools: cl.ConstantPool,
	}
	inst.Run()
}
