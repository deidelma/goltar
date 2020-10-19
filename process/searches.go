package jobs

// Search encapsulates data for a single Pubmed Search
//
// collection is the name of the MongoDB collection in use
// ands Terms to be joined by AND
// ors terms to be joined by OR
// authors list of authors in form ["hamid q", "o'byrne p"]
// startDate, endDate provide the years over which to search
type Search struct {
	collection string
	ands       []string
	ors        []string
	authors    []string
	startDate  string
	endDate    string
}

// Dates returns the startDate and endDate of the search
func (search *Search) Dates() (string, string) {
	return search.startDate, search.endDate
}

// OrTerms returns a slice of terms to be joined by OR
func (search *Search) OrTerms() []string {
	return search.ors
}

// Authors returns a slice of terms to be joined by AND
func (search *Search) Authors() []string {
	return search.ands
}

// Collection returns a slice of terms to be joined by AND
func (search *Search) Collection() string {
	return search.collection
}
