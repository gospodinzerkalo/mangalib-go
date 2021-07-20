package mangalib_go

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"strings"
)


func(m *mangalib) GetManhva() (*[]Resp, error) {
	s, err := m.doRequest("", "")

	res, err := parseBody(s)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return res, nil
}

func parseBody(body io.Reader) (*[]Resp, error) {
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