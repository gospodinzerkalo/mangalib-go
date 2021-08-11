package mangalib_go

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
)

func (m mangalib) GetBookmark(user User) (*Bookmark, error) {
	url := fmt.Sprintf("%s/bookmark/%d", BASEURL, user.ID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	resp, err := m.doRequest(req)
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

//TODO
func (m mangalib) GetUserInfo(user User) (*UserInfo, error) {
	url := fmt.Sprintf("%s/user/%d", BASEURL, user.ID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	//req.Header.Add("cookies", "_ga=GA1.2.457774219.1625997498; _ym_uid=162599749918906125; _ym_d=1625997499; _gid=GA1.2.341716699.1628696630; _ym_isad=1; _ym_visorc=b; _gat_gtag_UA_27663466_6=1; remember_web_59ba36addc2b2f9401580f014c7f58ea4e30989d=eyJpdiI6ImpyaE9qdHVvR3JvZUtETitoY3pzMHc9PSIsInZhbHVlIjoiaGJ5RFpicERjWkdLUTZRa0tQZFA0ZWRUVThyNVV4OHRDNDBheW9xMkkwd0lKRTFaR3FvUjd1eEk1cmdrcVF0cnRpdGRjTkV2OXVwbHBmRXFYSDFRbytYdTk2d3JMVEEzYnVLWjhhZ3pNSXViRlZPTXBIN25obTdzRW1kSVpxYThxT3lJQlowNmsxMkxUSFdGenRva0gzNXNhMUllRTVjcmw0SkNtY255VlVialNIOHpDSEViMGFsa2M0MGVVcFpwdm1UNFFXN0svcXlWSks2aUVpRktoeUxjdUZxaDFrKzI1SytJZmNVYnNrWT0iLCJtYWMiOiIzNWQwYTIyYjY0NjJhMzgwNWRmZDdiYmE2NTJlMWQyMzAzMmRhYTAzZDMzNTk3N2RjYTgzNjZlNjJhOGVkZWY4In0%3D; XSRF-TOKEN=eyJpdiI6IjlYMEhUb0c3OGZBam9uOFlET2VsWlE9PSIsInZhbHVlIjoiVVB1VGpERmdYUEkzdmRuYjZkUHAvekNDWk9lcFdrM1o4WWRLWDFDNWMzSWNhdHpKQVg0RHMvbmtxKzdMV3JmR0NCWi9LY3BHaGdYMFh4aUNDL2FVNUc2aWhOVHNyUDVMelhvZWZKTFo0SElreExmNnBsc3UvWitwMldOVTR6VHciLCJtYWMiOiI5NDI1NWU4MzY3NjViMjdkMDJlNmUxZDg0NTlmZjQ4M2I0MjU5NjI4Y2IzMjQyZjBjY2YzYWM4NGM5MTNhNGU5In0%3D; mangalib_session=eyJpdiI6IndMWTUzcUpiQzd4VVJicUVJSEdvYnc9PSIsInZhbHVlIjoiOGtmRS8xZVpGN2d5OUdycjFRbERqUlhqaWpQVi94ZmM0Q0lYOVZBa0pMU3BJZGFZTjlpb2FZTUNDcVFOV1FINDhGTlBrYVBvb3pOT2JKT2hOQXNBUWY2OUE3a09ONmxGVnhmWkJGSVM3ajNMd1FtVU9GSUFUMXBZSmUvMjR3U3UiLCJtYWMiOiI0YTdmZTNhNGMyYzBiZmNiMmQwODBhZTE0NTgxNTYzOTJhNjY3NDMwMDA1MDA1YWQzYzNlZDQwYTI0MmUwYTMwIn0%3D; _count=eyJpdiI6InlyYkVXVWlpODlaSTVjM2RyTWM5TXc9PSIsInZhbHVlIjoiMU1QNk9Kc1dsb1B0aGx6cEdVZ2xZUnEyQWREb3JDWjZnYWpOcERRbWJmcTNpTmtBM3hmR1BsL1c1WWdqc1hheiIsIm1hYyI6ImJhZmVmM2E2ZTdhMzE3YWZiMjhhYmMyOWU4ZGJkMmJkNGUxYTQxYTQ3ZTUxMmFkY2Q1ZTgwYzAxOGZkNDZiYmIifQ%3D%3D")
	resp, err := m.doRequest(req)
	if err != nil {
		return nil, err
	}

	return parseGetUserInfoBody(resp)
}

func parseGetUserInfoBody(body io.Reader) (*UserInfo, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	res := UserInfo{}
	fmt.Println(doc.Text())
	fmt.Println(doc.Find(".profile-user__rank .text-muted").Text())
	return &res, nil
}

type UserInfo struct {
	Level 		int64		`json:"level"`
}