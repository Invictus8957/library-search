package libby

import (
	"errors"
	"log"

	"github.com/invictus8957/library-search/internal/pkg/libbysinglelib"
)

// we want to make this a simple interface with a complex implementation that hides the aggregation
// necessary to use libby across multiple libraries since they don't seem to have a multi-library query api
// (or would they if I had multiple libraries setup?).

// for a single library we want to do all the query operations as well, and then aggregate
// in the main impl across the multiple libraries

type LibbyLibrary struct {
	ID        string // the string id used in searches
	WebsiteID int    // this is the website id needed to lookup the id string
}

type LibbyResult struct {
	Author            string
	Title             string
	LongTitle         string
	Library           string
	TotalCopies       int
	IsAvailable       bool
	AvailableCopies   int
	HoldsCount        int
	EstimatedWaitDays int
	Type              string // ebook or audiobook
}

type Libby interface {
	FindLibrary(string) ([]LibbyLibrary, error)
	Search(query string, maxResults int) ([]LibbyResult, error)
}

func NewLibby(libraryIDs []string) Libby {
	var libs []libbysinglelib.Lib
	for _, v := range libraryIDs {
		libs = append(libs, *libbysinglelib.NewLibbySingleLib(v))
	}
	return &LibbyImpl{
		libraries: libs,
	}
}

type LibbyImpl struct {
	libraries []libbysinglelib.Lib
}

// TODO: implement find library
func (li *LibbyImpl) FindLibrary(string) ([]LibbyLibrary, error) {
	return nil, errors.ErrUnsupported
}

func (li *LibbyImpl) Search(query string, maxResults int) ([]LibbyResult, error) {
	log.Println("Beginning top level search invocation.")
	var results []LibbyResult
	for _, lib := range li.libraries {
		res, err := lib.Search(query, maxResults)
		if err != nil {
			return nil, err
		}
		results = append(results, convertResults(res, lib.LibID)...)
	}
	return results, nil
}

func convertResults(r []libbysinglelib.LibbySearchResponseItem, libraryID string) []LibbyResult {
	var libbyResults []LibbyResult
	for _, item := range r {
		libbyResults = append(libbyResults, LibbyResult{
			Author:            item.Author,
			Title:             item.Title,
			Library:           libraryID,
			TotalCopies:       item.OwnedCopies,
			IsAvailable:       (item.OwnedCopies > 0 && item.AvailableCopies > 0),
			AvailableCopies:   item.AvailableCopies,
			HoldsCount:        item.HoldsCount,
			EstimatedWaitDays: item.EstimatedWaitDays,
			Type:              item.Type.ID,
		})
	}
	return libbyResults
}
