package pubmed

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const apiKey = "5f3d0627bf00f873f63c871a085d82b25b08"
const eUtils = "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/"

var searchURL = fmt.Sprintf("%sesearch.fcgi?db=pubmed&api_key=%s&usehistory=y", eUtils, apiKey)
var fetchURL = fmt.Sprintf("%sefetch.fcgi?db=pubmed&api_key=%s", eUtils, apiKey)

// Query is a container for the results of an esearch
type Query struct {
	WebEnv   string
	QueryKey string
	Count    int64
}

func extractQuery(data []byte, query *Query) error {
	var f interface{}
	err := json.Unmarshal(data, &f)
	if err != nil {
		return err
	}
	m := f.(map[string]interface{})
	sresult := m["esearchresult"].(map[string]interface{})
	query.Count, _ = strconv.ParseInt(sresult["count"].(string), 10, 64)
	query.QueryKey = sresult["querykey"].(string)
	query.WebEnv = sresult["webenv"].(string)
	// return fmt.Errorf("Extract query not implemented")
	return nil
}

// ESearch executes a search on Pubmed using the provided terms
// terms is a string in the form "asthma+AND+copd+OR+cancer"
func ESearch(terms string) (Query, error) {
	result := Query{}
	url := fmt.Sprintf("%s&term=%s&retmode=json", searchURL, terms)
	log.Print(url)
	res, err := http.Get(url)
	if err != nil {
		return result, err
	}
	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	err = extractQuery(data, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}
