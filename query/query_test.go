package pubmed

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/deidelma/goltar/dbg"
	jobs "github.com/deidelma/goltar/jobs"
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
func TestFetchJobSync(t *testing.T) {
	job, _ := jobs.ReadJobString(bigJobToml)
	search := job.Searches[0]
	terms := search.TermString()[0]
	q, err := ESearch(terms)
	if err != nil {
		t.Errorf("Error during search:[%v]", err)
	}
	log.Printf("Downloading %d records", q.Count)
	data, err := EFetchSync(q)
	if err != nil {
		t.Errorf("Error fetching data:[%v]", err)
	}
	recs := SplitXML(string(data), "PubmedArticle")
	if len(recs) != int(q.Count-5) {
		t.Errorf("Expected %d, received %d", q.Count-5, len(recs))
	}
	log.Printf("Received %d records", len(recs))
}
func TestGenerateSlices(t *testing.T) {
	max := 4227
	size := 500
	slices := generateSlices(max, size)

	dbg.Start()
	for i, n := range slices {
		dbg.Printf("%d: %d => %d", i, n, n+size)
	}
}

func TestParseXML(t *testing.T) {
	file := "testdata/asthma200.xml"
	data, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatalf("Unable to find data:%s", file)
	}
	recs := SplitXML(string(data), "PubmedArticle")
	if len(recs) != 200 {
		t.Errorf("Expected 200, receieved %d", len(recs))
	}
}

func TestSearchReturnsCorrectCount(t *testing.T) {
	terms := jobs.URLEncode("asthma AND leukotrienes AND o'byrne p[au]")
	dbg.Printf("Terms:%s", terms)
	q, err := ESearch(terms)
	if err != nil {
		t.Errorf("Failed to complete search: <%v>", err)
	}
	if q.Count != 37 {
		t.Errorf("Expected 37, received %d", q.Count)
	}
}

func TestFetchReturnsCorrectNumberOfRecords(t *testing.T) {
	terms := jobs.URLEncode("asthma AND leukotrienes AND o'byrne p[au]")
	q, _ := ESearch(terms)
	recs, err := EFetchRecs(q)
	if err != nil {
		t.Errorf("Failed to fetch records: <%v>", err)
	}
	if len(recs) != 37 {
		t.Errorf("Expected 37 recs, received %d", len(recs))
	}

}
