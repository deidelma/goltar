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
// years list of years to conduct the search [2001, 2002]
// if there are 3 entries and the last entry is a 1,
// treat the first two entries as an inclusive range [2001, 2010, 1]
// TODO: allow default collection in the header, handle ranges of years
type Search struct {
	Collection string
	Ands       []string
	Ors        []string
	Authors    []string
	Years      []int64
}

// Job encapsulates the data for a pubmed job.
// Currently only supports searches
//
// database is the name of the MongoDB database in use
//
type Job struct {
	Searches []Search
	Database string
	Name     string
}

// processSearch extracts the job information from a single
// element in a toml jobs file.
// Returns error if missing date or and terms information
func processSearch(t *toml.Tree) (Search, error) {
	result := Search{}
	ands := t.GetArray("ands")
	if ands == nil {
		return result, errors.New("missing AND terms")
	}
	result.Ands = ands.([]string)
	dateItems := t.GetArray("years").([]string)
	if dateItems == nil {
		return result, errors.New("missing years list")
	}

	return result, nil
}

// getSearches parses the job data looking for the key values
//
// ands -- an array of terms to be joined by AND [required]
// ors -- an array of terms to be joined by OR
// authors -- an array of author names joined by AND
// years -- an array of years [required]
func getSearches(searchTree []*toml.Tree) ([]Search, error) {
	result := []Search{}
	for _, aTree := range searchTree {
		s := Search{}
		ands := aTree.GetArray("ands")
		if ands == nil {
			return result, errors.New("missing AND clause in search")
		}
		s.Ands = ands.([]string)

		years := aTree.GetArray("years")
		if years == nil {
			return result, errors.New("missing years specification")
		}

		s.Years = years.([]int64)

		authors := aTree.GetArray("authors")
		if authors != nil {
			s.Authors = authors.([]string)
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
	value.Name = name.(string)
	value.Database = db.(string)
	value.Searches = searches
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
