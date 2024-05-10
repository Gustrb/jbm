package classfile_test

import (
	"bytes"
	"testing"

	"github.com/Gustrb/jbm/src/core"
	"github.com/Gustrb/jbm/src/utils"
)

var ClazzWithAttributeClassFile, _ = utils.ReadFileContent("../fixtures/ClazzWithAttribute.class")

func TestShouldParseTheAttributesOfTheClass(t *testing.T) {
	cf, err := core.ClassFileFromReader(bytes.NewReader(ClazzWithAttributeClassFile))

	if err != nil {
		t.Fatalf("Error reading class file: %s", err)
	}

	if len(cf.ConstantPool) != 14 {
		t.Fatalf("Expected constant pool of size 1, got %d", len(cf.ConstantPool))
	}

	expected := [14]string{
		"ConstantPoolIndexableInfo{ ClassIndex: 2, NameAndTypeIndex: 3 }",
		"ClassInfo{ NameIndex: 4 }",
		"NameAndTypeInfo{ NameIndex: 5, DescriptorIndex: 6 }",
		"UTF8Info{ Bytes: java/lang/Object }",
		"UTF8Info{ Bytes: <init> }",
		"UTF8Info{ Bytes: ()V }",
		"ClassInfo{ NameIndex: 8 }",
		"UTF8Info{ Bytes: ClazzWithAttribute }",
		"UTF8Info{ Bytes: attr }",
		"UTF8Info{ Bytes: Ljava/lang/String; }",
		"UTF8Info{ Bytes: Code }",
		"UTF8Info{ Bytes: LineNumberTable }",
		"UTF8Info{ Bytes: SourceFile }",
		"UTF8Info{ Bytes: ClazzWithAttribute.java }",
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
		t.Fatalf("Expected this class to be 7, got %d", cf.ThisClass)
	}

	if cf.SuperClass != 2 {
		t.Fatalf("Expected super class to be 2, got %d", cf.SuperClass)
	}

	if len(cf.Interfaces) != 0 {
		t.Fatalf("Expected interfaces to be empty, got %d", len(cf.Interfaces))
	}

	if len(cf.Fields) != 1 {
		t.Fatalf("Expected fields to have 1 element, got %d", len(cf.Fields))
	}

	if cf.Fields[0].AccessFlags != 1 {
		t.Fatalf("Expected field access flags to be 1, got %d", cf.Fields[0].AccessFlags)
	}

	if cf.Fields[0].NameIndex != 9 {
		t.Fatalf("Expected field name index to be 9, got %d", cf.Fields[0].NameIndex)
	}

	if cf.Fields[0].DescriptorIndex != 10 {
		t.Fatalf("Expected field descriptor index to be 10, got %d", cf.Fields[0].DescriptorIndex)
	}

	if len(cf.Fields[0].Attributes) != 0 {
		t.Fatalf("Expected field attributes to be empty, got %d", len(cf.Fields[0].Attributes))
	}
}
