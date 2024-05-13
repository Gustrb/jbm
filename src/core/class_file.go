package core

import (
	"bytes"
	"fmt"
)

// Spec: https://docs.oracle.com/javase/specs/jvms/se8/html/jvms-4.html

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

// MagicNumber is the magic number of a Java class file. It is always 0xCAFEBABE.
const MagicNumber uint32 = 0xCAFEBABE

const (
	ACC_PUBLIC     uint16 = 0x0001
	ACC_FINAL      uint16 = 0x0010
	ACC_SUPER      uint16 = 0x0020
	ACC_INTERFACE  uint16 = 0x0200
	ACC_ABSTRACT   uint16 = 0x0400
	ACC_SYNTHETIC  uint16 = 0x1000
	ACC_ANNOTATION uint16 = 0x2000
	ACC_ENUM       uint16 = 0x4000
)

// ConstantPoolInfo represents an element inside the `ConstantPool`
// of a Java class file.
type ConstantPoolInfo struct {
	// Tag is the type of the constant pool entry.
	Tag uint8
	// Info is the actual data of the constant pool entry, which can be of different types
	// depending on the value of `Tag`.
	Info interface{}
}

// ClassInfo represents a CONSTANT_Class_info structure in the constant pool.
// It is used to represent a class or an interface.
type ClassInfo struct {
	// NameIndex is the index of a UTF-8 entry in the constant pool that represents the name of the class.
	NameIndex uint16
}

// ConstantPoolIndexableInfo represents a CONSTANT_Fieldref_info, CONSTANT_Methodref_info or
// CONSTANT_InterfaceMethodref_info structure in the constant pool.
//
// It is used to represent a field, method or interface method respectively.
type ConstantPoolIndexableInfo struct {
	// ClassIndex is the index of a CONSTANT_Class_info structure in the constant pool.
	ClassIndex uint16
	// NameAndTypeIndex is the index of a CONSTANT_NameAndType_info structure in the constant pool.
	NameAndTypeIndex uint16
}

// NameAndTypeInfo represents a CONSTANT_NameAndType_info structure in the constant pool.
// It is used to represent a field or method name and type descriptor.
type NameAndTypeInfo struct {
	// NameIndex is the index of a UTF-8 entry in the constant pool that represents the name of the field or method.
	NameIndex uint16
	// DescriptorIndex is the index of a UTF-8 entry in the constant pool that
	// represents the descriptor of the field or method.
	// The descriptor is a string representing the type of the field or method.
	DescriptorIndex uint16
}

// StringInfo represents a CONSTANT_String_info structure in the constant pool.
type StringInfo struct {
	// StringIndex is the index of a UTF-8 entry in the constant pool that represents the string value.
	StringIndex uint16
}

type Numeric32BitsInfo struct {
	Value uint32
}

// UTF8Info represents a CONSTANT_Utf8_info structure in the constant pool.
// It is used to represent a string value.
type UTF8Info struct {
	// Bytes is the UTF-8 encoded string.
	Bytes []byte
}

// AttributeInfo represents an attribute of a class file, field or method.
type AttributeInfo struct {
	// AttributeNameIndex is the index of a UTF-8 entry in the constant pool that represents the name of the attribute.
	AttributeNameIndex uint16
	// Info is the actual data of the attribute, the structure of which depends on the attribute name.
	// for example, the `Code` attribute has its own structure.
	// The `Info` field is a byte slice that contains the raw data of the attribute.
	Info []byte
}

// FieldInfo represents a field of a class.
// It contains the access flags, name, descriptor and attributes of the field.
type FieldInfo struct {
	// AccessFlags is a mask of flags used to denote access permissions to and properties of the field.
	AccessFlags uint16
	// NameIndex is the index of a UTF-8 entry in the constant pool that represents the name of the field.
	NameIndex uint16
	// DescriptorIndex is the index of a UTF-8 entry in the constant pool that represents the descriptor of the field.
	// The descriptor is a string representing the type of the field.
	DescriptorIndex uint16
	// Attributes is a list of attributes of the field.
	Attributes []AttributeInfo
}

// MethodInfo represents a method of a class.
type MethodInfo struct {
	// AccessFlags is a mask of flags used to denote access permissions to and properties of the method.
	AccessFlags uint16
	// NameIndex is the index of a UTF-8 entry in the constant pool that represents the name of the method.
	NameIndex uint16
	// DescriptorIndex is the index of a UTF-8 entry in the constant pool that represents the descriptor of the method.
	// The descriptor is a string representing the type of the method.
	DescriptorIndex uint16
	// Attributes is a list of attributes of the method.
	Attributes []AttributeInfo
}

