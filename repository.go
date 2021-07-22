package mangalib_go

import "io"

type Repository interface {
	GetManga(manga Manga) (*[]Resp, error)
	getChapters(manga Manga) (*FirstChapter, error)
	SearchManga(search Search) (*[]SearchResult, error)

	doRequest(url, method string) (io.Reader, error)
}