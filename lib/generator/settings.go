package generator

import "github.com/lspaccatrosi16/go-cli-tools/input"

type GenerateSettings struct {
	NumberType  string
	StringType  string
	PackageName string
	EnumType    string
}

func NewSettings(defaults bool) GenerateSettings {
	var nt, st, pn, et string
	if defaults {
		nt = "float64"
		st = "string"
		pn = "types"
		et = "int"
	} else {
		nt = input.GetInput("Go Number Type")
		st = input.GetInput("Go String Type")
		et = input.GetInput("Go Enum Type")
		pn = input.GetInput("Go Package Name")
	}

	return GenerateSettings{
		NumberType:  nt,
		StringType:  st,
		PackageName: pn,
		EnumType:    et,
	}
}
