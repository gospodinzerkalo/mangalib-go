package mangalib_go

import (
	"io"
	"net/http"
)

type Repository interface {
	GetManga(manga Manga) (*MangaResponse, error)
	getChapters(manga Manga) (*FirstChapter, error)
	Search(search Search) (*[]SearchResult, error)
	GetUpdates() (*[]UpdateResult, error)
	GetGenres() (*GenresResult, error)
	GetBookmark(user User) (*Bookmark, error)
	GetUserInfo(user User) (*UserInfo, error)

	doRequest(req *http.Request) (io.Reader, error)
}