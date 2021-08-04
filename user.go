package mangalib_go

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (m mangalib) GetBookmark(user User) (*Bookmark, error) {
	url := fmt.Sprintf("%s/bookmark/%d", BASEURL, user.ID)
	resp, err := m.doRequest(url, http.MethodGet)
	if err != nil {
		return nil, err
	}

	var bookmark bookmarkResponse
	if err := json.NewDecoder(resp).Decode(&bookmark); err != nil {
		return nil, err
	}

	var res Bookmark
	for _, b := range bookmark.Items {
		switch b.Status {
		case 1:
			res.Reading = append(res.Reading, b)
		case 2:
			res.Plan = append(res.Plan, b)
		case 3:
			res.Thrown = append(res.Thrown, b)
		case 4:
			res.Read = append(res.Read, b)
		case 5:
			res.Loved = append(res.Loved, b)
		}
	}

	return &res, nil
}

type User struct {
	ID 		int64
}

type bookmarkResponse struct {
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

type Bookmark struct {
	Read 	[]Item		`json:"read"`
	Reading []Item		`json:"reading"`
	Thrown 	[]Item		`json:"thrown"`
	Loved	[]Item		`json:"loved"`
	Plan 	[]Item		`json:"plan"`
}