package mangalib_go

import (
	"io"
	"net/http"
)

const (
	BASEURL = "https://mangalib.me"
	TEST = "https://mangalib.me/seirei-gensouki-konna-sekai-de-deaeta-kimi-ni-minazuki-futago/v5/c31.5?bid=242&page=1"
	IMGURL = "https://img3.cdnlibs.link//manga"
)

func(m mangalib) doRequest(url, method string) (io.Reader, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}