package pubmed

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/antchfx/xmlquery"
	xj "github.com/basgys/goxml2json"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func TestConvertFileToJson(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/asthma_article_one.xml")
	check(err)
	xml := strings.NewReader(string(data))
	json, err := xj.Convert(xml)
	check(err)
	fmt.Println(json.String())

}
func TestConvertFileToJson2(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/asthma_article_two.xml")
	check(err)
	xml := strings.NewReader(string(data))
	json, err := xj.Convert(xml)
	check(err)
	fmt.Println(json.String())
}

func TestFindPMID(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/asthma_article_one.xml")
	check(err)
	xml := strings.NewReader(string(data))
	doc, err := xmlquery.Parse(xml)
	check(err)
	pmid, err := xmlquery.Query(doc, "//PMID")
	check(err)
	fmt.Printf("PMID:%s\n", pmid.InnerText())
	atitle, err := xmlquery.Query(doc, "//ArticleTitle")
	check(err)
	if atitle != nil {
		fmt.Printf("ArticleTitle:%s\n", atitle.InnerText())
	}
	aText, err := xmlquery.QueryAll(doc, "//AbstractText")
	check(err)
	if aText == nil {
		return
	}
	for _, t := range aText {
		fmt.Printf("abstract text element: %s\n", t.InnerText())
	}
}
