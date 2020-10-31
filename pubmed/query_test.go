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
authors=["hamid q"]
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
	if q.Count != 17 {
		t.Errorf("Expected 17, received %d", q.Count)
	}

}
