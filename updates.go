package mangalib_go

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"strings"
)

//GetUpdates - get updates of mangas from main page
func (m mangalib) GetUpdates() (*[]UpdateResult, error) {
	req, err := http.NewRequest(http.MethodGet, BASEURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := m.doRequest(req)
	if err != nil {
		return nil, err
	}

	return parseGetUpdatesBody(resp, &m)
}

type UpdateResult struct {
	Name 		string
	NameRus 	string
	Type 		string
	Link 		string
	UpdatesDate string
	mLib 		Repository
}

func parseGetUpdatesBody(body io.Reader, mLib Repository) (*[]UpdateResult, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	var res []UpdateResult

	doc.Find(".updates__item").Each(func(i int, selection *goquery.Selection) {
		link, _ := selection.Find(".updates__left a").Attr("href")
		typ := selection.Find(".updates__left .updates__type").Text()
		nameRus := selection.Find(".updates__name a").Text()
		date := selection.Find(".updates__date").Text()
		res = append(res, UpdateResult{
			NameRus: nameRus,
			Type: typ,
			Link: link,
			UpdatesDate: date,
			mLib: mLib,
		})
	})

	return &res, nil
}

func (u UpdateResult) Get() (*MangaResponse, error) {
	if u.Link == "" {
		return nil, nil
	}

	lis := strings.Split(u.Link, "/")
	name := lis[len(lis) - 1]
	resp, err := u.mLib.GetManga(Manga{Name: name})
	if err != nil {
		return nil, err
	}
	return resp, nil
}