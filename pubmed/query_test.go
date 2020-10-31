package pubmed

import (
	"log"
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

const bigJobToml = `
[goltar]
database="goltar"
name="test"
collection="asthma"
[[searches]]
ands=["asthma"]
years=[2010]
`

func TestFetchJob(t *testing.T) {
	job, _ := jobs.ReadJobString(bigJobToml)
	search := job.Searches[0]
	terms := search.TermString()[0]
	q, err := ESearch(terms)
	if err != nil {
		t.Errorf("Error during search:[%v]", err)
	}
	log.Printf("Downloading %d records", q.Count)
	recs, err := EFetchRecs(q)
	if err != nil {
		t.Errorf("Error fetching data:[%v]", err)
	}
	if len(recs) != int(q.Count-5) {
		t.Errorf("Expected %d, received %d", q.Count-5, len(recs))
	}
	log.Printf("Received %d records", len(recs))
}

func TestGenerateSlices(t *testing.T) {
	max := 4227
	size := 500
	slices := generateSlices(max, size)

	for i, n := range slices {
		log.Printf("%d: %d => %d", i, n, n+size)
	}
}
