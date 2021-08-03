package mangalib_go

import "io"

type Repository interface {
	GetManga(manga Manga) (*MangaResponse, error)
	getChapters(manga Manga) (*FirstChapter, error)
	Search(search Search) (*[]SearchResult, error)
	GetUpdates() (*[]UpdateResult, error)
	GetGenres() (*GenresResult, error)
	GetBookmark(user User) (*BookmarkResponse, error)

	doRequest(url, method string) (io.Reader, error)
}