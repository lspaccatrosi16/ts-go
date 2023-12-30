package types

import (
	"golang.org/x/exp/slices"
)

var custTypes = map[string]*TsNode{}

func RegisterCustomType(name string, t *TsNode) {
	custTypes[name] = t
}

func RetrieveCustomType(name string) (*TsNode, bool) {
	v, ok := custTypes[name]
	return v, ok
}

func AllCustomTypes() []string {
	arr := []string{}

	for k := range custTypes {
		arr = append(arr, k)
	}

	slices.Sort(arr)

	return arr
}
