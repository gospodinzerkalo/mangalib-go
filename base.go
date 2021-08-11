package mangalib_go

import (
	"crypto/tls"
	"net/http"
)

type mangalib struct {
	client 		http.Client
}


func NewMangalib() Repository {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{Transport: tr}

	return &mangalib{client: client}
}



