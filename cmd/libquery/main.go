package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/invictus8957/library-search/internal/pkg/libby"
)

type Config struct {
	maxResults int
	q          string
}

func parseFlags() *Config {
	c := &Config{}

	flag.IntVar(&c.maxResults, "max", 10, "Max results to return.")
	flag.StringVar(&c.q, "query", "", "Query string to search.")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "Example:\n")
		fmt.Fprintf(os.Stderr, "%s -max 10 -query foo", os.Args[0])
	}

	flag.Parse()

	if c.q == "" {
		fmt.Fprintf(os.Stderr, "Query string cannot be empty.\n")
		flag.Usage()
	}
	return c
}

func main() {
	c := parseFlags()

	const libID = "lexpublib"
	l := libby.NewLibby([]string{libID})
	results, err := l.Search(c.q, c.maxResults)
	if err != nil {
		log.Fatalf("Error searching for results in libby: %v", err)
	}
	log.Printf("Results: %v", results)
}