// ClassFile represents the structure of a Java class file.
type ClassFile struct {
	// Magic is the magic number of the class file.
	// It is always 0xCAFEBABE.
	Magic uint32
	// MinorVersion is the minor version of the class file.
	// It is used to indicate changes to the class file that are not compatible with previous versions.
	MinorVersion uint16
	// MajorVersion is the major version of the class file.
	// It is used to indicate changes to the class file that are not compatible with previous versions.
	MajorVersion uint16
	// ConstantPool is the constant pool of the class file.
	// The constant pool is a table of structures representing various constants.
	ConstantPool []ConstantPoolInfo
	// AccessFlags is a mask of flags used to denote access permissions to and properties of the class.
	// The flags are used to denote if the class is public, final, etc.
	AccessFlags uint16
	// ThisClass is the index of a CONSTANT_Class_info structure in the constant pool.
	// It represents the class or interface defined by the class file.
	ThisClass uint16
	// SuperClass is the index of a CONSTANT_Class_info structure in the constant pool.
	// It represents the direct superclass of the class defined by the class file.
	SuperClass uint16
	// Interfaces is a list of indices of CONSTANT_Class_info structures in the constant pool.
	// Each index represents an interface implemented by the class.
	Interfaces []uint16
	// Fields is a list of fields of the class.
	// Each field contains the access flags, name, descriptor and attributes of the field.
	Fields []FieldInfo
	// Methods is a list of methods of the class.
	Methods []MethodInfo
	// Attributes is a list of attributes of the class.
	Attributes []AttributeInfo
}

var (
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

// Validate validates the class file
// It checks if the magic number is correct, if the access flags are valid, if the this_class and super_class, etc...
// are valid.
//
// It is really heavy and should be used only for debugging purposes, as it is not necessary to fully validate a class,
// since you need to go through the whole class hierarchy to fully validate it.
func (c *ClassFile) Validate() error {
	if err := c.ValidateMagicNumber(); err != nil {
		return err
	}

	if err := c.ValidateAccessFlags(); err != nil {
		return err
	}

	if err := c.ValidateThisClass(); err != nil {
		return err
	}

	// TODO: implement
	return nil
}

func (c *ClassFile) ValidateMagicNumber() error {
	if c.Magic != MagicNumber {
		return fmt.Errorf("invalid magic number: 0x%x", c.Magic)
	}

	return nil
}

func (c *ClassFile) ValidateAccessFlags() error {
	validFlags := ACC_PUBLIC | ACC_FINAL | ACC_SUPER | ACC_INTERFACE | ACC_ABSTRACT | ACC_SYNTHETIC | ACC_ANNOTATION | ACC_ENUM
	if c.AccessFlags&^validFlags != 0 {
		return fmt.Errorf("invalid access flags: 0x%x", c.AccessFlags)
	}

	// If the ACC_INTERFACE flag is set, the ACC_ABSTRACT flag must also be set, and the ACC_FINAL, ACC_SUPER,
	// and ACC_ENUM flags set must not be set.
	if c.AccessFlags&ACC_INTERFACE != 0 {
		if c.AccessFlags&ACC_ABSTRACT == 0 {
			return fmt.Errorf("interface must have abstract flag set")
		}

		if c.AccessFlags&ACC_FINAL != 0 {
			return fmt.Errorf("interface must not have final flag set")
		}

		if c.AccessFlags&ACC_SUPER != 0 {
			return fmt.Errorf("interface must not have super flag set")
		}

		if c.AccessFlags&ACC_ENUM != 0 {
			return fmt.Errorf("interface must not have enum flag set")
		}
	}

	// If the ACC_INTERFACE flag is not set, any of the other flags in Table 4.1-A may be set except ACC_ANNOTATION.
	// However, such a class file must not have both its ACC_FINAL and ACC_ABSTRACT flags set.
	if c.AccessFlags&ACC_INTERFACE == 0 {
		if c.AccessFlags&ACC_FINAL != 0 && c.AccessFlags&ACC_ABSTRACT != 0 {
			return fmt.Errorf("class must not have both final and abstract flags set")
		}

		if c.AccessFlags&ACC_ANNOTATION != 0 {
			return fmt.Errorf("class must not have annotation flag set")
		}
	}

	return nil
}

func (c *ClassFile) ValidateThisClass() error {
	// ThisClass should be a valid ConstantPool index and it should always be of tag CONSTANT_Class_info
	if c.ThisClass == 0 || c.ThisClass > uint16(len(c.ConstantPool)) {
		return fmt.Errorf("invalid this class index: %d", c.ThisClass)
	}

	if c.ConstantPool[c.ThisClass-1].Tag != CONSTANT_Class {
		return fmt.Errorf("this class should be a CONSTANT_Class_info")
	}

	return nil
}

func ExecuteClassFile(reader *bytes.Reader) error {
	_, err := ClassFileFromReader(reader)

	if err != nil {
		return err
	}

	return nil
}
