package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/bigwhite/instrument_trace/instrumenter"
	"github.com/bigwhite/instrument_trace/instrumenter/ast"
)

var (
	wrote bool
)

func init() {
	flag.BoolVar(&wrote, "w", false, "write result to (source) file instead of stdout")
}

func usage() {
	fmt.Println("instrument [-w] xxx.go")
	flag.PrintDefaults()
}

func main() {
	fmt.Println(os.Args)
	flag.Usage = usage
	flag.Parse()

	if len(os.Args) < 2 {
		usage()
		return
	}

	var file string
	if len(os.Args) == 3 {
		file = os.Args[2]
	}

	if len(os.Args) == 2 {
		file = os.Args[1]
	}
	if filepath.Ext(file) != ".go" {
		usage()
		return
	}

	var ins instrumenter.Instrumenter
	ins = ast.New("github.com/bigwhite/instrument_trace", "trace", "Trace")
	newSrc, err := ins.Instrument(file)
	if err != nil {
		panic(err)
	}

	if newSrc == nil {
		// add nothing to the source file. no change
		fmt.Printf("no trace added for %s\n", file)
		return
	}

	if !wrote {
		fmt.Println(string(newSrc))
		return
	}

	// write to the source file
	if err = ioutil.WriteFile(file, newSrc, 0666); err != nil {
		fmt.Printf("write %s error: %v\n", file, err)
		return
	}
	fmt.Printf("instrument trace for %s ok\n", file)
}
