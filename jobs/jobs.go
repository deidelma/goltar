package jobs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

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
type Search struct {
	Collection string
	Ands       []string
	Ors        []string
	Authors    []string
	Years      []int64
}

// itemsToString concatenates strings using the provided conjunction
// suffix is then appended to each item before concatenation
// replaces spaces by '+'
// returns "" if items is empty or nil
func itemsToStringSuffix(items []string, conjunction string, suffix string) string {
	if items == nil || len(items) == 0 {
		return ""
	}
	result := []string{}
	first := fmt.Sprintf("%s%s", strings.ReplaceAll(items[0], " ", "+"), suffix)
	result = append(result, first)
	for _, item := range items[1:] {
		result = append(result, conjunction)
		result = append(result, fmt.Sprintf("%s%s", strings.ReplaceAll(item, " ", "+"), suffix))
	}
	return strings.Join(result, "+")

}

// itemsToString concatenates strings using the provided conjunction
// without a suffix (suffix == "")
// replaces space by '+'
// returns "" if items is empty or nil
func itemsToString(items []string, conjuction string) string {
	return itemsToStringSuffix(items, conjuction, "")
}

// URLEncode ensures that terms are properly URL encoded
func URLEncode(s string) string {
	result := strings.ReplaceAll(s, "'", "%27")
	result = strings.ReplaceAll(result, " ", "+")
	return result
}

// TermString returns strings suitable for use in a Pubmed
// search using the NLM entrez system.
//
// typical term string:
//	asthma+AND+leukotrienes+OR+interleukin-4+AND+o%27byrne+p[Au]+2003[pdat]
// 	returns one string for each year in Search.Years
func (search *Search) TermString() []string {
	result := []string{}
	for _, year := range search.Years {
		andSegment := itemsToString(search.Ands, "AND")
		orSegment := itemsToString(search.Ors, "OR")
		authorSegment := itemsToStringSuffix(search.Authors, "AND", "[Au]")
		items := []string{}
		items = append(items, andSegment)
		if orSegment != "" {
			items = append(items, fmt.Sprintf("+OR+%s", orSegment))
		}
		if authorSegment != "" {
			items = append(items, fmt.Sprintf("+AND+%s", authorSegment))
		}
		items = append(items, fmt.Sprintf("+AND+%d[pdat]", year))
		s := URLEncode(strings.Join(items, ""))
		// result = append(result, strings.Join(items, ""))
		result = append(result, s)
	}
	return result
}

// Job encapsulates the data for a pubmed job.
// Currently only supports searches
//
// database is the name of the MongoDB database in use
//
type Job struct {
	Searches   []Search
	Database   string
	Name       string
	Collection string // optional collection to use with all searches
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

// parseYears returns a slice of ints
// input can be with or without a hyphen
// "2001" or "2000-2020"
func parseYears(value string) ([]int64, error) {
	var result = []int64{}
	s := strings.Trim(value, " ")
	if strings.Contains(s, "-") {
		items := strings.Split(value, "-")
		yr1, err1 := strconv.ParseInt(items[0], 10, 64)
		if err1 != nil {
			return result, fmt.Errorf("Unable to convert string to int <%s>", items[0])
		}
		yr2, err2 := strconv.ParseInt(items[1], 10, 64)
		if err2 != nil {
			return result, fmt.Errorf("Unable to convert string to int <%s>", items[1])
		}
		for y := yr1; y <= yr2; y++ {
			result = append(result, y)
		}
		return result, nil
	}
	yr, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return result, fmt.Errorf("Unable to convert to int <%s>", value)
	}
	result = append(result, yr)
	return result, nil
}

// getYears returns a slice of values corresponding to the years
// over which this search will take place
// inputs come in two forms:
// 	[2001, 2002, 2003] an toml array of years as ints
// 	"2001-2005" a string representing a range of ints
// 	"2001" if no hyphen is present, the string has a single int
func getYears(t *toml.Tree) ([]int64, error) {
	result := []int64{}
	item := t.Get("years")
	if item == nil {
		return nil, fmt.Errorf("No years value given")
	}
	switch item.(type) {
	case string:
		years, err := parseYears(item.(string))
		if err != nil {
			return result, err
		}
		result = years
	default:
		result = t.GetArray("years").([]int64)
	}
	return result, nil
}

// getSearches parses the job data looking for the key values
//
// ands -- an array of terms to be joined by AND [required]
// ors -- an array of terms to be joined by OR
// authors -- an array of author names joined by AND
// years -- an array of years [required]
func getSearches(searchTree []*toml.Tree, masterCollection string) ([]Search, error) {
	result := []Search{}
	for _, aTree := range searchTree {
		s := Search{}
		ands := aTree.GetArray("ands")
		if ands == nil {
			return result, errors.New("missing AND clause in search")
		}
		s.Ands = ands.([]string)

		// years := aTree.GetArray("years")
		years, err := getYears(aTree)
		if err != nil {
			return result, err
		}
		if years == nil {
			return result, errors.New("missing years specification")
		}
		s.Years = years

		authors := aTree.GetArray("authors")
		if authors != nil {
			s.Authors = authors.([]string)
		}

		collection := aTree.Get("collection")
		if collection == nil {
			s.Collection = masterCollection
		} else {
			s.Collection = collection.(string)
		}
		result = append(result, s)
	}
	return result, nil
}

// tomlToJob takes an array of bytes in UTF-8 format and converts it to a job
func tomlToJob(data []byte) (Job, error) {
	value := Job{}
	t, err := toml.LoadBytes(data)
	if err != nil {
		return Job{}, err
	}
	db := t.Get("goltar.database")
	if db == nil {
		return value, errors.New("missing database name in jobs file")
	}
	name := t.Get("goltar.name")
	if name == nil {
		return value, errors.New("missing job name in jobs file")
	}

	// if no collection given, then use the search name
	coll := t.Get("goltar.collection")
	var collection string
	if coll != nil {
		collection = coll.(string)
	} else {
		collection = name.(string)
	}

	ss := t.GetArray("searches")
	if ss == nil {
		return value, errors.New("no searches provided in jobs file")
	}
	searches, err := getSearches(ss.([]*toml.Tree), collection)
	if err != nil {
		return value, err
	}
	value.Name = name.(string)
	value.Database = db.(string)
	value.Collection = collection
	value.Searches = searches
	return value, nil
}

// ReadJobString is convenience interface to the parser.
// It takes data as string and executes tomlToJob in order
// to convert it to a job.
func ReadJobString(data string) (Job, error) {
	return tomlToJob([]byte(data))
}

// ReadJobFile reads a toml file containing job information
// and then attempts to convert it to a Job struct
// path is a properly qualified path to the file
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
