package mangalib_go

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
)

func (m mangalib) GetGenres() (*GenresResult, error){
	req, err := m.doRequest(BASEURL, http.MethodGet)
	if err != nil {
		return nil, err
	}

	resp, err := parseGenresBody(req)
	if err != nil {
		return nil, err
	}

	return &GenresResult{Genres: resp}, nil
}

func parseGenresBody(body io.Reader) ([]Genres, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	var genres []Genres
	doc.Find(".tags-short a").Each(func(i int, selection *goquery.Selection) {
		title, _ := selection.Attr("title")
		link, _ := selection.Attr("href")
		genres = append(genres, Genres{
			Title: title,
			Link:  link,
		})
	})

	return genres, nil
}

type GenresResult struct {
	Genres []Genres
}

type Genres struct {
	Title 	string
	Link 	string
}