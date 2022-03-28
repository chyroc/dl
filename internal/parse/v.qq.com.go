package parse

import (
	"github.com/chyroc/dl/internal/parse/adapter/qqvideo"
	"github.com/chyroc/dl/internal/resource"
)

func NewVQqCom() Parser {
	return &vQqCom{}
}

type vQqCom struct{}

func (r *vQqCom) Kind() string {
	return "v.qq.com"
}

func (r *vQqCom) ExampleURLs() []string {
	return []string{"https://v.qq.com/txp/iframe/player.html?vid=j0822mqey5h"}
}

func (r *vQqCom) Parse(uri string) (resource.Resource, error) {
	video, err := qqvideo.New().Extract(uri)
	if err != nil {
		return nil, err
	}

	spec := []*resource.Specification{}
	for specType, v := range video.Streams {
		if len(v.Parts) != 1 {
			panic("v.qq.com: len(v.Parts)!=1")
		}
		spec = append(spec, &resource.Specification{
			Size:       v.Size,
			Definition: resource.MayConvertDefinition(specType),
			URL:        v.Parts[0].URL,
		})
	}
	return resource.NewURLWithSpecification(video.Title, spec), nil
}
