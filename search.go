package mangalib_go

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// types of search
const (
	MANGA = "manga"
	PEOPLE = "people"
	PUBLISHER = "publisher"
	TEAM = "team"
	USER = "user"
)

//Search You can search in english and russian both
type Search struct {
	Q 		string
	Type 	string
}

type SearchResult struct {
	ID 		int64		`json:"id"`
	Slug 	string 		`json:"slug"`
	Cover 	string 		`json:"cover"`
	RusName string 		`json:"rus_name"`
	EngName string 		`json:"eng_name"`
	RateAvg	string 		`json:"rate_avg"`
}

func (m mangalib) Search(search Search) (*[]SearchResult, error) {
	url := fmt.Sprintf("%s/search?type=%s&q=%s", BASEURL, search.Type, search.Q)
	resp, err := m.doRequest(url, http.MethodGet)
	if err != nil {
		return nil, err
	}

	var res []SearchResult
	if err := json.NewDecoder(resp).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}