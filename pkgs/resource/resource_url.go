package resource

import (
	"io"
	"net/http"
)

type urlResource struct {
	title string
	url   string
	specs []*Specification
}

func NewURL(title string, url string) Resourcer {
	return &urlResource{
		title: title,
		url:   url,
	}
}

func NewURLWithSpecification(title string, specifications []*Specification) Resourcer {
	return &urlResource{
		title: title,
		specs: specifications,
	}
}

func (r *urlResource) Title() string {
	return r.title
}

func (r *urlResource) Reader() (int64, io.ReadCloser, error) {
	url := r.url
	if len(r.specs) > 0 {
		url = SpecificationList(r.specs).GetMax().URL
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, nil, err
	}

	resp, err := downloadHttpClient.Do(req)
	if err != nil {
		return 0, nil, err
	}

	return resp.ContentLength, resp.Body, nil
}
