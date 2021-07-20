package mangalib_go

import (
	"io"
	"net/http"
)

const (
	BASEURL = "https://mangalib.me/"
	TEST = "https://mangalib.me/seirei-gensouki-konna-sekai-de-deaeta-kimi-ni-minazuki-futago/v5/c31.5?bid=242&page=1"
)

func(m mangalib) doRequest(url, method string) (io.Reader, error) {
	resp, err := http.Get(TEST)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

