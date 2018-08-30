package main

import (
	"flag"
	"go-to-speech/pkg"
)

func main() {
	var foo [10]int
	quietFlag := flag.Bool("q", false, "Don't output speech")

	if foo[0] > 0 {

	}
	flag.Parse()

	pkg.ShutUp = *quietFlag

	for _, filename := range flag.Args() {
		pkg.SpeakGoFile(filename)
	}
}
