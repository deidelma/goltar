package jobs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/pelletier/go-toml"
)

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

func getSearches(searchTree []*toml.Tree) ([]Search, error) {
	result := []Search{}
	for _, aTree := range searchTree {
		ands := aTree.GetArray("ands")
		if ands == nil {
			return result, errors.New("missing AND clause in search")
		}
		startDate := aTree.Get("startDate")
		if startDate == nil {
			return result, errors.New("missing start date in search")
		}
		s := Search{}
		s.ands = ands.([]string)
		s.startDate = startDate.(string)
		ors := aTree.GetArray("ors")
		if ors != nil {
			s.ors = ors.([]string)
		}
		authors := aTree.GetArray("authors")
		if authors != nil {
			s.authors = authors.([]string)
		}
		endDate := aTree.Get("endDate")
		if endDate != nil {
			s.endDate = endDate.(string)
		}
		result = append(result, s)
	}
	return result, nil
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
		return value, errors.New("missing database name in jobs file")
	}
	name := t.Get("goltar.name")
	if name == nil {
		return value, errors.New("missing job name in jobs file")
	}
	ss := t.GetArray("searches")
	if ss == nil {
		return value, errors.New("no searches provided in jobs file")
	}
	searches, err := getSearches(ss.([]*toml.Tree))
	if err != nil {
		return value, err
	}
	value.name = name.(string)
	value.database = db.(string)
	value.searches = searches
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
