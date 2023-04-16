package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	assert := assert.New(t)

	cl := ClassFile{}
	cl.Run()

	assert.Equal(0, int(cl.MinorVersion))
	assert.Equal(55, int(cl.MajorVersion))
	assert.Equal(29, int(cl.ConstantPoolCount))

	{
		cla, _ := cl.ConstantPool[cl.ThisClass-1].(ConstClass)
		utf8, _ := cl.ConstantPool[cla.NameIdx].(ConstUtf8)
		name := fmt.Sprintf("%s", utf8.Bytes)
		assert.Equal("java/lang/Object", name)
	}

	{
		cla, _ := cl.ConstantPool[cl.SuperClass-1].(ConstClass)
		utf8, _ := cl.ConstantPool[cla.NameIdx].(ConstUtf8)
		name := fmt.Sprintf("%s", utf8.Bytes)
		assert.Equal("java/lang/System", name)
	}

	assert.Equal(0, int(cl.InterfacesCount))
	assert.Equal(0, len(cl.Interfaces))
	assert.Equal(2, int(cl.MethodsCount))
	assert.Equal(2, len(cl.Methods))
	{
		utf8, _ := cl.ConstantPool[cl.Methods[0].NameIdx-1].(ConstUtf8)
		name := fmt.Sprintf("%s", utf8.Bytes)
		assert.Equal("<init>", name)
	}
	{
		utf8, _ := cl.ConstantPool[cl.Methods[1].NameIdx-1].(ConstUtf8)
		name := fmt.Sprintf("%s", utf8.Bytes)
		assert.Equal("main", name)
	}

	assert.Equal(0, int(cl.FieldsCount))
	assert.Equal(1, int(cl.AttributesCount))
}
