package resource

import (
	"fmt"
	"strings"
)

// https://www.image-engineering.de/library/technotes/991-separating-sd-hd-full-hd-4k-and-8k
type Definition string

const (
	DefinitionSD     Definition = "sd"
	DefinitionHD     Definition = "hd"
	DefinitionFullHD Definition = "full-hd"
	DefinitionUHD    Definition = "uhd"
	Definition4K     Definition = "4k"
	Definition8K     Definition = "8k"
)

func MayConvertDefinition(s string) Definition {
	res, _ := ConvertDefinition(s)
	return res
}

func ConvertDefinition(s string) (Definition, error) {
	switch strings.TrimSpace(strings.ToLower(s)) {
	case "sd", "mp4sd", "240p", "360p":
		return DefinitionSD, nil
	case "hd", "mp4hd", "mp4hd2v2", "720p":
		return DefinitionHD, nil
	case "uhd":
		return DefinitionUHD, nil
	case "fhd", "fullhd", "full-hd", "1080p":
		return DefinitionFullHD, nil
	case "4k", "2160p":
		return Definition4K, nil
	case "8k":
		return Definition8K, nil
	}
	return "", fmt.Errorf("%q is invalid definition", s)
}
