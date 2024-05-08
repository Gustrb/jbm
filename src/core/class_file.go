package core

import (
	"bytes"
	"fmt"

	"github.com/Gustrb/jbm/src/utils"
)

// Spec: https://docs.oracle.com/javase/specs/jvms/se8/html/jvms-4.html

type ConstantPoolInfo struct {
	Tag  uint8
	Info interface{}
}

type ClassInfo struct {
	NameIndex uint16
}

type ConstantPoolIndexableInfo struct {
	ClassIndex       uint16
	NameAndTypeIndex uint16
}

type NameAndTypeInfo struct {
	NameIndex       uint16
	DescriptorIndex uint16
}

type StringInfo struct {
	StringIndex uint16
}

type Numeric32BitsInfo struct {
	Value uint32
}

type UTF8Info struct {
	Bytes []byte
}

type FieldInfo struct{}
type MethodInfo struct{}
type AttributeInfo struct{}

type ClassFile struct {
	Magic        uint32
	MinorVersion uint16
	MajorVersion uint16
	ConstantPool []ConstantPoolInfo
	AccessFlags  uint16
	ThisClass    uint16
	SuperClass   uint16
	Interfaces   []uint16
	Fields       []FieldInfo
	Methods      []MethodInfo
	Attributes   []AttributeInfo
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

var (
	ErrInvalidMagicNumber      = fmt.Errorf("invalid magic number")
	ErrInvalidConstantPoolSize = fmt.Errorf("invalid constant pool size")
)

func ClassFileFromReader(reader *bytes.Reader) (ClassFile, error) {
	bigEndianReader := utils.NewBigEndianReaderFromReader(reader)
	classFile := ClassFile{}

	// the first 4 bytes are the magic number
	magic, err := bigEndianReader.ReadUint32()
	if err != nil {
		return classFile, err
	}

	if magic != MagicNumber {
		return classFile, ErrInvalidMagicNumber
	}

	classFile.Magic = magic

	// the next 2 bytes are the minor version
	minorVersion, err := bigEndianReader.ReadUint16()
	if err != nil {
		return classFile, err
	}

	classFile.MinorVersion = minorVersion

	// the next 2 bytes are the major version
	majorVersion, err := bigEndianReader.ReadUint16()
	if err != nil {
		return classFile, err
	}

	classFile.MajorVersion = majorVersion

	// The value of the constant_pool_count item is equal to the number of entries in the constant_pool table plus one
	constantPoolCount, err := bigEndianReader.ReadUint16()
	if err != nil {
		return classFile, err
	}

	if constantPoolCount == 1 || constantPoolCount == 0 {
		return classFile, ErrInvalidConstantPoolSize
	}

	constantPool := make([]ConstantPoolInfo, constantPoolCount-1)
	for i := 0; i < int(constantPoolCount)-2; i++ {
		cpInfo, err := constantPoolFromReader(bigEndianReader)
		if err != nil {
			return classFile, err
		}

		constantPool[i] = cpInfo
	}

	classFile.ConstantPool = constantPool

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

		cpInfo.Info = ClassInfo{nameIndex}

	} else if tag == CONSTANT_Fieldref || tag == CONSTANT_Methodref || tag == CONSTANT_InterfaceMethodref {
		classIndex, err := reader.ReadUint16()
		if err != nil {
			return cpInfo, err
		}

		nameAndTypeIndex, err := reader.ReadUint16()

		if err != nil {
			return cpInfo, err
		}

		cpInfo.Info = ConstantPoolIndexableInfo{classIndex, nameAndTypeIndex}
	} else if tag == CONSTANT_String {
		stringIndex, err := reader.ReadUint16()

		if err != nil {
			return cpInfo, err
		}

		cpInfo.Info = StringInfo{stringIndex}
	} else if tag == CONSTANT_Integer || tag == CONSTANT_Float {
		value, err := reader.ReadUint32()

		if err != nil {
			return cpInfo, err
		}

		cpInfo.Info = Numeric32BitsInfo{value}
	} else if tag == CONSTANT_NameAndType {
		nameIndex, err := reader.ReadUint16()
		if err != nil {
			return cpInfo, err
		}

		descriptorIndex, err := reader.ReadUint16()
		if err != nil {
			return cpInfo, err
		}

		cpInfo.Info = NameAndTypeInfo{nameIndex, descriptorIndex}
	} else if tag == CONSTANT_Utf8 {
		length, err := reader.ReadUint16()

		if err != nil {
			return cpInfo, err
		}

		b, err := reader.ReadBytes(int(length))

		if err != nil {
			return cpInfo, err
		}

		cpInfo.Info = UTF8Info{Bytes: b}

	} else {
		return cpInfo, fmt.Errorf("invalid constant pool tag: %d", tag)
	}

	cpInfo.Tag = tag

	return cpInfo, nil
}

func ExecuteClassFile(reader *bytes.Reader) error {
	classFile, err := ClassFileFromReader(reader)

	if err != nil {
		return err
	}

	for i, cpInfo := range classFile.ConstantPool {
		fmt.Printf("Constant pool entry %d\n", i)

		if cpInfo.Tag == CONSTANT_Class {
			classInfo := cpInfo.Info.(ClassInfo)
			fmt.Printf("Class name index: %d\n", classInfo.NameIndex)
		} else if cpInfo.Tag == CONSTANT_Fieldref || cpInfo.Tag == CONSTANT_Methodref || cpInfo.Tag == CONSTANT_InterfaceMethodref {
			indexableInfo := cpInfo.Info.(ConstantPoolIndexableInfo)
			fmt.Printf("Class index: %d\n", indexableInfo.ClassIndex)
			fmt.Printf("Name and type index: %d\n", indexableInfo.NameAndTypeIndex)
		} else if cpInfo.Tag == CONSTANT_String {
			stringInfo := cpInfo.Info.(StringInfo)
			fmt.Printf("String index: %d\n", stringInfo.StringIndex)

			fmt.Printf("String: %s\n", string(classFile.ConstantPool[stringInfo.StringIndex-1].Info.(UTF8Info).Bytes))
		} else if cpInfo.Tag == CONSTANT_Integer || cpInfo.Tag == CONSTANT_Float {
			n32Info := cpInfo.Info.(Numeric32BitsInfo)
			fmt.Printf("Value: %d\n", n32Info.Value)
		} else if cpInfo.Tag == CONSTANT_NameAndType {
			nameAndTypeInfo := cpInfo.Info.(NameAndTypeInfo)
			fmt.Printf("Name index: %d\n", nameAndTypeInfo.NameIndex)
			fmt.Printf("Descriptor index: %d\n", nameAndTypeInfo.DescriptorIndex)
		} else if cpInfo.Tag == CONSTANT_Utf8 {
			utf8Info := cpInfo.Info.(UTF8Info)
			fmt.Printf("Bytes: %s\n", string(utf8Info.Bytes))
		}
	}

	return nil
}
