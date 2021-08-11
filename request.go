package mangalib_go

import (
	"errors"
	"io"
	"net/http"
)

const (
	BASEURL = "https://mangalib.me"
	TEST = "https://mangalib.me/seirei-gensouki-konna-sekai-de-deaeta-kimi-ni-minazuki-futago/v5/c31.5?bid=242&page=1"
	IMGURL = "https://img3.cdnlibs.link//manga"
)

func(m mangalib) doRequest(req *http.Request) (io.Reader, error) {

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

//checkImageServer
func checkImageServer(url string, ch chan string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return errors.New("cant get from server")
	}
	ch <- url
	return nil
}