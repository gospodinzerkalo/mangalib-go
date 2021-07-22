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
	_ = fmt.Sprintf("%s")
	s, err := m.doRequest("", "")

	res, err := parseGetMangaBody(s)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return res, nil
}

func (m *mangalib) GetChapters(manga Manga) (interface{}, error) {
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
			firstChapter.Name = f
			return
		}
	})

	return &firstChapter, nil
}

func parseGetMangaBody(body io.Reader) (*[]Resp, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	imgUrl := "https://img3.cdnlibs.link//manga/seirei-gensouki-konna-sekai-de-deaeta-kimi-ni-minazuki-futago/chapters/"


	id, _  := doc.Find("#comments").Attr("data-post-id")
	fmt.Println(id)
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
				fmt.Sprintf("%s%s/%s", imgUrl, id, img.U),
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