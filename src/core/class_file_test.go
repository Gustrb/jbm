package core_test

import (
	"bytes"
	"testing"

	"github.com/Gustrb/jbm/src/core"
)

func TestShouldValidateMagicNumber(t *testing.T) {
	cf := core.ClassFile{Magic: 0xCAFEBABE}

	if err := cf.Validate(); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestShouldFailIfMagicNumberIsInvalid(t *testing.T) {
	cf := core.ClassFile{Magic: 0xDEADBEEF}

	if err := cf.Validate(); err.Error() != "invalid magic number: 0xdeadbeef" {
		t.Errorf("Expected 'invalid magic number 0xdeadbeef', got %v", err)
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
