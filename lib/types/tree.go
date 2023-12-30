package types

import (
	"bytes"
	"fmt"
)

type NodeType int

const (
	Object NodeType = iota
	Inline
)

func (n NodeType) String() string {
	switch n {
	case Object:
		return "object"
	case Inline:
		return "inline"
	default:
		return "Invalid Type"
	}
}

type NodeInfo struct {
	JsonName  string
	FieldName string
}

func (n *NodeInfo) XmlOpen(t string) string {
	if t != "" {
		return fmt.Sprintf("<%s json=%s type=%s>", n.FieldName, n.JsonName, t)
	} else {
		return fmt.Sprintf("<%s json=%s>", n.FieldName, n.JsonName)
	}
}

func (n *NodeInfo) XmlClose() string {
	return fmt.Sprintf("</%s>", n.FieldName)
}

type TsField struct {
	NodeInfo
	SubType  *TsNode
	TypeData string
}

func (f *TsField) String() string {
	buf := bytes.NewBuffer(nil)

	fmt.Fprint(buf, f.XmlOpen(""))

	if f.SubType != nil {
		fmt.Fprint(buf, f.SubType.String())
	} else {
		fmt.Fprint(buf, f.TypeData)
	}

	fmt.Fprint(buf, f.XmlClose())

	return buf.String()
}

type TsNode struct {
	NodeInfo
	Type   NodeType
	Fields []*TsField
}

func (n *TsNode) String() string {
	buf := bytes.NewBuffer(nil)

	fmt.Fprint(buf, n.XmlOpen(n.Type.String()))

	for _, f := range n.Fields {
		fmt.Fprint(buf, f.String())
	}

	fmt.Fprint(buf, n.XmlClose())
	return buf.String()
}
