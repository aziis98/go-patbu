package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/aziis98/go-patbu"

	"github.com/spf13/pflag"
)

var stdin bool

func main() {
	fs := &pflag.FlagSet{}
	fs.BoolVar(&stdin, "stdin", false,
		`map value from stdin`)

	err := fs.Parse(os.Args[1:])
	if err != nil {
		if err != pflag.ErrHelp {
			log.Fatal(err)
		}
		os.Exit(2)
	}

	var pattern, builder patbu.Patbu

	if fs.NArg() < 2 {
		log.Fatalf(`expected at least two arguments, got %v`, fs.NArg())
	}

	pattern, err = patbu.Parse(fs.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	builder, err = patbu.Parse(fs.Arg(1))
	if err != nil {
		log.Fatal(err)
	}

	var in string
	if stdin {
		if fs.NArg() != 2 {
			log.Fatalf(`expected two arguments in stdin mode, got %v`, fs.NArg())
		}

		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}

		in = strings.TrimSpace(string(data))
	} else {
		if fs.NArg() != 3 {
			log.Fatalf(`expected three arguments, got %v`, fs.NArg())
		}

		in = fs.Arg(2)
	}

	context, err := pattern.Match(in)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf(`Context: %#v`, context)

	out, err := builder.Build(context)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(out)
}
