package pubmed

import (
	"fmt"
	"io"
	"strings"

	"github.com/antchfx/xmlquery"
	"github.com/deidelma/goltar/dbg"
)

// pubmed exports functionality that allows the parsing of an
// xml file into a Article struct and provides for it translation
// into a json or bson representation suitable for use with
// MongDB

// JournalInfo uniquely identifies a journal or equivalent
type JournalInfo struct {
	Title           string
	ISOAbbreviation string
}

// PubInfo provides information about a specific publication
type PubInfo struct {
	Date       string
	Vol        string
	Pages      string
	Issue      string
	EntrezDate string
}

// Author aims to uniquely identify an author
type Author struct {
	FirstName      string
	LastName       string
	Initials       string
	Affiliation    string
	CollectiveName string
}

// Article represents all of the retained information
// from a PubmedArticle entry in an xml record
//
type Article struct {
	PMID     string
	Title    string
	Abstract string
	Journal  JournalInfo
	Issue    PubInfo
	Authors  []Author
	Keywords []string
}

func parsePubInfo(article *Article, doc *xmlquery.Node) {
	volume, _ := getText(doc, "//JournalIssue/Volume")
	issue, _ := getText(doc, "//JournalIssue/Issue")
	date, _ := getText(doc, "//PubDate//MedlineDate")
	if date == "" {
		date, _ = getText(doc, "//PubDate//Year")
	}
	pagination, _ := getText(doc, "//Pagination/MedlinePgn")
	article.Issue.Date = date
	article.Issue.Pages = pagination
	article.Issue.Issue = issue
	article.Issue.Vol = volume
}

// parseAbstract loads abstract information even if it is
// stored in multiple parts
func parseAbstract(article *Article, nodes []*xmlquery.Node) {
	buffer := strings.Builder{}
	for _, node := range nodes {
		buffer.WriteString(node.InnerText())
		buffer.WriteString("\n")
	}
	article.Abstract = buffer.String()
}

// parseAuthors populates the article.Authors slice corresponding to the
// list of authors in the xml
func parseAuthors(article *Article, doc *xmlquery.Node) {
	nodes, _ := xmlquery.QueryAll(doc, "//AuthorList/Author")
	for _, node := range nodes {
		author := Author{}
		author.FirstName, _ = getText(node, "//ForeName")
		author.LastName, _ = getText(node, "//LastName")
		author.Initials, _ = getText(node, "//Initials")
		author.CollectiveName, _ = getText(node, "//CollectiveName")
		author.Affiliation, _ = getText(node, "//AffiliationInfo/Affiliation")
		article.Authors = append(article.Authors, author)
	}
}

// parseKeywords populates the article.Keywords slice corresponding
// to the accepted keywords in the xml record
func parseKeywords(article *Article, doc *xmlquery.Node) {
	kws := []string{}
	nodes, _ := xmlquery.QueryAll(doc, "//MeshHeading")
	for _, node := range nodes {
		item, _ := getText(node, "//DescriptorName")
		if len(item) > 0 {
			kws = append(kws, item)
		}
	}
	nodes, _ = xmlquery.QueryAll(doc, "//Keyword")
	for _, node := range nodes {
		kw := node.InnerText()
		kws = append(kws, kw)
	}
	article.Keywords = append(article.Keywords, kws...)
}

// textOne attempts to get a single element from the document
// path is of the form "//PMID"
// retrurns "" if no data is found for this path
// otherwise returns the InnerText of the node corresponding
// to the path
func getText(doc *xmlquery.Node, path string) (string, error) {
	node, err := xmlquery.Query(doc, path)
	if err != nil {
		return "", err
	}
	if node == nil {
		return "", err
	}
	return node.InnerText(), nil
}

// ParseXML reads the provided XML and attempts
// to populate an Article structure from it.
// Returns the Article
//
func ParseXML(xml io.Reader) (Article, error) {
	dbg.Printf("Entering ParseXML")
	result := Article{}
	doc, err := xmlquery.Parse(xml)
	if err != nil {
		return result, err
	}
	pmid, err := getText(doc, "//PMID")
	if err != nil {
		return result, err
	}
	// pmid is a required field
	if len(pmid) == 0 {
		return result, fmt.Errorf("Unable to find PMID")
	}
	result.PMID = pmid
	// assume that if there is a PMID, the xml record is valid so skip err check
	result.Title, _ = getText(doc, "//ArticleTitle")
	abstractNodes, _ := xmlquery.QueryAll(doc, "//Abstract/AbstractText")
	if abstractNodes == nil {
		result.Abstract, _ = getText(doc, "//Abstract")
	} else {
		parseAbstract(&result, abstractNodes)
	}
	parsePubInfo(&result, doc)
	result.Journal.Title, _ = getText(doc, "//Journal/Title")
	result.Journal.ISOAbbreviation, _ = getText(doc, "//ISOAbbreviation")
	parseAuthors(&result, doc)
	parseKeywords(&result, doc)

	dbg.Printf("Leaving ParseXML")
	return result, nil
}
