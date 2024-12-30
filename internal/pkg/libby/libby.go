package libby

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
}

type Libby interface {
	FindLibrary(string) []LibbyLibrary
	SearchByAuthor(string) []LibbyResult
	SearchByTitle(string) []LibbyResult
}

// probably want this to not be a const because testing
var libbyURL = "https://thunder.api.overdrive.com"
