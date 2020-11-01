package pubmed

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const apiKey = "5f3d0627bf00f873f63c871a085d82b25b08"
const eUtils = "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/"
const sliceSize = 1000

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
	//log.Print(url)
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

// myClient  returnws a modified http client with a larger connection pool
// based on:
//http://tleyden.github.io/blog/2016/11/21/tuning-the-go-http-client-library-for-load-testing/
func myClient() *http.Client {
	defaultRoundTripper := http.DefaultTransport
	defaultTransportPointer, ok := defaultRoundTripper.(*http.Transport)
	if !ok {
		panic(fmt.Sprintf("defaultRoundTripper not an *http.Transport"))
	}
	defaultTransport := *defaultTransportPointer // deref it to get a copy of the struct
	defaultTransport.MaxIdleConns = 100
	defaultTransport.MaxIdleConnsPerHost = 100
	return &http.Client{Transport: &defaultTransport}
}

// fetchSlice returns the result of executing efetch.cgi to receive the records
// from start to start + nrecs
func fetchSlice(q Query, start int, nrecs int) ([]byte, error) {
	run := 1
	retries := 200
	notDone := true
	result := []byte{}
	url := fmt.Sprintf("%s&WebEnv=%s&query_key=%s&retstart=%d&retmax=%d&retmode=xml",
		fetchURL, q.WebEnv, q.QueryKey, start, nrecs)
	// log.Printf("URL:%s", url)
	for notDone {
		// res, err := http.Get(url)
		res, err := myClient().Get(url)
		if err != nil {
			return result, err
		}
		if res.StatusCode == 200 {
			defer res.Body.Close()
			result, err = ioutil.ReadAll(res.Body)
			break
		} else {
			run++
			// log.Printf("Trying again <%d>", run)
			time.Sleep(1)
		}
		if run >= retries {
			return result, fmt.Errorf("too many retries")
		}
	}
	return result, nil
}

// geneateSlices returns a slice of int containing
// the starting points for search slices
func generateSlices(max int, sliceSize int) []int {
	result := []int{}
	for i := 0; i < max; i += sliceSize {
		result = append(result, i)
	}
	return result
}

// EFetchSync returns the results of call efetch.cgi synchronously
// This is intended to be used for testing purposes only
// and not for production
func EFetchSync(q Query) (string, error) {
	result := strings.Builder{}
	if q.Count == 0 {
		return result.String(), nil
	}
	sliceStarts := generateSlices(int(q.Count), sliceSize)

	// var wg sync.WaitGroup
	for i, slice := range sliceStarts {
		xml, err := fetchSlice(q, slice, sliceSize)
		if err != nil {
			log.Fatalf("Unable to download slice: %v", err)
		}
		result.WriteString(string(xml))
		recs := parseXML(string(xml), "PubmedArticle")
		log.Printf("%d) found %d recs after fetchSlice", i, len(recs))
		time.Sleep(1)
	}
	log.Printf("Downloaded %d bytes of xml", len(result.String()))
	return result.String(), nil
}

// EFetchRecs returns the results of call efetch.cgi
// returns a slice of XML records
func EFetchRecs(q Query) ([]string, error) {
	result := []string{}
	if q.Count == 0 {
		return result, nil
	}
	sliceStarts := generateSlices(int(q.Count), sliceSize)

	var wg sync.WaitGroup
	for i, slice := range sliceStarts {
		wg.Add(1)
		go func(q Query, slice int, sliceSize int, i int) {
			xml, err := fetchSlice(q, slice, sliceSize)
			if err != nil {
				log.Fatalf("Unable to download slice: %v", err)
			}
			recs := parseXML(string(xml), "PubmedArticle")
			log.Printf("%d) found %d recs after fetchSlice", i, len(recs))
			result = append(result, recs...)
			wg.Done()
		}(q, slice, sliceSize, i)
	}
	wg.Wait()
	return result, nil
}

// EFetch returns the results of call efetch.cgi
func EFetch(q Query) (string, error) {
	result := strings.Builder{}
	if q.Count == 0 {
		return result.String(), nil
	}
	sliceStarts := generateSlices(int(q.Count), sliceSize)

	var wg sync.WaitGroup
	for i, slice := range sliceStarts {
		wg.Add(1)
		go func(q Query, slice int, sliceSize int, i int) {
			// log.Printf("%d) Downloading slice %d size %d", i, slice, sliceSize)
			xml, err := fetchSlice(q, slice, sliceSize)
			if err != nil {
				log.Fatalf("Unable to download slice: %v", err)
			}
			// log.Printf("%d) Adding data of %d bytes to buffer", i, len(xml))
			result.WriteString(string(xml))
			recs := parseXML(string(xml), "PubmedArticle")
			log.Printf("%d) found %d recs after fetchSlice", i, len(recs))
			// if len(recs) < sliceSize {
			// 	log.Fatalln(string(xml))
			// }
			wg.Done()
		}(q, slice, sliceSize, i)
	}
	wg.Wait()
	log.Printf("Downloaded %d bytes of xml", len(result.String()))
	return result.String(), nil
}

// parseXML splits the xml string into blocks demarked by
// the provided tag
func parseXML(xml string, tag string) []string {
	result := []string{}
	start := fmt.Sprintf("<%s>", tag)
	end := fmt.Sprintf("</%s>", tag)
	rawLines := strings.Split(xml, "\n")
	inArticle := false
	buffer := []string{}
	for _, line := range rawLines {
		// line := strings.Trim(rawLine, " ")
		if strings.Contains(line, start) {
			inArticle = true
			buffer = append(buffer, line)
		} else if strings.Contains(line, end) {
			inArticle = false
			buffer = append(buffer, line)
			articleStr := strings.Join(buffer, "\n")
			result = append(result, articleStr)
			buffer = []string{}
		} else if inArticle {
			buffer = append(buffer, line)
		}
	}
	return result
}

// FetchRecords performs the search indicated by the terms.
//
// terms is in proper url format
// returns a slice containing XML documents
func FetchRecords(terms string) ([]string, error) {
	result := []string{}
	q, err := ESearch(terms)
	if err != nil {
		return result, err
	}
	// xml, err := EFetch(q)
	xml, err := EFetchSync(q)
	if err != nil {
		return result, err
	}
	result = parseXML(xml, "PubmedArticle")
	return result, nil
}
