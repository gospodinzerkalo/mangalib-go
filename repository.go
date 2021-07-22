package mangalib_go

import "io"

type Repository interface {
	GetManga(manga Manga) (*[]Resp, error)
	GetChapters(manga Manga) (interface{}, error)

	doRequest(url, method string) (io.Reader, error)
}