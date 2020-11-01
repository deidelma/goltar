package pubmed

// pubmed exports functionality that allows the parsing of an
// xml file into a Article struct and provides for it translation
// into a json or bson representation suitable for use with
// MongDB

// JournalInfo uniquely identifies a journal or equivalent
type JournalInfo struct {
	Title           string
	ISOAbbreviation string
}

// PubInfo provides information about a specific publication
type PubInfo struct {
	Date       string
	Vol        string
	Pages      string
	Issue      string
	EntrezDate string
}

// Author aims to uniquely identify an author
type Author struct {
	FirstName      string
	LastName       string
	Initials       string
	Affiliation    string
	CollectiveName string
}

// Article represents all of the retained
type Article struct {
	PMID     string
	Title    string
	Abstract string
	Journal  JournalInfo
	Issue    PubInfo
	Authors  []Author
	Keywords []string
}
