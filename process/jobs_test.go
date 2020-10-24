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
