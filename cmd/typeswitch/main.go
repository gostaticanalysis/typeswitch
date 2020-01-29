package main

import (
	"github.com/gostaticanalysis/typeswitch"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(typeswitch.Analyzer) }
