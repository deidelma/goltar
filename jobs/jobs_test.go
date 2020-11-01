package jobs

import "testing"

var data1 = `
[goltar]
name="david"
database="asthma"

[[searches]]
ands=["asthma"]
years=[2001]
`

func TestSimpleSearch(t *testing.T) {
	j, err := tomlToJob([]byte(data1))
	if err != nil {
		t.Errorf("Failed to read job")
	}
	if v := j.Name; v != "david" {
		t.Errorf("Expected 'david', received '%v'", v)
	}
	if n := len(j.Searches); n != 1 {
		t.Errorf("Expected only one search in this job")
	}
	if ands := j.Searches[0].Ands; len(ands) != 1 {
		t.Errorf("Expected only one and term")
	}
	if term := j.Searches[0].Ands[0]; term != "asthma" {
		t.Errorf("Expected 'asthma', received '%s'", term)
	}

}

var data2 = `
[goltar]
name="david"
database="goltar"

[[searches]]
ands=["asthma", "copd"]
authors=["martin jg","hamid q"]
years=[2000, 2001]

[[searches]]
ands=["asthma", "copd"]
ors = ["leukotrienes"]
authors=["martin jg","minshall e"]
years=[2000, 2001, 1]

[[searches]]
ands=["pulmonary fibrosis", "steroids"]
ors = ["methotrexate", "azothioprine"]
authors=["colman n", "zackon h", "cosio m"]
years=[2005, 2010, 2011]
`

func TestMultipleSearch(t *testing.T) {
	j, err := tomlToJob([]byte(data2))
	if err != nil {
		t.Errorf("Unable to process data2")
	}
	t.Logf("Loaded job named %s", j.Name)
	searches := j.Searches
	if len(searches) != 3 {
		t.Errorf("Failed to find 3 searches in %s", j.Name)
	}
	s1 := searches[0]
	if ors := s1.Ors; ors != nil {
		t.Error("Erroneous or!")
	}
	if len(s1.Authors) != 2 {
		t.Errorf("Wrong number of authors")
	}
	s2 := searches[1]
	if len(s2.Years) != 3 {
		t.Error("Wrong number of years")
	}
	if s2.Years[2] != 1 {
		t.Error("Missing sequence marker")
	}

	s3 := searches[2]
	if len(s3.Authors) != 3 {
		t.Error("Wrong number of authors in 3rd search")
	}
	if s3.Years[2] != 2011 {
		t.Error("Wrong year in third search")
	}
}

var data3 = `
[goltar]
name="david"
database="goltar"

[[searches]]
ands = ["asthma", "copd"]
years = "2001-2010"
`

func TestYearRange(t *testing.T) {
	job, err := tomlToJob([]byte(data3))
	if err != nil {
		t.Error("Unable to parse data3")
	}
	s := job.Searches[0]
	if len(s.Years) != 10 {
		t.Errorf("Wrong number of years.  Expected 10, got %d", len(s.Years))
	}
}

func TestParseYears(t *testing.T) {
	yrs, err := parseYears("2001")
	if err != nil {
		t.Errorf("Parser failure: %s", "2001")
	}
	if len(yrs) != 1 {
		t.Errorf("Expected 1 year, got %d", len(yrs))
	}

	yrs, err = parseYears("2001-2003")
	if err != nil {
		t.Errorf("Parser failure: %s", "2001-2003")
	}
	if len(yrs) != 3 {
		t.Errorf("Expected 3 years, got %d", len(yrs))
	}

	yrs, err = parseYears("bob")
	if err == nil {
		t.Errorf("Failed to detect illegal year")
	}

	yrs, err = parseYears("2000-")
	if err == nil {
		t.Errorf("Failed to detect illegal value for year")
	}
}

func TestCollection(t *testing.T) {
	job, _ := tomlToJob([]byte(data3))
	if job.Collection != job.Name {
		t.Error("Failed to assign default value to master collection")
	}
	s := job.Searches[0]
	if s.Collection != job.Name {
		t.Error("Failed to assign default collection to search")
	}
}

var data4 = `
[goltar]
name="bob"
collection="asthma"
database = "goltar"
[[searches]]
ands = ["asthma"]
years=[2001]

[[searches]]
ands=["copd"]
years="2005-2007"
collection="copd"
`

func TestCollectionAssignment(t *testing.T) {
	job, _ := tomlToJob([]byte(data4))
	if job.Collection != "asthma" {
		t.Error("Failed to assign default value to master collection")
	}
	s := job.Searches[0]
	if s.Collection != job.Collection {
		t.Error("Failed to assign default collection to search")
	}
	s = job.Searches[1]
	if s.Collection != "copd" {
		t.Errorf("Expected 'copd', recevied '%s'", s.Collection)
	}
}

var badData1 = `
[zoltar]
name="david"
database="david"
[[searches]]
ands=["asthma"]
years[2001]
`

var badData2 = `
[goltar]
name="david"
database="david"
[[searches]]
ands=["asthma"]
# years[2001]
`

func TestBadInputs(t *testing.T) {
	_, err := tomlToJob([]byte(badData1))
	if err == nil {
		t.Error("Failed to detect invalid file: badData1")
	}
	_, err1 := tomlToJob([]byte(badData2))
	if err1 == nil {
		t.Error("Failed to detect invalid file: badData2")
	}
}

func TestItemsToString(t *testing.T) {
	sample := []string{"asthma", "copd", "sam smith"}
	s := itemsToString(sample, "AND")
	if s != "asthma+AND+copd+AND+sam+smith" {
		t.Errorf("Expected:'asthma+AND+copd+AND+sam+smith', received '%s'", s)
	}
	sample = []string{"o'byrne p", "martin jg"}
	s = itemsToString(sample, "OR")
	if s != "o'byrne+p+OR+martin+jg" {
		t.Errorf("Expected:'o'byrne+p+OR+martin+j', received '%s'", s)
	}
	sample = []string{"asthma"}
	s = itemsToString(sample, "AND")
	if s != "asthma" {
		t.Errorf("Expected 'asthma', received:'%s'", s)
	}
	sample = []string{}
	s = itemsToString(sample, "AND")
	if s != "" {
		t.Error("Mishandled empty array")
	}
	s = itemsToString(nil, "AND")
	if s != "" {
		t.Error("Mishandled null array")
	}
}

func TestSingleTermString(t *testing.T) {
	s := Search{
		Ands:       []string{"asthma"},
		Ors:        []string{"copd"},
		Authors:    nil,
		Years:      []int64{2001},
		Collection: "",
	}
	terms := s.TermString()
	if len(terms) != 1 {
		t.Error("Wrong number of term strings, expected 1")
	}
	if terms[0] != "asthma+OR+copd+AND+2001[pdat]" {
		t.Errorf("Term string failed:%s", terms)
	}
}

func TestMultipleTermString(t *testing.T) {
	s := Search{
		Ands:       []string{"asthma"},
		Ors:        []string{"copd", "pulmonary fibrosis"},
		Authors:    []string{"martin jg", "o'byrne p", "hamid q"},
		Years:      []int64{2001, 2002, 2003},
		Collection: "",
	}
	terms := s.TermString()
	if len(terms) != 3 {
		t.Error("Wrong number of term strings, expected 3")
	}
	if terms[2] != "asthma+OR+copd+OR+pulmonary+fibrosis+AND+martin+jg[Au]+AND+o\\'byrne+p[Au]+AND+hamid+q[Au]+AND+2003[pdat]" {
		t.Errorf("Term string failed:%s", terms)
	}
}
