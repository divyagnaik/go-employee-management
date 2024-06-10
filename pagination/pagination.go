package pagination

import (
	"strconv"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http"
)

type Meta struct {
	Page Page `json:"page"`
}

type Page struct {
	Offset int `json:"offset"`
	Size   int `json:"size"`
}

const (
	maxItems           = -1
	offset             = 0
	pageSizeLowerLimit = 0
	PageSizeUpperLimit = 20
	PageParam          = "page."
	SizeParam          = "size"
	OffsetParam        = "offset"
)

func initialisePage(pageOffset, size string) (*Page, error) {
	p := &Page{Size: PageSizeUpperLimit, Offset: offset}

	if pageOffset != "" {
		n, err := strconv.Atoi(pageOffset)
		if err != nil {
			return nil, http.ErrorInvalidParam{Params: []string{PageParam + OffsetParam}}
		}

		p.Offset = n
	}

	if size != "" && size != "-1" {
		n, err := strconv.Atoi(size)
		if err != nil {
			return nil, http.ErrorInvalidParam{Params: []string{PageParam + SizeParam}}
		}

		p.Size = n
	}

	return p, nil
}

func (p Page) validate() error {
	if p.Size < maxItems || p.Size == pageSizeLowerLimit {
		return http.ErrorInvalidParam{Params: []string{PageParam + SizeParam}}
	}

	if p.Offset < offset {
		return http.ErrorInvalidParam{Params: []string{PageParam + OffsetParam}}
	}

	return nil
}

func Pagination(ctx *gofr.Context) (*Page, error) {
	pageSize := ctx.Param("page.size")

	pageOffset := ctx.Param("page.offset")

	page, err := initialisePage(pageOffset, pageSize)
	if err != nil {
		return nil, err
	}

	err = page.validate()
	if err != nil {
		return nil, err
	}

	return page, nil
}
