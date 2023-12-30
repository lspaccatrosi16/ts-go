package util

import (
	"bytes"
	"fmt"

	"github.com/go-xmlfmt/xmlfmt"
	"github.com/lspaccatrosi16/ts-go/lib/types"
)

func FormatIr(arr []*types.TsNode) string {
	buf := bytes.NewBuffer(nil)

	for _, el := range arr {
		unformatted := el.String()

		formatted := xmlfmt.FormatXML(unformatted, "", "   ")
		fmt.Fprintln(buf, formatted)
	}
	return buf.String()
}
