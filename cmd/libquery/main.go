package main

import (
	"log"

	"github.com/invictus8957/library-search/internal/pkg/libby"
)

func main() {
	const libID = "lexpublib"
	l := libby.NewLibby([]string{libID})
	results, err := l.Search("robert henderson", 10)
	if err != nil {
		log.Fatalf("Error searching for results in libby: %v", err)
	}
	log.Printf("Results: %v", results)
}
