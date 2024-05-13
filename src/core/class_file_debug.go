package core

import "fmt"

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
