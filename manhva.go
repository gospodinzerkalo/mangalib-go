package mangalib_go

import (
	"encoding/json"
	"errors"
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

	// first get last script in tag head
	headLastScript := doc.Find("head script").Last()
	scripts := strings.Split(strings.TrimSpace(headLastScript.Text()), "\n")
	if len(scripts) < 2 {
		return nil, errors.New("Cat find last scripts in <head>")
	}

	// there are 2 scripts, we need second for parsing image servers
	secondScript := strings.Split(strings.TrimSpace(scripts[1]), "=")
	if len(secondScript) < 2 {
		return nil, errors.New("cant parse script")
	}
	var servers ImageServers
	serversData := strings.TrimSpace(secondScript[1])
	if err := json.Unmarshal([]byte(serversData[:len(serversData) - 1]), &servers); err != nil {
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
			url, err := fetchImagesFromServers(servers.Servers, fmt.Sprintf("/manga/%s/chapters/%s/%s", name, id, img.U))
			if err != nil {
				return nil, err
			}
			res = append(res, Resp{
				url,
			})
		}
	}

	return &res, nil
}

type ImageServers struct {
	Servers 	Servers	`json:"servers"`
}
//Servers Images servers
type Servers struct {
	Main 		string		`json:"main"`
	Secondary 	string 		`json:"secondary"`
	Compress 	string 		`json:"compress"`
	Fourth 		string 		`json:"fourth"`
}

type Resp struct {
	Url 		string 		`json:"url"`
}

type Image struct {
	P 	int64 	`json:"p"`
	U 	string 	`json:"u"`
}

//fetchImagesFromServers check all servers and return available image from server
func fetchImagesFromServers(servers Servers, q string) (string, error) {
	ch := make(chan string, 0)

	go checkImageServer(servers.Main + "/" + q, ch)
	go checkImageServer(servers.Secondary + "/" + q, ch)
	go checkImageServer(servers.Compress + "/" + q, ch)
	go checkImageServer(servers.Fourth + "/" + q, ch)

	for  {
		select {
		case res := <- ch:
			return res, nil
		}
	}
}

