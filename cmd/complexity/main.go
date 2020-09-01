package main

import (
	"github.com/shoooooman/complexity"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(complexity.Analyzer) }

