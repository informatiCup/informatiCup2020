package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/BitFlipp/informatiCup2020/game"
	"github.com/spf13/pflag"
)

const (
	version = "2.2.0"
)

func main() {
	flags := &pflag.FlagSet{}
	f := flags.StringP("log-file-path", "o", "", `Log file path. If not set or empty, the standard output will be used`)
	s := flags.Int64P("random-seed", "s", 0, "Random seed. If not set or 0, the Unix timestamp (UTC) in nanoseconds will be used")
	t := flags.IntP("request-timeout", "t", 10*1000, "Request timeout in milliseconds, >= 0. If 0, the HTTP client will wait indefinitely")
	u := flags.StringP("endpoint-url", "u", "http://localhost:50123", "Endpoint URL, must not be empty")
	b := new(bytes.Buffer)
	flags.SetOutput(b)

	usage := func() {
		fmt.Printf("informatiCup 2020 command line tool v%s\n\nFlags:\n", version)
		fmt.Println(flags.FlagUsages())
		os.Exit(2)
	}

	if err := flags.Parse(os.Args[1:]); err != nil {
		usage()
	}

	if *t < 0 {
		usage()
	}

	if *u == "" {
		usage()
	}

	sd := time.Now().UTC().UnixNano()
	if *s != 0 {
		sd = *s
	}
	rand.Seed(sd)

	lf := os.Stdout
	if *f != "" {
		if of, err := os.Create(*f); err != nil {
			log.Fatalf("failed to create log file: %s", err)
		} else {
			lf = of
		}
	}

	g := game.New()
	if err := g.Run(*u, *t, lf); err != nil {
		log.Fatalf("failed to run game: %s", err)
	}
}
