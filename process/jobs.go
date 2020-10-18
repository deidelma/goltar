package jobs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/pelletier/go-toml"
)

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

// Job encapsulates the data for a pubmed job.
// Currently only supports searches
//
// database is the name of the MongoDB database in use
//
type Job struct {
	searches []Search
	database string
	name     string
}

// Name returns the name associated with the job.
func (job *Job) Name() string {
	return job.name
}

// Searches returns the searches associated with the Job.
func (job *Job) Searches() []Search {
	return job.searches
}

func tomlToJob(data []byte) (Job, error) {
	value := Job{}
	t, err := toml.LoadBytes(data)
	if err != nil {
		return Job{}, err
	}
	for _, k := range t.Keys() {
		fmt.Printf("Key: %s\n", k)
	}
	db := t.Get("goltar.database")
	if db == nil {
		return value, errors.New("Missing database name in jobs file")
	}
	name := t.Get("goltar.name")
	if name == nil {
		return value, errors.New("Missing job name in jobs file")
	}
	ss := t.GetArray("searches")
	if ss == nil {
		log.Println("No searches provided in jobs file.")
		return value, nil
	}
	searches := []Search{}
	sarray := ss.([]*toml.Tree)
	for _, tt := range sarray {
		ands := tt.GetArray("ands")
		if ands == nil {
			log.Println("Missing and clause in search")
			return value, nil
		}
		startDate := tt.Get("startDate")
		if startDate == nil {
			log.Println("Missing start date in search")
			return value, nil
		}
		s := Search{}
		s.ands = ands.([]string)
		s.startDate = startDate.(string)
		ors := tt.GetArray("ors")
		if ors != nil {
			s.ors = ors.([]string)
		}
		authors := tt.GetArray("authors")
		if authors != nil {
			s.authors = authors.([]string)
		}
		endDate := tt.Get("endDate")
		if endDate != nil {
			s.endDate = endDate.(string)
		}
		searches = append(searches, s)

		value.name = name.(string)
		value.database = db.(string)
		value.searches = searches
	}

	// searches := []Search{}

	// for _, stree :=  {

	// }
	value.name = name.(string)
	value.database = db.(string)
	// value.searches = searches.([]Search)
	return value, nil
}

// ReadJobFile reads a toml file containing job information
// and then attempts to convert it to a Job struct
func ReadJobFile(path string) (Job, error) {
	log.Printf("Loading file %s\n", path)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return Job{}, err
	}

	// convert the toml data to job format
	job, err := tomlToJob(data)
	if err != nil {
		return Job{}, err
	}
	return job, nil

}
