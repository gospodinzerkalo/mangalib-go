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

func(m *mangalib) GetManga(manga Manga) (*MangaResponse, error) {
	// firstly parse section chapter page for fetching url of first chapter
	firstCh, err := m.getChapters(manga)
	if err != nil {
		return nil, err
	}

	resp, err := m.doRequest(firstCh.Name, http.MethodGet)
	if err != nil {
		return nil, err
	}

	res, mInfo, err := parseGetMangaBody(resp, manga.Name)
	if err != nil {
		return nil, err
	}

	respManga := MangaResponse{Urls: *res, page: 1, url: firstCh.Name, name: manga.Name, next: mInfo.Next.Url}

	return &respManga, nil
}

// TODO
func (mr *MangaResponse) NextChapter() (*MangaResponse, error){
	//if mr.page == 0 {
	//	return nil, nil
	//}
	//
	//url := strings.Split(mr.url, "/")
	//if len(url) < 1 {
	//	return nil, nil
	//}
	//
	//newUrl := strings.Join(url[:len(url) - 1], "/")
	//newUrl = fmt.Sprintf("%s/c%d", newUrl, mr.page)
	//fmt.Println(newUrl)
	if mr.next == "" {
		return nil, nil
	}
	resp, err := http.Get(mr.next)
	if err != nil {
		return nil, err
	}

	res, _, err := parseGetMangaBody(resp.Body, mr.name)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	respManga := MangaResponse{Urls: *res, page: mr.page + 1, name: mr.name, url: mr.url}

	return &respManga, nil
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

func parseGetMangaBody(body io.Reader, name string) (*[]Resp, *mangaInfo, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, nil, err
	}

	// first get last script in tag head
	headLastScript := doc.Find("head script").Last()
	scripts := strings.Split(strings.TrimSpace(headLastScript.Text()), "\n")
	if len(scripts) < 2 {
		return nil, nil, errors.New("Cat find last scripts in <head>")
	}

	// there are 2 scripts, we need second for parsing image servers
	secondScript := strings.Split(strings.TrimSpace(scripts[1]), "=")
	if len(secondScript) < 2 {
		return nil, nil, errors.New("cant parse script")
	}
	var servers mangaInfo
	serversData := strings.TrimSpace(secondScript[1])
	if err := json.Unmarshal([]byte(serversData[:len(serversData) - 1]), &servers); err != nil {
		fmt.Println("ERR 1", secondScript)
		return nil, nil, err
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
			fmt.Println("ERR 2")
			return nil, nil, err
		}
		for _, img := range list {
			url, err := fetchImagesFromServers(servers.Servers, fmt.Sprintf("/manga/%s/chapters/%s/%s", name, id, img.U))
			if err != nil {
				return nil, nil, err
			}
			res = append(res, Resp{
				url,
			})
		}
	}

	return &res, &servers, nil
}

type mangaInfo struct {
	Servers 	Servers	`json:"servers"`
	Next 		next 	`json:"next"`
}

type next struct {
	Url 	string	`json:"url"`
}

//Servers Images servers
type Servers struct {
	Main 		string		`json:"main"`
	Secondary 	string 		`json:"secondary"`
	Compress 	string 		`json:"compress"`
	Fourth 		string 		`json:"fourth"`
}

type MangaResponse struct {
	page 		int64
	name 		string
	url 		string
	next 		string
	Urls		[]Resp
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
	//go checkImageServer(servers.Fourth + "/" + q, ch)

	for  {
		select {
		case res := <- ch:
			return res, nil
		}
	}
}

