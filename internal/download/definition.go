package download

import (
	"fmt"
	"strings"
)

func MayConvertDefinition(s string) Definition {
	res, _ := ConvertDefinition(s)
	return res
}

func ConvertDefinition(s string) (Definition, error) {
	switch strings.TrimSpace(strings.ToLower(s)) {
	case "sd":
		return DefinitionSD, nil
	case "hd":
		return DefinitionHD, nil
	case "uhd":
		return DefinitionUHD, nil
	case "fhd", "fullhd", "full-hd":
		return DefinitionFullHD, nil
	case "4k":
		return Definition4K, nil
	case "8k":
		return Definition8K, nil
	}
	return "", fmt.Errorf("%q is invalid definition", s)
}
