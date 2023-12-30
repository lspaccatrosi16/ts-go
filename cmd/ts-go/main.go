package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/lspaccatrosi16/go-cli-tools/logging"
	"github.com/lspaccatrosi16/ts-go/lib/generator"
	"github.com/lspaccatrosi16/ts-go/lib/parser"
	"github.com/lspaccatrosi16/ts-go/lib/util"
)

var fp = flag.String("i", "", "Path to the input file")
var op = flag.String("o", "", "Path to the output file (otherwise same as input)")
var verbose = flag.Bool("v", false, "Display Verbose Logging")
var help = flag.Bool("h", false, "Shows the help message")
var showXml = flag.Bool("x", false, "Print intermediate XML representation of type system")
var useDefaults = flag.Bool("d", false, "Use default go types")

func main() {
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	logger := logging.GetLogger()
	logger.SetVerbose(*verbose)

	if *fp == "" {
		fmt.Println("must provide a file path")
		os.Exit(1)
	}

	if *op == "" {
		ex := filepath.Ext(*fp)
		*op = (*fp)[:len(*fp)-len(ex)] + ".go"
	}

	src, err := os.Open(*fp)
	if err != nil {
		panic(err)
	}

	defer src.Close()

	buf := bytes.NewBuffer(nil)

	io.Copy(buf, src)

	tree := parser.ParseInput(buf.String())

	if *showXml {
		logger.Log(util.FormatIr(tree))
	}

	settings := generator.NewSettings(*useDefaults)

	dst, err := os.Create(*op)
	if err != nil {
		panic(err)
	}

	defer dst.Close()

	generated := generator.Generate(settings, tree)
	io.Copy(dst, generated)
}
