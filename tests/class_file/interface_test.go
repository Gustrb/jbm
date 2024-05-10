package classfile_test

import (
	"bytes"
	"testing"

	"github.com/Gustrb/jbm/src/core"
	"github.com/Gustrb/jbm/src/utils"
)

var InterfaceClassFile, _ = utils.ReadFileContent("../fixtures/Interface.class")

func TestShouldAllowValidClassFileWithInterfaceDefinition(t *testing.T) {
	cf, err := core.ClassFileFromReader(bytes.NewReader(InterfaceClassFile))

	if err != nil {
		t.Fatalf("Error reading class file: %s", err)
	}

	if len(cf.ConstantPool) != 8 {
		t.Fatalf("Expected constant pool of size 8, got %d", len(cf.ConstantPool))
	}

	expected := [8]string{
		"ClassInfo{ NameIndex: 2 }",
		"UTF8Info{ Bytes: Interface }",
		"ClassInfo{ NameIndex: 4 }",
		"UTF8Info{ Bytes: java/lang/Object }",
		"UTF8Info{ Bytes: testMethod }",
		"UTF8Info{ Bytes: ()V }",
		"UTF8Info{ Bytes: SourceFile }",
		"UTF8Info{ Bytes: Interface.java }",
	}

	for i, e := range expected {
		if cf.ConstantPool[i].String() != e {
			t.Fatalf("Expected constant pool entry %d to be %s, got %s", i, e, cf.ConstantPool[i].String())
		}
	}

	if cf.AccessFlags != 1537 {
		t.Fatalf("Expected access flags to be 1537, got %d", cf.AccessFlags)
	}

	if cf.ThisClass != 1 {
		t.Fatalf("Expected this class to be 1, got %d", cf.ThisClass)
	}

	if cf.SuperClass != 3 {
		t.Fatalf("Expected super class to be 2, got %d", cf.SuperClass)
	}

	// We don't care implement any interfaces
	if len(cf.Interfaces) != 0 {
		t.Fatalf("Expected interfaces to be 0, got %d", len(cf.Interfaces))
	}

	if len(cf.Fields) != 0 {
		t.Fatalf("Expected fields to be empty, got %d", len(cf.Fields))
	}

	if len(cf.Methods) != 1 {
		t.Fatalf("Expected methods to be 1, got %d", len(cf.Methods))
	}

	if cf.Methods[0].NameIndex != 5 {
		t.Fatalf("Expected method name index to be 5, got %d", cf.Methods[0].NameIndex)
	}

	if cf.Methods[0].DescriptorIndex != 6 {
		t.Fatalf("Expected method descriptor index to be 6, got %d", cf.Methods[0].DescriptorIndex)
	}

	if len(cf.Attributes) != 1 {
		t.Fatalf("Expected attributes to be 1, got %d", len(cf.Attributes))
	}
}
