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

func TestParseXML(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/asthma_article_one.xml")
	check(err)
	xml := strings.NewReader(string(data))
	article, err := ParseXML(xml)
	check(err)
	if article.PMID != "32374540" {
		t.Errorf("Incorrect PMID: %s", article.PMID)
	}
	if !strings.HasPrefix(article.Title, "[Asthma and COPD") {
		t.Errorf("Incorrect title: %s", article.Title)
	}
	if !strings.Contains(article.Abstract, "Numerous patients with asthma") {
		t.Errorf("Incorrect abstract: %s ", article.Abstract)
	}
	pinfo := article.Issue
	if pinfo.Vol != "16" {
		t.Errorf("Incorrect volume number: %s", pinfo.Vol)
	}
	if article.Journal.Title != "Revue medicale suisse" {
		t.Errorf("Wrong journal title: %s", article.Journal.Title)
	}
	if article.Journal.ISOAbbreviation != "Rev Med Suisse" {
		t.Errorf("Wrong journal iso: %s", article.Journal.ISOAbbreviation)
	}
	if article.Authors[0].LastName != "Daccord" {
		t.Errorf("Wrong last name: %s", article.Authors[0].LastName)
	}
	if article.Authors[0].FirstName != "Cécile" {
		t.Errorf("Wrong first name: %s", article.Authors[0].FirstName)
	}
	if article.Keywords[3] != "Humans" {
		t.Errorf("Wrong keyword: %s", article.Keywords[3])
	}
}

func TestParseXMLTwo(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/asthma_article_two.xml")
	check(err)
	xml := strings.NewReader(string(data))
	article, err := ParseXML(xml)
	check(err)
	if article.PMID != "31747880" {
		t.Errorf("Incorrect PMID: %s", article.PMID)
	}
	if !strings.HasPrefix(article.Title, "The role of secreted Hsp90") {
		t.Errorf("Incorrect title: %s", article.Title)
	}
	if !strings.Contains(article.Abstract, "The dysfunction of airway") {
		t.Errorf("Incorrect abstract: %s ", article.Abstract)
	}
	pinfo := article.Issue
	if pinfo.Vol != "19" {
		t.Errorf("Incorrect volume number: %s", pinfo.Vol)
	}
	if article.Journal.Title != "BMC pulmonary medicine" {
		t.Errorf("Wrong journal title: %s", article.Journal.Title)
	}
	if article.Journal.ISOAbbreviation != "BMC Pulm Med" {
		t.Errorf("Wrong journal iso: %s", article.Journal.ISOAbbreviation)
	}
	if article.Authors[0].LastName != "Ye" {
		t.Errorf("Wrong last name: %s", article.Authors[0].LastName)
	}
	if article.Authors[0].FirstName != "Cuiping" {
		t.Errorf("Wrong first name: %s", article.Authors[0].FirstName)
	}
	if article.Keywords[3] != "Cadherins" {
		t.Errorf("Wrong keyword: %s", article.Keywords[3])
	}
	i := len(article.Keywords) - 1
	if article.Keywords[i] != "Secreted Hsp90α" {
		t.Errorf("Wrong keyword: %s", article.Keywords[i])
	}
	if article.Keywords[i-1] != "HDM" {
		t.Errorf("Wrong keyword: %s", article.Keywords[i-1])
	}
}
