package core_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/Gustrb/jbm/src/core"
)

func TestShouldValidateMagicNumber(t *testing.T) {
	cf := core.ClassFile{Magic: 0xCAFEBABE}

	if err := cf.ValidateMagicNumber(); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestShouldFailIfMagicNumberIsInvalid(t *testing.T) {
	cf := core.ClassFile{Magic: 0xDEADBEEF}

	if err := cf.ValidateMagicNumber(); err.Error() != "invalid magic number: 0xdeadbeef" {
		t.Errorf("Expected 'invalid magic number 0xdeadbeef', got %v", err)
	}
}

func TestShouldNotAllowRandomAccessFlags(t *testing.T) {
	numbersToTest := []uint16{69, 420, 1337, 9001}
	cf := core.ClassFile{}

	for _, number := range numbersToTest {
		cf.AccessFlags = number

		if err := cf.ValidateAccessFlags(); err.Error() != fmt.Sprintf("invalid access flags: 0x%x", number) {
			t.Errorf("Expected 'invalid access flags: 0x%x', got %v", number, err)
		}
	}
}

func TestItShouldNotAllowAccessFlagsToBeInterfaceAndNotAbstract(t *testing.T) {
	cf := core.ClassFile{AccessFlags: core.ACC_INTERFACE}

	if err := cf.ValidateAccessFlags(); err.Error() != "interface must have abstract flag set" {
		t.Errorf("Expected 'interface must have abstract flag set', got %v", err)
	}

	cf.AccessFlags = core.ACC_INTERFACE | core.ACC_ABSTRACT

	if err := cf.ValidateAccessFlags(); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestItShouldNotAllowAccessFlagsToBeInterfaceAndFinal(t *testing.T) {
	cf := core.ClassFile{AccessFlags: core.ACC_INTERFACE | core.ACC_ABSTRACT | core.ACC_FINAL}

	if err := cf.ValidateAccessFlags(); err.Error() != "interface must not have final flag set" {
		t.Errorf("Expected 'interface must not have final flag set', got %v", err)
	}
}

func TestItShouldNotAllowAccessFlagsToBeInterfaceAndSuper(t *testing.T) {
	cf := core.ClassFile{AccessFlags: core.ACC_INTERFACE | core.ACC_ABSTRACT | core.ACC_SUPER}

	if err := cf.ValidateAccessFlags(); err.Error() != "interface must not have super flag set" {
		t.Errorf("Expected 'interface must not have super flag set', got %v", err)
	}
}

func TestItShouldNotAllowAccessFlagsToBeInterfaceAndEnum(t *testing.T) {
	cf := core.ClassFile{AccessFlags: core.ACC_INTERFACE | core.ACC_ABSTRACT | core.ACC_ENUM}

	if err := cf.ValidateAccessFlags(); err.Error() != "interface must not have enum flag set" {
		t.Errorf("Expected 'interface must not have enum flag set', got %v", err)
	}
}

func TestItShouldNotAllowAccessFlagsToBeAnnotationAndNotInterface(t *testing.T) {
	cf := core.ClassFile{AccessFlags: core.ACC_ANNOTATION}

	if err := cf.ValidateAccessFlags(); err.Error() != "class must not have annotation flag set" {
		t.Errorf("Expected 'class must not have annotation flag set', got %v", err)
	}
}

func TestItShouldNotAllowAccessFlagsToBeFinalAndAbstract(t *testing.T) {
	cf := core.ClassFile{AccessFlags: core.ACC_FINAL | core.ACC_ABSTRACT}

	if err := cf.ValidateAccessFlags(); err.Error() != "class must not have both final and abstract flags set" {
		t.Errorf("Expected 'class must not have both final and abstract flags set', got %v", err)
	}
}

func TestItShouldAllowValidAccessFlags(t *testing.T) {
	cf := core.ClassFile{AccessFlags: core.ACC_PUBLIC}

	if err := cf.ValidateAccessFlags(); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	cf.AccessFlags = core.ACC_PUBLIC | core.ACC_FINAL
	if err := cf.ValidateAccessFlags(); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	cf.AccessFlags = core.ACC_PUBLIC | core.ACC_SUPER
	if err := cf.ValidateAccessFlags(); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestShouldNotAllowInvalidThisClass(t *testing.T) {
	cf := core.ClassFile{ThisClass: 0}

	if err := cf.ValidateThisClass(); err.Error() != "invalid this class index: 0" {
		t.Errorf("Expected 'invalid this class index: 0', got %v", err)
	}

	cf.ThisClass = 1
	cf.ConstantPool = make([]core.ConstantPoolInfo, 1)
	cf.ConstantPool[0] = core.ConstantPoolInfo{Tag: core.CONSTANT_Utf8}

	if err := cf.ValidateThisClass(); err.Error() != "this class should be a CONSTANT_Class_info" {
		t.Errorf("Expected 'this class should be a CONSTANT_Class_info', got %v", err)
	}
}

func TestShouldAllowValidThisClass(t *testing.T) {
	cf := core.ClassFile{ThisClass: 1, ConstantPool: make([]core.ConstantPoolInfo, 1)}
	cf.ConstantPool[0] = core.ConstantPoolInfo{Tag: core.CONSTANT_Class}

	if err := cf.ValidateThisClass(); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestItShouldReadTheMinorAndMajorVersions(t *testing.T) {
	b := []byte{0xCA, 0xFE, 0xBA, 0xBE, 0x00, 0x00, 0x00, 0x34}

	// suppress the error, since we don't want to put a whole
	// fake bytecode stub here
	cf, _ := core.ClassFileFromReader(bytes.NewReader(b))

	if cf.MinorVersion != 0 {
		t.Errorf("Expected minor version 0, got %d", cf.MinorVersion)
	}

	if cf.MajorVersion != 52 {
		t.Errorf("Expected major version 52, got %d", cf.MajorVersion)
	}
}

func TestItShouldFailIfTheConstantPoolSizeIsNotValid(t *testing.T) {
	b := []byte{0xCA, 0xFE, 0xBA, 0xBE, 0x00, 0x00, 0x00, 0x34, 0x00, 0x01}

	_, err := core.ClassFileFromReader(bytes.NewReader(b))
	if err != core.ErrInvalidConstantPoolSize {
		t.Errorf("Expected ErrInvalidConstantPoolSize, got %v", err)
	}

	b = []byte{0xCA, 0xFE, 0xBA, 0xBE, 0x00, 0x00, 0x00, 0x34, 0x00, 0x00}
	_, err = core.ClassFileFromReader(bytes.NewReader(b))

	if err != core.ErrInvalidConstantPoolSize {
		t.Errorf("Expected ErrInvalidConstantPoolSize, got %v", err)
	}
}
