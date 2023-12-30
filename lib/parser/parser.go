package parser

import (
	"strings"

	"github.com/lspaccatrosi16/go-cli-tools/logging"
	"github.com/lspaccatrosi16/ts-go/lib/types"
)

type runtime struct {
	lines  []string
	curPos int
}

func (r *runtime) getNextLine() (string, bool) {
	r.curPos++
	return r.getCurLine()
}

func (r *runtime) getCurLine() (string, bool) {
	if r.curPos >= len(r.lines) {
		return "", false
	} else {
		return r.lines[r.curPos], true
	}
}

func ParseInput(input string) []*types.TsNode {
	logger := logging.GetLogger()

	lines := strings.Split(input, "\n")
	rt := &runtime{lines: lines, curPos: 0}

	nodes := []*types.TsNode{}

	for {
		next, ok := rt.getCurLine()
		if !ok {
			break
		}

		elements := strings.Split(next, " ")

		if elements[0] == "export" {
			elements = elements[1:]
		}

		var node *types.TsNode
		var name string
		var failed bool

		if len(elements) < 2 {
			failed = true
		} else {
			term := elements[0]
			name = elements[1]

			switch term {
			case "interface":
				node = parseObject(rt, name)
			case "type":
				node = parseType(rt, name, strings.Join(elements[3:], " "))
				rt.getNextLine()
			default:
				failed = true
			}
		}

		if failed {
			logger.Debug("WARN: could not parse line: '" + next + "'")
			rt.getNextLine()
		} else {
			types.RegisterCustomType(name, node)
			nodes = append(nodes, node)
		}
	}
	return nodes
}

func parseType(rt *runtime, name, typeInfo string) *types.TsNode {
	if strings.HasPrefix(typeInfo, "{") {
		return parseObject(rt, name)
	} else {
		ni := types.NodeInfo{
			JsonName:  name,
			FieldName: makeGoName(name),
		}

		field := types.TsField{
			NodeInfo: ni,
		}

		for _, c := range typeInfo {
			if c == ';' {
				break
			}

			field.TypeData += string(c)
		}

		node := types.TsNode{
			NodeInfo: ni,
			Fields:   []*types.TsField{&field},
			Type:     types.Inline,
		}
		return &node
	}
}

func parseObject(rt *runtime, name string) *types.TsNode {
	node := types.TsNode{
		NodeInfo: types.NodeInfo{
			JsonName:  name,
			FieldName: makeGoName(name),
		},
		Fields: []*types.TsField{},
		Type:   types.Object,
	}

	for {
		next, ok := rt.getNextLine()
		if !ok {
			break
		}

		if next == "" {
			continue
		}

		if strings.Contains(next, "}") {
			break
		}

		node.Fields = append(node.Fields, parseObjectField(rt))
	}

	return &node
}

func parseObjectField(rt *runtime) *types.TsField {
	cur, ok := rt.getCurLine()
	if !ok {
		return nil
	}

	components := strings.Split(cur, ":")
	for i := 0; i < len(components); i++ {
		components[i] = strings.Trim(components[i], " \n\r\t")
	}

	name := components[0]
	typeInfo := strings.Join(components[1:], ":")

	t := parseType(rt, name, typeInfo)
	return t.Fields[0]
}

func makeGoName(str string) string {
	if len(str) == 0 {
		return ""
	}
	return strings.ToUpper(string(str[0])) + string(str[1:])
}
