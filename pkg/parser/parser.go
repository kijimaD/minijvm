package parser

type ClassFile struct {
	magic              uint // マジックナンバー
	minorVersion       uint
	majorVersion       uint
	constantPool_count uint     // constantPoolの長さに1足した数
	constantPool       []cpInfo // 定数プール。クラス名やメソッド名、文字列などを定義
	accessFlags        uint     // クラスあるいはインターフェースの情報、アクセス制御に関するフラグ
	thisClass          uint     // このクラスあるいはインターフェースが何なのか。constant_poolで定義されているはずのこのクラス情報のインデックスが入る
	superClass         uint     // 親クラスを示すconstant_poolのインデックスが入る
	interfacesCount    uint
	interfaces         []uint // このクラスが実装しているインターフェース情報。定数プールに定義されているインターフェースのインデックスが入る
	fieldsCount        uint
	fields             []fieldInfo // 各フィールドの定義情報
	methodsCount       uint
	methods            methodInfo
	attributesCount    uint
	attributes         attributeInfo // その他の付加情報
}

type cpInfo struct {
	tag  uint
	info []uint
}

type fieldInfo struct {
}

// func (cpu *CPU) opcode(line uint16) uint16 {
// 	return (line >> 11)
// }

// func (cpu *CPU) opregA(line uint16) uint16 {
// 	return ((line >> 8) & 0x0007)
// }
