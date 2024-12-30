package libbysinglelib

import "github.com/invictus8957/library-search/internal/pkg/libby"

type libbySingleLib struct {
}

func (*libbySingleLib) SearchByAuthor(string) []libby.LibbyResult {
	return nil
}

func (*libbySingleLib) SearchByTitle(string) []libby.LibbyResult {
	return nil
}
