package core

import (
	"bytes"
	"fmt"

	"github.com/Gustrb/jbm/src/utils"
)

// Spec: https://docs.oracle.com/javase/specs/jvms/se8/html/jvms-4.html

type ConstantPoolInfo struct{}
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

	return classFile, nil
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
