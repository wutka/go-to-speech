package main

import (
	"flag"
	"go-to-speech/pkg"
)

func main() {
	quietFlag := flag.Bool("q", false, "Don't output speech")
	skipImportsFlag := flag.Bool("noimports", false, "Don't read imports")

	flag.Parse()

	pkg.ShutUp = *quietFlag
	pkg.SkipImports = *skipImportsFlag

	for _, filename := range flag.Args() {
		pkg.SpeakGoFile(filename)
	}
}
