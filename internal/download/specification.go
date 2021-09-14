package download

type Specification struct {
	Size       int64      `json:"size"`
	Definition Definition `json:"definition"`
	URL        string     `json:"url"`
}

type SpecificationList []*Specification

func (r SpecificationList) GetMax() *Specification {
	var size int64
	var pkg *Specification
	for _, v := range r {
		if v.Size > size {
			size = v.Size
			pkg = v
		}
	}
	return pkg
}
