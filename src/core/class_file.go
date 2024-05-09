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

func (c *ConstantPoolInfo) String() string {
	switch c.Tag {
	case CONSTANT_Class:
		return fmt.Sprintf("ClassInfo{ NameIndex: %d }", c.Info.(ClassInfo).NameIndex)
	case CONSTANT_Fieldref, CONSTANT_Methodref, CONSTANT_InterfaceMethodref:
		return fmt.Sprintf("ConstantPoolIndexableInfo{ ClassIndex: %d, NameAndTypeIndex: %d }",
			c.Info.(ConstantPoolIndexableInfo).ClassIndex, c.Info.(ConstantPoolIndexableInfo).NameAndTypeIndex)
	case CONSTANT_String:
		return fmt.Sprintf("StringInfo{ StringIndex: %d }", c.Info.(StringInfo).StringIndex)
	case CONSTANT_Integer, CONSTANT_Float:
		return fmt.Sprintf("Numeric32BitsInfo{ Value: %d }", c.Info.(Numeric32BitsInfo).Value)
	case CONSTANT_NameAndType:
		return fmt.Sprintf("NameAndTypeInfo{ NameIndex: %d, DescriptorIndex: %d }",
			c.Info.(NameAndTypeInfo).NameIndex, c.Info.(NameAndTypeInfo).DescriptorIndex)
	case CONSTANT_Utf8:
		return fmt.Sprintf("UTF8Info{ Bytes: %s }", c.Info.(UTF8Info).Bytes)
	default:
		return fmt.Sprintf("Unknown constant pool tag: %d", c.Tag)
	}
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
	Tags                       = map[uint8]string{
		CONSTANT_Class:              "CONSTANT_Class",
		CONSTANT_Fieldref:           "CONSTANT_Fieldref",
		CONSTANT_Methodref:          "CONSTANT_Methodref",
		CONSTANT_InterfaceMethodref: "CONSTANT_InterfaceMethodref",
		CONSTANT_String:             "CONSTANT_String",
		CONSTANT_Integer:            "CONSTANT_Integer",
		CONSTANT_Float:              "CONSTANT_Float",
		CONSTANT_Long:               "CONSTANT_Long",
		CONSTANT_Double:             "CONSTANT_Double",
		CONSTANT_NameAndType:        "CONSTANT_NameAndType",
		CONSTANT_Utf8:               "CONSTANT_Utf8",
		CONSTANT_MethodHandle:       "CONSTANT_MethodHandle",
		CONSTANT_MethodType:         "CONSTANT_MethodType",
		CONSTANT_InvokeDynamic:      "CONSTANT_InvokeDynamic",
	}
)

// ClassFileFromReader reads a class file from a bytes.Reader
//
// It assumes the content of the reader is big-endian, so it uses a utils.BigEndianReader to read the data
// from the reader.
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

	classFile.ConstantPool = make([]ConstantPoolInfo, constantPoolCount-1)
	for i := 0; i < len(classFile.ConstantPool); i++ {
		cpInfo, err := classFile.constantPoolFromReader(bigEndianReader)
		if err != nil {
			return classFile, err
		}

		classFile.ConstantPool[i] = cpInfo
	}

	accessFlags, err := bigEndianReader.ReadUint16()
	if err != nil {
		return classFile, err
	}

	classFile.AccessFlags = accessFlags

	thisClass, err := bigEndianReader.ReadUint16()
	if err != nil {
		return classFile, err
	}

	classFile.ThisClass = thisClass

	superClass, err := bigEndianReader.ReadUint16()
	if err != nil {
		return classFile, err
	}

	classFile.SuperClass = superClass

	interfacesCount, err := bigEndianReader.ReadUint16()
	if err != nil {
		return classFile, err
	}

	classFile.Interfaces = make([]uint16, interfacesCount)
	for i := 0; i < len(classFile.Interfaces); i++ {
		interfaceIndex, err := bigEndianReader.ReadUint16()
		if err != nil {
			return classFile, err
		}

		classFile.Interfaces[i] = interfaceIndex
	}

	if err := classFile.Validate(); err != nil {
		return classFile, err
	}

	return classFile, nil
}

func (c *ClassFile) Validate() error {
	// TODO: implement
	return nil
}

