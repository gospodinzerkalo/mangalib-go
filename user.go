package mangalib_go

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (m mangalib) GetBookmark(user User) (*BookmarkResponse, error) {
	url := fmt.Sprintf("%s/bookmark/%d", BASEURL, user.ID)
	resp, err := m.doRequest(url, http.MethodGet)
	if err != nil {
		return nil, err
	}

	var bookmark BookmarkResponse
	if err := json.NewDecoder(resp).Decode(&bookmark); err != nil {
		return nil, err
	}

	return &bookmark, nil
}

type User struct {
	ID 		int64
}

type BookmarkResponse struct {
	Items	[]Item 	`json:"items"`
}

type Item struct {
	MangaName 		string 		`json:"manga_name"`
	RusName 		string 		`json:"rus_name"`
	Slug 			string 		`json:"slug"`
	MangaId 		int64		`json:"manga_id"`
	Cover 			string 		`json:"cover"`
	SiteId			int64		`json:"site_id"`
	BookId 			int64		`json:"book_id"`
	Status 			int64		`json:"status"`
	CreatedAt 		string 		`json:"created_at"`
}