package main

import (
	"customanalyzer/checkpubfuncname" // Adjust import path to your module

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(checkpubfuncname.Analyzer)
}
