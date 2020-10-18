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
	sarray := ss.(*[]go-toml.Tree)

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
