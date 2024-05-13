package core

import (
	"bytes"
	"fmt"

	"github.com/Gustrb/jbm/src/utils"
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

	fieldsCount, err := bigEndianReader.ReadUint16()
	if err != nil {
		return classFile, err
	}

	classFile.Fields = make([]FieldInfo, fieldsCount)
	for i := 0; i < len(classFile.Fields); i++ {
		f, err := classFile.fieldInfoFromReader(bigEndianReader)
		if err != nil {
			return classFile, err
		}

		classFile.Fields[i] = f
	}

	methodsCount, err := bigEndianReader.ReadUint16()
	if err != nil {
		return classFile, err
	}

	classFile.Methods = make([]MethodInfo, methodsCount)
	for i := 0; i < len(classFile.Methods); i++ {
		m, err := classFile.methodInfoFromReader(bigEndianReader)
		if err != nil {
			return classFile, err
		}

		classFile.Methods[i] = m
	}

	attributesCount, err := bigEndianReader.ReadUint16()
	if err != nil {
		return classFile, err
	}

	classFile.Attributes = make([]AttributeInfo, attributesCount)
	for i := 0; i < len(classFile.Attributes); i++ {
		a, err := classFile.attributeInfoFromReader(bigEndianReader)
		if err != nil {
			return classFile, err
		}

		classFile.Attributes[i] = a
	}

	if err := classFile.Validate(); err != nil {
		return classFile, err
	}

	return classFile, nil
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

func (c *ClassFile) fieldInfoFromReader(reader *utils.BigEndianReader) (FieldInfo, error) {
	fInfo := FieldInfo{}

	accessFlags, err := reader.ReadUint16()
	if err != nil {
		return fInfo, err
	}

	fInfo.AccessFlags = accessFlags

	nameIndex, err := reader.ReadUint16()
	if err != nil {
		return fInfo, err
	}

	fInfo.NameIndex = nameIndex

	descriptorIndex, err := reader.ReadUint16()
	if err != nil {
		return fInfo, err
	}

	fInfo.DescriptorIndex = descriptorIndex

	attributesCount, err := reader.ReadUint16()
	if err != nil {
		return fInfo, err
	}

	fInfo.Attributes = make([]AttributeInfo, attributesCount)
	for i := 0; i < len(fInfo.Attributes); i++ {
		attr, err := c.attributeInfoFromReader(reader)
		if err != nil {
			return fInfo, err
		}

		fInfo.Attributes[i] = attr
	}

	return fInfo, nil
}

func (c *ClassFile) attributeInfoFromReader(reader *utils.BigEndianReader) (AttributeInfo, error) {
	attr := AttributeInfo{}

	nameIndex, err := reader.ReadUint16()
	if err != nil {
		return attr, nil
	}

	attr.AttributeNameIndex = nameIndex

	attributeLength, err := reader.ReadUint32()
	if err != nil {
		return attr, nil
	}

	info, err := reader.ReadBytes(int(attributeLength))
	if err != nil {
		return attr, err
	}

	attr.Info = info

	return attr, nil
}

// methodInfoFromReader reads a method_info structure from the reader, it assumes the reader
// is correctly positioned at the beggining of the structure.
//
// Here, we don't validate the method_info structure, we just return an error if there is any
// kind of IO problem, + we return the incomplete `MethodInfo` structure.
func (c *ClassFile) methodInfoFromReader(reader *utils.BigEndianReader) (MethodInfo, error) {
	mInfo := MethodInfo{}

	accessFlags, err := reader.ReadUint16()
	if err != nil {
		return mInfo, err
	}

	mInfo.AccessFlags = accessFlags

	nameIndex, err := reader.ReadUint16()
	if err != nil {
		return mInfo, err
	}

	mInfo.NameIndex = nameIndex

	descriptorIndex, err := reader.ReadUint16()
	if err != nil {
		return mInfo, err
	}

	mInfo.DescriptorIndex = descriptorIndex

	attributesCount, err := reader.ReadUint16()
	if err != nil {
		return mInfo, err
	}

	mInfo.Attributes = make([]AttributeInfo, attributesCount)
	for i := 0; i < len(mInfo.Attributes); i++ {
		attr, err := c.attributeInfoFromReader(reader)
		if err != nil {
			return mInfo, err
		}

		mInfo.Attributes[i] = attr
	}

	return mInfo, nil
}
