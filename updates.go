package mangalib_go

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
)

//GetUpdates - get updates of mangas from main page
func (m mangalib) GetUpdates() (*[]UpdateResult, error) {
	req, err := m.doRequest(BASEURL, http.MethodGet)
	if err != nil {
		return nil, err
	}

	return parseGetUpdatesBody(req)
}

type UpdateResult struct {
	Name 		string
	Type 		string
	Link 		string
}

func parseGetUpdatesBody(body io.Reader) (*[]UpdateResult, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	var res []UpdateResult

	doc.Find(".updates__item").Each(func(i int, selection *goquery.Selection) {
		link, _ := selection.Find(".updates__left a").Attr("href")
		typ := selection.Find(".updates__left .updates__type").Text()
		name := selection.Find(".updates__name a").Text()
		res = append(res, UpdateResult{
			Name: name,
			Type: typ,
			Link: link,
		})
	})

	return &res, nil
}
