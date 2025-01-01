package libbysinglelib

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"
)

// MockHTTPClient implements a mock HTTP client for testing
type MockHTTPClient struct {
	DoFunc  func(req *http.Request) (*http.Response, error)
	Timeout time.Duration
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

// TestNewLibbySingleLib tests the constructor
func TestNewLibbySingleLib(t *testing.T) {
	lib := NewLibbySingleLib("test-library-id")

	if lib.LibID != "test-library-id" {
		t.Errorf("Expected LibID to be 'test-library-id', got %s", lib.LibID)
	}

	if lib.librarySearchURL != "https://thunder.api.overdrive.com/v2/libraries/test-library-id/media" {
		t.Errorf("Incorrect search URL: %s", lib.librarySearchURL)
	}

	httpClient, _ := lib.httpClient.(*http.Client)
	if httpClient.Timeout != defaultHTTPTimeout {
		t.Errorf("Expected timeout to be %v, got %v", defaultHTTPTimeout, httpClient.Timeout)
	}
}

// mockSearchResponse creates a mock response for testing
func mockSearchResponse(items []LibbySearchResponseItem, nextPage *int) *http.Response {
	links := libbySearchResponseLinks{
		Self: &libbyPageInfo{Page: 1},
	}
	if nextPage != nil {
		links.Next = &libbyPageInfo{Page: *nextPage}
	}

	resp := libbySearchResponse{
		Items: items,
		Links: links,
	}

	jsonData, _ := json.Marshal(resp)
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBuffer(jsonData)),
	}
}

func TestSinglePageSearchRequest(t *testing.T) {
	mockClient := &MockHTTPClient{}
	lib := NewLibbySingleLib("test-lib")
	lib.httpClient = mockClient

	testCases := []struct {
		name          string
		pageSize      int
		mockResponse  *http.Response
		expectedError bool
		expectedItems int
		expectedNext  bool
	}{
		{
			name:     "Success with no next page",
			pageSize: 50,
			mockResponse: mockSearchResponse([]LibbySearchResponseItem{
				{Title: "Book 1", Author: "Author 1"},
				{Title: "Book 2", Author: "Author 2"},
			}, nil),
			expectedError: false,
			expectedItems: 2,
			expectedNext:  false,
		},
		{
			name:          "Exceeds max page size",
			pageSize:      maxPageSize + 1,
			mockResponse:  nil,
			expectedError: true,
			expectedItems: 0,
			expectedNext:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
				return tc.mockResponse, nil
			}

			items, links, err := lib.singlePageSearchRequest("test query", 1, tc.pageSize)

			if tc.expectedError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tc.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if len(items) != tc.expectedItems {
				t.Errorf("Expected %d items, got %d", tc.expectedItems, len(items))
			}

			if tc.expectedNext && links.Next == nil {
				t.Error("Expected next page link but got none")
			}

			if !tc.expectedNext && links != nil && links.Next != nil {
				t.Error("Expected no next page link but got one")
			}
		})
	}
}

func TestSearch(t *testing.T) {
	mockClient := &MockHTTPClient{}
	lib := NewLibbySingleLib("test-lib")
	lib.httpClient = mockClient

	nextPage := 2
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		// First page response
		if req.URL.Query().Get(pageNumberParamName) == "1" {
			return mockSearchResponse([]LibbySearchResponseItem{
				{Title: "Book 1"},
				{Title: "Book 2"},
			}, &nextPage), nil
		}
		// Second page response
		return mockSearchResponse([]LibbySearchResponseItem{
			{Title: "Book 3"},
		}, nil), nil
	}

	results, err := lib.Search("test query", 3)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	expectedTitles := []string{"Book 1", "Book 2", "Book 3"}
	for i, result := range results {
		if result.Title != expectedTitles[i] {
			t.Errorf("Expected title %s, got %s", expectedTitles[i], result.Title)
		}
	}
}

func TestSearchMaxResults(t *testing.T) {
	mockClient := &MockHTTPClient{}
	lib := NewLibbySingleLib("test-lib")
	lib.httpClient = mockClient

	nextPage := 2
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return mockSearchResponse([]LibbySearchResponseItem{
			{Title: "Book 1"},
			{Title: "Book 2"},
			{Title: "Book 3"},
		}, &nextPage), nil
	}

	results, err := lib.Search("test query", 2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results (max results), got %d", len(results))
	}
}
