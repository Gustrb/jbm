package classfile_test

import (
	"bytes"
	"testing"

	"github.com/Gustrb/jbm/src/core"
	"github.com/Gustrb/jbm/src/utils"
)

var ImplEmptyInterfaceClassFile, _ = utils.ReadFileContent("../fixtures/ImplEmptyInterface.class")

func TestShouldAllowClassFileImplementingAnEmptyInterface(t *testing.T) {
	cf, err := core.ClassFileFromReader(bytes.NewReader(ImplEmptyInterfaceClassFile))

	if err != nil {
		t.Fatalf("Error reading class file: %s", err)
	}

	if len(cf.ConstantPool) != 14 {
		t.Fatalf("Expected constant pool of size 8, got %d", len(cf.ConstantPool))
	}

	expected := [14]string{
		"ConstantPoolIndexableInfo{ ClassIndex: 2, NameAndTypeIndex: 3 }",
		"ClassInfo{ NameIndex: 4 }",
		"NameAndTypeInfo{ NameIndex: 5, DescriptorIndex: 6 }",
		"UTF8Info{ Bytes: java/lang/Object }",
		"UTF8Info{ Bytes: <init> }",
		"UTF8Info{ Bytes: ()V }",
		"ClassInfo{ NameIndex: 8 }",
		"UTF8Info{ Bytes: ImplEmptyInterface }",
		"ClassInfo{ NameIndex: 10 }",
		"UTF8Info{ Bytes: EmptyInterface }",
		"UTF8Info{ Bytes: Code }",
		"UTF8Info{ Bytes: LineNumberTable }",
		"UTF8Info{ Bytes: SourceFile }",
		"UTF8Info{ Bytes: ImplEmptyInterface.java }",
	}

	for i, e := range expected {
		if cf.ConstantPool[i].String() != e {
			t.Fatalf("Expected constant pool entry %d to be %s, got %s", i, e, cf.ConstantPool[i].String())
		}
	}

	if cf.AccessFlags != 33 {
		t.Fatalf("Expected access flags to be 33, got %d", cf.AccessFlags)
	}

	if cf.ThisClass != 7 {
		t.Fatalf("Expected this class to be 1, got %d", cf.ThisClass)
	}

	if cf.SuperClass != 2 {
		t.Fatalf("Expected super class to be 3, got %d", cf.SuperClass)
	}

	if len(cf.Interfaces) != 1 {
		t.Fatalf("Expected interfaces to be 1, got %d", len(cf.Interfaces))
	}

	if cf.Interfaces[0] != 9 {
		t.Fatalf("Expected interface to be 9, got %d", cf.Interfaces[0])
	}
}
