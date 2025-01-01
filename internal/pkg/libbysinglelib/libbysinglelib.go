package libbysinglelib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
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

type LibbySearchResponseItem struct {
	Title             string                  `json:"title"`                // short title
	SortTitle         string                  `json:"sortTitle"`            // full title of multi-phrase
	Author            string                  `json:"firstCreatorName"`     // First Last
	SortAuthor        string                  `json:"firstCreatorSortName"` // last, first
	IsOwned           bool                    `json:"isOwned"`              // if they have a copy
	OwnedCopies       int                     `json:"ownedCopies"`
	AvailableCopies   int                     `json:"availableCopies"`
	HoldsCount        int                     `json:"holdsCount"`
	EstimatedWaitDays int                     `json:"estimatedWaitDays"`
	Type              libbySearchResponseType `json:"type"` // audiobook or ebook
}

type Lib struct {
	LibID            string
	librarySearchURL string
	httpClient       *http.Client
}

type libbySearchResponseLinks struct {
	Self *libbyPageInfo `json:"self"`
	Next *libbyPageInfo `json:"next"`
}

type libbyPageInfo struct {
	Page int `json:"page"`
}

type libbySearchResponse struct {
	Items []LibbySearchResponseItem `json:"items"`
	Links libbySearchResponseLinks  `json:"links"`
}

type libbySearchResponseType struct {
	ID string `json:"id"`
}

func NewLibbySingleLib(libraryID string) *Lib {

	lsb := &Lib{
		LibID:            libraryID,
		librarySearchURL: "https://thunder.api.overdrive.com/v2/libraries/" + libraryID + "/media",
	}
	lsb.httpClient = &http.Client{
		Timeout: defaultHTTPTimeout,
	}
	return lsb
}

func (lsb *Lib) Search(q string, maxResults int) ([]LibbySearchResponseItem, error) {
	return lsb.searchGetAllPagesUntilMax(q, maxResults)
}

func (lsb *Lib) searchGetAllPagesUntilMax(q string, maxResults int) ([]LibbySearchResponseItem, error) {
	results, links, err := lsb.singlePageSearchRequest(q, 1, maxPageSize)
	if err != nil {
		return nil, err
	}
	for {
		if len(results) < maxResults && links.Next != nil {
			intermediateResults, intermediateLinks, err := lsb.singlePageSearchRequest(q, links.Next.Page, maxPageSize)
			if err != nil {
				return nil, err
			}
			links = intermediateLinks
			results = append(results, intermediateResults...)

			continue
		}
		break
	}
	if len(results) == maxResults {
		return results, nil
	}
	if len(results) > maxResults {
		results = results[:maxResults]
		return results, nil
	}
	return results, nil
}

func (lsb *Lib) singlePageSearchRequest(q string, pageNum int, pageSize int) ([]LibbySearchResponseItem, *libbySearchResponseLinks, error) {
	if pageSize > maxPageSize {
		return nil, nil, fmt.Errorf("request page size of %v exceeded max page size of %v", pageSize, maxPageSize)
	}
	r, err := http.NewRequest("GET", lsb.librarySearchURL, nil)
	if err != nil {
		log.Printf("error creating single page search request, %s\n", err)
		return nil, nil, err
	}
	queryVals := r.URL.Query()
	// add all the constant vals
	for k := range constantSearchQueryParams {
		var sb strings.Builder // the api uses a comma-separated string as opposed to fmt=x&fmt=y
		for i, v := range constantSearchQueryParams[k] {
			sb.WriteString(v)
			if i != len(constantSearchQueryParams[k])-1 {
				sb.WriteString(",")
			}
		}
		queryVals.Add(k, sb.String())
	}
	queryVals.Add("query", q)
	queryVals.Add(pageNumberParamName, strconv.Itoa(pageNum))
	queryVals.Add(pageSizeParamName, strconv.Itoa(pageSize))

	// the api doesn't like encoded commas in the query string
	// so we have to hack around it
	qStr := queryVals.Encode()
	r.URL.RawQuery = strings.ReplaceAll(qStr, "%2C", ",")
	log.Printf("query string: %s\n", r.URL.RawQuery)

	resp, err := lsb.httpClient.Do(r)
	if err != nil {
		log.Printf("Error while fetching page of results: %v", err)
		return nil, nil, err
	}
	defer resp.Body.Close()
	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	//log.Printf("raw body: %v", string(rawBody))
	var data libbySearchResponse
	err = json.NewDecoder(bytes.NewReader(rawBody)).Decode(&data)
	if err != nil {
		return nil, nil, err
	}
	return data.Items, &data.Links, nil
}