func (c *ClassFile) constantPoolFromReader(reader *utils.BigEndianReader) (ConstantPoolInfo, error) {
	cpInfo := ConstantPoolInfo{}
	tag, err := reader.ReadUint8()

	if err != nil {
		return cpInfo, err
	}

	if tag == CONSTANT_Class {
		return c.readConstantPoolClassInfo(reader)
	}

	if tag == CONSTANT_Fieldref || tag == CONSTANT_Methodref || tag == CONSTANT_InterfaceMethodref {
		return c.readConstantPoolObjectRefInfo(reader, tag)
	}

	if tag == CONSTANT_String {
		return c.readConstantPoolStringInfo(reader)
	}

	if tag == CONSTANT_Integer || tag == CONSTANT_Float {
		return c.readConstantPoolNumeric32BitsInfo(reader, tag)
	}

	if tag == CONSTANT_NameAndType {
		return c.readConstantPoolNameAndTypeInfo(reader)
	}

	if tag == CONSTANT_Utf8 {
		return c.readConstantPoolUTF8Info(reader)
	}

	return cpInfo, fmt.Errorf("invalid constant pool tag: %d", tag)
}

func (c *ClassFile) readConstantPoolClassInfo(reader *utils.BigEndianReader) (ConstantPoolInfo, error) {
	cpInfo := ConstantPoolInfo{
		Tag: CONSTANT_Class,
	}

	nameIndex, err := reader.ReadUint16()
	if err != nil {
		return cpInfo, err
	}

	cpInfo.Info = ClassInfo{nameIndex}

	return cpInfo, nil
}

func (c *ClassFile) readConstantPoolObjectRefInfo(reader *utils.BigEndianReader, tag uint8) (ConstantPoolInfo, error) {
	cpInfo := ConstantPoolInfo{
		Tag: tag,
	}

	classIndex, err := reader.ReadUint16()
	if err != nil {
		return cpInfo, err
	}

	nameAndTypeIndex, err := reader.ReadUint16()
	if err != nil {
		return cpInfo, err
	}

	cpInfo.Info = ConstantPoolIndexableInfo{classIndex, nameAndTypeIndex}

	return cpInfo, nil
}

func (c *ClassFile) readConstantPoolStringInfo(reader *utils.BigEndianReader) (ConstantPoolInfo, error) {
	cpInfo := ConstantPoolInfo{
		Tag: CONSTANT_String,
	}

	stringIndex, err := reader.ReadUint16()
	if err != nil {
		return cpInfo, err
	}

	cpInfo.Info = StringInfo{stringIndex}

	return cpInfo, nil
}

func (c *ClassFile) readConstantPoolNumeric32BitsInfo(reader *utils.BigEndianReader, tag uint8) (ConstantPoolInfo, error) {
	cpInfo := ConstantPoolInfo{
		Tag: tag,
	}

	value, err := reader.ReadUint32()
	if err != nil {
		return cpInfo, err
	}

	cpInfo.Info = Numeric32BitsInfo{value}

	return cpInfo, nil
}

func (c *ClassFile) readConstantPoolNameAndTypeInfo(reader *utils.BigEndianReader) (ConstantPoolInfo, error) {
	cpInfo := ConstantPoolInfo{
		Tag: CONSTANT_NameAndType,
	}

	nameIndex, err := reader.ReadUint16()
	if err != nil {
		return cpInfo, err
	}

	descriptorIndex, err := reader.ReadUint16()
	if err != nil {
		return cpInfo, err
	}

	cpInfo.Info = NameAndTypeInfo{nameIndex, descriptorIndex}

	return cpInfo, nil
}

func (c *ClassFile) readConstantPoolUTF8Info(reader *utils.BigEndianReader) (ConstantPoolInfo, error) {
	cpInfo := ConstantPoolInfo{
		Tag: CONSTANT_Utf8,
	}

	length, err := reader.ReadUint16()
	if err != nil {
		return cpInfo, err
	}

	b, err := reader.ReadBytes(int(length))
	if err != nil {
		return cpInfo, err
	}

	cpInfo.Info = UTF8Info{Bytes: b}

	return cpInfo, nil
}

func ExecuteClassFile(reader *bytes.Reader) error {
	_, err := ClassFileFromReader(reader)

	if err != nil {
		return err
	}

	return nil
}
