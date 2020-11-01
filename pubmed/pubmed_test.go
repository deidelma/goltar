package pubmed

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

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
