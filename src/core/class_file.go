package core

import (
	"bytes"
	"fmt"

	"github.com/Gustrb/jbm/src/utils"
)

// Spec: https://docs.oracle.com/javase/specs/jvms/se8/html/jvms-4.html

type ConstantPoolInfo struct {
	tag  uint8
	info interface{}
}

type ClassInfo struct {
	nameIndex uint16
}

type ConstantPoolIndexableInfo struct {
	classIndex       uint16
	nameAndTypeIndex uint16
}

type NameAndTypeInfo struct {
	nameIndex       uint16
	descriptorIndex uint16
}

type StringInfo struct {
	stringIndex uint16
}

type Numeric32BitsInfo struct {
	value uint32
}

type UTF8Info struct {
	bytes []byte
}

type FieldInfo struct{}
type MethodInfo struct{}
type AttributeInfo struct{}

type ClassFile struct {
	magic        uint32
	minorVersion uint16
	majorVersion uint16
	constantPool []ConstantPoolInfo
	accessFlags  uint16
	thisClass    uint16
	superClass   uint16
	interfaces   []uint16
	fields       []FieldInfo
	methods      []MethodInfo
	attributes   []AttributeInfo
}

const MagicNumber uint32 = 0xCAFEBABE

// Constant pool tags docs:
// https://docs.oracle.com/javase/specs/jvms/se8/html/jvms-4.html#jvms-4.4-140
const (
	CONSTANT_Class              uint8 = 7
	CONSTANT_Fieldref           uint8 = 9
	CONSTANT_Methodref          uint8 = 10
	CONSTANT_InterfaceMethodref uint8 = 11
	CONSTANT_String             uint8 = 8
	CONSTANT_Integer            uint8 = 3
	CONSTANT_Float              uint8 = 4
	CONSTANT_Long               uint8 = 5
	CONSTANT_Double             uint8 = 6
	CONSTANT_NameAndType        uint8 = 12
	CONSTANT_Utf8               uint8 = 1
	CONSTANT_MethodHandle       uint8 = 15
	CONSTANT_MethodType         uint8 = 16
	CONSTANT_InvokeDynamic      uint8 = 18
)

func classFileFromReader(reader *bytes.Reader) (ClassFile, error) {
	bigEndianReader := utils.NewBigEndianReaderFromReader(reader)
	classFile := ClassFile{}

	// the first 4 bytes are the magic number
	magic, err := bigEndianReader.ReadUint32()
	if err != nil {
		return classFile, err
	}

	if magic != MagicNumber {
		return classFile, fmt.Errorf("invalid magic number: %x", magic)
	}

	classFile.magic = magic

	// the next 2 bytes are the minor version
	minorVersion, err := bigEndianReader.ReadUint16()
	if err != nil {
		return classFile, err
	}

	classFile.minorVersion = minorVersion

	// the next 2 bytes are the major version
	majorVersion, err := bigEndianReader.ReadUint16()
	if err != nil {
		return classFile, err
	}

	classFile.majorVersion = majorVersion

	// The value of the constant_pool_count item is equal to the number of entries in the constant_pool table plus one
	constantPoolCount, err := bigEndianReader.ReadUint16()
	if err != nil {
		return classFile, err
	}

	constantPool := make([]ConstantPoolInfo, constantPoolCount-1)
	for i := 0; i < int(constantPoolCount)-2; i++ {
		cpInfo, err := constantPoolFromReader(bigEndianReader)
		if err != nil {
			return classFile, err
		}

		constantPool[i] = cpInfo
	}

	return classFile, nil
}

func constantPoolFromReader(reader *utils.BigEndianReader) (ConstantPoolInfo, error) {
	cpInfo := ConstantPoolInfo{}
	tag, err := reader.ReadUint8()

	if err != nil {
		return cpInfo, err
	}

	if tag == CONSTANT_Class {
		nameIndex, err := reader.ReadUint16()

		if err != nil {
			return cpInfo, err
		}

		cpInfo.info = ClassInfo{nameIndex}

	} else if tag == CONSTANT_Fieldref || tag == CONSTANT_Methodref || tag == CONSTANT_InterfaceMethodref {
		classIndex, err := reader.ReadUint16()
		if err != nil {
			return cpInfo, err
		}

		nameAndTypeIndex, err := reader.ReadUint16()

		if err != nil {
			return cpInfo, err
		}

		cpInfo.info = ConstantPoolIndexableInfo{classIndex, nameAndTypeIndex}
		fmt.Printf("classIndex: %d\n", classIndex)
		fmt.Printf("nameAndTypeIndex: %d\n", nameAndTypeIndex)
	} else if tag == CONSTANT_String {
		stringIndex, err := reader.ReadUint16()

		if err != nil {
			return cpInfo, err
		}

		cpInfo.info = StringInfo{stringIndex}
	} else if tag == CONSTANT_Integer || tag == CONSTANT_Float {
		value, err := reader.ReadUint32()

		if err != nil {
			return cpInfo, err
		}

		cpInfo.info = Numeric32BitsInfo{value}
	} else if tag == CONSTANT_NameAndType {
		nameIndex, err := reader.ReadUint16()
		if err != nil {
			return cpInfo, err
		}

		descriptorIndex, err := reader.ReadUint16()
		if err != nil {
			return cpInfo, err
		}

		cpInfo.info = NameAndTypeInfo{nameIndex, descriptorIndex}
	} else if tag == CONSTANT_Utf8 {
		length, err := reader.ReadUint16()

		if err != nil {
			return cpInfo, err
		}

		b, err := reader.ReadBytes(int(length))

		if err != nil {
			return cpInfo, err
		}

		cpInfo.info = UTF8Info{bytes: b}

	} else {
		return cpInfo, fmt.Errorf("invalid constant pool tag: %d", tag)
	}

	cpInfo.tag = tag

	return cpInfo, nil
}

func ExecuteClassFile(reader *bytes.Reader) error {
	classFile, err := classFileFromReader(reader)

	if err != nil {
		return err
	}

	fmt.Printf("Magic: %x\n", classFile.magic)
	fmt.Printf("Minor version: %d\n", classFile.minorVersion)
	fmt.Printf("Major version: %d\n", classFile.majorVersion)

	return nil
}
