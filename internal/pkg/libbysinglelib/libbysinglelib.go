package libbysinglelib

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/invictus8957/library-search/internal/pkg/libby"
)

const defaultHTTPTimeout = 5 * time.Second
const maxPageSize = 100 // mucked around to find this
const pageNumberParamName = "page"
const pageSizeParamName = "perPage"

// format=ebook-overdrive,ebook-media-do,ebook-overdrive-provisional,audiobook-overdrive,audiobook-overdrive-provisional,magazine-overdrive
var constantSearchQueryParams = url.Values{
	"format": []string{
		"ebook-overdrive",
		"ebook-media-do",
		"ebook-overdrive-provisional",
		"audiobook-overdrive",
		"audiobook-overdrive-provisional",
		"magazine-overdrive",
	},
}

type libbySearchResponse struct {
	Items []libbySearchResponseItem  `json:"items"`
	Links []libbySearchResponseLinks `json:"links"`
}

type libbySearchResponseItem struct {
	// TODO -- pick fields
}

/*
	"links": {
	    "self": {
	      "page": 1,
	      "pageText": "1"
	    },
	    "first": {
	      "page": 1,
	      "pageText": "1"
	    },
	    "last": {
	      "page": 1,
	      "pageText": "1"
	    }
	  }
*/
type libbySearchResponseLinks struct {
	Last libbyPageInfo `json:"last"`
}

type libbyPageInfo struct {
	Page number `json:"page"`
}

type libbySingleLib struct {
	LibrarySearchURL string
	httpClient       *http.Client
}

func NewLibbySingleLib(libSearchURL string) *libbySingleLib {

	lsb := &libbySingleLib{
		LibrarySearchURL: libSearchURL,
	}
	lsb.httpClient = &http.Client{
		Timeout: defaultHTTPTimeout,
	}
	return lsb
}

func (*libbySingleLib) SearchByAuthor(string) ([]libby.LibbyResult, error) {
	return nil, errors.ErrUnsupported
}

func (*libbySingleLib) SearchByTitle(string) ([]libby.LibbyResult, error) {
	return nil, errors.ErrUnsupported
}

func (lsb *libbySingleLib) singlePageSearchRequest(q string, pageNum int, pageSize int) ([]libby.LibbyResult, error) {
	if pageSize > maxPageSize {
		return nil, fmt.Errorf("request page size of %v exceeded max page size of %v", pageSize, maxPageSize)
	}
	r, err := http.NewRequest("GET", lsb.LibrarySearchURL, nil)
	if err != nil {
		log.Printf("error creating single page search request, %s\n", err)
		return nil, err
	}
	queryVals := r.URL.Query()
	// add all the constant vals
	for k, _ := range constantSearchQueryParams {
		for v := range constantSearchQueryParams {
			queryVals.Add(k, v)
		}
	}
	queryVals.Add("query", q)
	queryVals.Add(pageNumberParamName, strconv.Itoa(pageNum))
	queryVals.Add(pageSizeParamName, strconv.Itoa(pageSize))

	resp, err := lsb.httpClient.Do(r)
	if err != nil {
		log.Printf("Error while fetching page of results: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	// FIXME: actually finish this
}
