package classfile_test

import (
	"bytes"
	"testing"

	"github.com/Gustrb/jbm/src/core"
	"github.com/Gustrb/jbm/src/utils"
)

var EmptyClassFile, _ = utils.ReadFileContent("../fixtures/Empty.class")

func TestShouldAllowAValidJavaProgram(t *testing.T) {
	cf, err := core.ClassFileFromReader(bytes.NewReader(EmptyClassFile))

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
		"UTF8Info{ Bytes: Empty }",
		"UTF8Info{ Bytes: Code }",
		"UTF8Info{ Bytes: LineNumberTable }",
		"UTF8Info{ Bytes: main }",
		"UTF8Info{ Bytes: ([Ljava/lang/String;)V }",
		"UTF8Info{ Bytes: SourceFile }",
		"UTF8Info{ Bytes: Empty.java }",
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
}
