package generator

import (
	"bytes"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/lspaccatrosi16/ts-go/lib/types"
)

type Generator struct {
	Nodes          []*types.TsNode
	FlattenedNodes []*types.TsNode
	Settings       GenerateSettings
	EnumTypes      map[string]*types.GoEnum
	StructTypes    map[string]*types.GoStruct
	Types          map[string]*types.GoType
	EnumOrder      []string
	StructOrder    []string
	TypeOrder      []string
	SymTab         map[string]string
}

func (g *Generator) ProduceFile() io.Reader {
	buf := bytes.NewBuffer(nil)

	fmt.Fprintf(buf, "package %s\n", g.Settings.PackageName)

	slices.Sort(g.TypeOrder)

	for _, t := range g.TypeOrder {
		fmt.Fprintln(buf, g.Types[t].Code())
	}

	slices.Sort(g.EnumOrder)

	for _, e := range g.EnumOrder {
		fmt.Fprintln(buf, g.EnumTypes[e].Code())
	}

	slices.Sort(g.StructOrder)

	for _, s := range g.StructOrder {
		fmt.Fprintln(buf, g.StructTypes[s].Code())
	}

	return buf
}

func (g *Generator) Analyse() {
	for _, node := range g.Nodes {
		g.FlattenNode(node)
	}
}

func (g *Generator) ProduceTypes() {
	for _, node := range g.FlattenedNodes {
		switch node.Type {
		case types.Object:
			str := types.GoStruct{Name: g.ResolveSymbol(node.FieldName)}
			for _, f := range node.Fields {
				str.Fields = append(str.Fields, types.GoStructField{
					Name:    formatSym(f.FieldName),
					Type:    g.ParseTypeData(f.TypeData, false),
					JsonTag: f.JsonName,
				})
			}
			g.StructTypes[node.FieldName] = &str
			g.StructOrder = append(g.StructOrder, node.FieldName)
		case types.Inline:
			field := node.Fields[0]
			if strings.Contains(field.TypeData, "|") {
				elements := strings.Split(field.TypeData, "|")
				values := []types.GoEnumVal{}
				foundType := ""

				for i, e := range elements {
					elements[i] = strings.Trim(e, " \n\t\r")
					if elements[i] == "" {
						continue
					}

					values = append(values, types.GoEnumVal{
						Ident: fmt.Sprintf("%s_%s", g.ResolveSymbol(node.FieldName), cleanSym(elements[i])),
						Value: elements[i],
					})

					if foundType == "" {
						foundType = inferTypeFromElement(elements[i])
					} else {
						if foundType != inferTypeFromElement(elements[i]) {
							foundType = "INVALID_TYPE"
							break
						}
					}
				}

				enu := types.GoEnum{
					Name:    g.ResolveSymbol(node.FieldName),
					VarType: g.Settings.EnumType,
					ValType: g.ParseTypeData(foundType, true),
					Values:  values,
				}

				g.EnumTypes[node.FieldName] = &enu
				g.EnumOrder = append(g.EnumOrder, node.FieldName)
			} else {
				typ := types.GoType{
					Name:    g.ResolveSymbol(node.FieldName),
					VarType: g.ParseTypeData(field.TypeData, false),
				}

				g.Types[node.FieldName] = &typ
				g.TypeOrder = append(g.TypeOrder, node.FieldName)
			}
		}

	}
}

func (g *Generator) ParseTypeData(str string, basicOnly bool) string {
	str = strings.Trim(str, " \n\r\t")

	switch str {
	case "number":
		return g.Settings.NumberType
	case "string":
		return g.Settings.StringType
	case "boolean":
		return "bool"
	default:
		if basicOnly {
			return "\"INVALID TYPE\""
		}
	}

	if strings.HasSuffix(str, "[]") {
		return "[]" + g.ParseTypeData(str[:len(str)-2], false)
	}

	if strings.HasPrefix(str, "Array<") {
		return "[]" + g.ParseTypeData(str[len("Array<"):len(str)-1], false)
	}

	if strings.HasPrefix(str, "Record<") {
		str = str[len("Record<") : len(str)-1]
		components := strings.Split(str, ",")
		if len(components) != 2 {
			return "INVALID_TYPE"
		}
		key := g.ParseTypeData(components[0], false)
		val := g.ParseTypeData(components[1], false)

		return fmt.Sprintf("map[%s]%s", key, val)
	}

	return str
}

func (g *Generator) FlattenNode(node *types.TsNode) {
	if node == nil {
		return
	}
	g.FlattenedNodes = append(g.FlattenedNodes, node)
	for _, f := range node.Fields {
		g.FlattenNode(f.SubType)
	}
}

func (g *Generator) ResolveSymbol(str string) string {
	if v, ok := g.SymTab[str]; ok {
		return v
	}

	resolved := formatSym(str)

	g.SymTab[str] = resolved
	return resolved
}

func Generate(settings GenerateSettings, nodes []*types.TsNode) io.Reader {
	generator := Generator{
		Nodes:          nodes,
		FlattenedNodes: []*types.TsNode{},
		Settings:       settings,
		EnumTypes:      map[string]*types.GoEnum{},
		StructTypes:    map[string]*types.GoStruct{},
		Types:          map[string]*types.GoType{},
		EnumOrder:      []string{},
		StructOrder:    []string{},
		TypeOrder:      []string{},
		SymTab:         map[string]string{},
	}

	generator.Analyse()
	generator.ProduceTypes()

	return generator.ProduceFile()
}

func inferTypeFromElement(str string) string {
	if str == "true" || str == "false" {
		return "boolean"
	} else if _, err := strconv.ParseFloat(str, 64); err == nil {
		return "number"
	} else if str[0] == '\'' || str[0] == '"' {
		return "string"
	} else {
		return "INVALID_TYPE"
	}
}

func formatSym(str string) string {
	var resolved string
	if (str[0] >= 'a' && str[0] <= 'z') || (str[0] >= 'A' && str[0] <= 'Z') {
		resolved = str
	} else if str[0] == '\'' || str[len(str)-1] == '"' {
		resolved = str[1 : len(str)-1]
	} else {
		resolved = "V_" + str
	}
	return resolved
}

func cleanSym(str string) string {
	if str[0] == '\'' || str[0] == '"' {
		return str[1 : len(str)-1]
	}
	return str
}
