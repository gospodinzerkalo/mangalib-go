package mangalib_go

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"strings"
)

type Manga struct {
	Name 	string
}

func(m *mangalib) GetManga(manga Manga) (*[]Resp, error) {
	// firstly parse section chapter page for fetching url of first chapter
	firstCh, err := m.getChapters(manga)
	if err != nil {
		return nil, err
	}

	resp, err := m.doRequest(firstCh.Name, http.MethodGet)
	if err != nil {
		return nil, err
	}

	res, err := parseGetMangaBody(resp, manga.Name)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return res, nil
}

func (m *mangalib) getChapters(manga Manga) (*FirstChapter, error) {
	url := fmt.Sprintf("%s/%s?section=chapters", BASEURL, manga.Name)
	resp, err := m.doRequest(url, http.MethodGet)
	if err != nil {
		return nil, err
	}

	//data, err := io.ReadAll(resp)
	//if err != nil {
	//	return nil, err
	//}
	//fmt.Println(string(data))
	//return nil, nil
	return parseGetMangaSectionChaptersBody(resp)
}

type FirstChapter struct {
	Name 	string
}

func parseGetMangaSectionChaptersBody(body io.Reader) (*FirstChapter, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	//TODO
	//firstChapter, _ := doc.Find(".media-sidebar__buttons .section a").Attr("href")
	firstChapter := FirstChapter{}

	doc.Find("a").Each(func(i int, selection *goquery.Selection) {
		if selection.HasClass("button button_block button_primary") {
			f, _ := selection.Attr("href")
			fmt.Println("First: ", f)
			firstChapter.Name = f
			return
		}
	})

	return &firstChapter, nil
}

func parseGetMangaBody(body io.Reader, name string) (*[]Resp, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	id, _  := doc.Find("#comments").Attr("data-post-id")
	s := doc.Find("#pg").Text()
	l := strings.Split(s, "=")
	var res []Resp
	if len(l) > 1 {
		ll := strings.TrimSpace(l[1])
		ss := ll[:len(ll) - 1]
		var list []Image
		if err := json.Unmarshal([]byte(ss), &list); err != nil {
			return nil, err
		}
		for _, img := range list {
			res = append(res, Resp{
				fmt.Sprintf("%s/%s/chapters/%s/%s", IMGURL, name, id, img.U),
			})
		}
	}

	return &res, nil
}

type Resp struct {
	Url 		string 		`json:"url"`
}

type Image struct {
	P 	int64 	`json:"p"`
	U 	string 	`json:"u"`
}