package mangalib_go

import "io"

type Repository interface {
	GetManga() (*[]Resp, error)

	doRequest(url, method string) (io.Reader, error)
}
