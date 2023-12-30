package types

import (
	"bytes"
	"fmt"
)

type GoEnumVal struct {
	Ident string
	Value string
}

type GoEnum struct {
	Name    string
	VarType string
	ValType string
	Values  []GoEnumVal
}

func (e *GoEnum) Code() string {
	buf := bytes.NewBuffer(nil)
	fmt.Fprintf(buf, "type %s %s\n", e.Name, e.VarType)

	fmt.Fprintf(buf, "const (\n")
	for i := 0; i < len(e.Values); i++ {
		str := e.Values[i].Ident
		if i == 0 {
			str += fmt.Sprintf(" %s = iota", e.Name)
		}
		fmt.Fprintln(buf, "\t"+str)
	}
	fmt.Fprintln(buf, ")")

	fmt.Fprintf(buf, "\nfunc(v %s) Val() %s {\n\tswitch v {\n", e.Name, e.ValType)
	for i := 0; i < len(e.Values); i++ {
		fmt.Fprintf(buf, "\tcase %s:\n\t\treturn %s\n", e.Values[i].Ident, e.Values[i].Value)
	}
	fmt.Fprintln(buf, "\tdefault:\n\t\treturn \"Unknown Value\"")
	fmt.Fprintln(buf, "\t}\n}")

	return buf.String()
}

type GoStructField struct {
	Name    string
	Type    string
	JsonTag string
}

type GoStruct struct {
	Name   string
	Fields []GoStructField
}

func (s *GoStruct) Code() string {
	buf := bytes.NewBuffer(nil)

	fmt.Fprintf(buf, "type %s struct {\n", s.Name)

	for i := 0; i < len(s.Fields); i++ {
		fmt.Fprintf(buf, "\t%s %s `json:\"%s\"`\n", s.Fields[i].Name, s.Fields[i].Type, s.Fields[i].JsonTag)
	}

	fmt.Fprintf(buf, "}")

	return buf.String()
}

type GoType struct {
	Name    string
	VarType string
}

func (t *GoType) Code() string {
	return fmt.Sprintf("type %s %s", t.Name, t.VarType)
}
