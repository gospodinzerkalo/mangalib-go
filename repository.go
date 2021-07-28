package mangalib_go

import "io"

type Repository interface {
	GetManga(manga Manga) (*MangaResponse, error)
	getChapters(manga Manga) (*FirstChapter, error)
	Search(search Search) (*[]SearchResult, error)
	GetUpdates() (*[]UpdateResult, error)

	doRequest(url, method string) (io.Reader, error)
}