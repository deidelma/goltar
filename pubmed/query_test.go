package pubmed

import (
	"testing"

	jobs "github.com/deidelma/goltar/process"
)

const jobToml = `
[goltar]
database="goltar"
name="test"
collection="asthma"
[[searches]]
ands=["asthma"]
authors=["o'byrne p"]
years=[2000]
`

func TestConnectToServer(t *testing.T) {
	job, _ := jobs.ReadJobString(jobToml)
	search := job.Searches[0]
	terms := search.TermString()[0]
	q, err := ESearch(terms)
	if err != nil {
		t.Errorf("Error connecting to server:[%v]", err)
	}
	if q.Count != 14 {
		t.Errorf("Expected 14, received %d", q.Count)
	}

}
