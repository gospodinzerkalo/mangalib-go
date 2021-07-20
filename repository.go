package mangalib_go

import "io"

type Repository interface {
	GetManhva() (*[]Resp, error)

	doRequest(url, method string) (io.Reader, error)
}
