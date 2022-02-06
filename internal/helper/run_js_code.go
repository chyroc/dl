package helper

import (
	"fmt"

	"rogchap.com/v8go"
)

func RunJsCode(code string) (string, error) {
	ctx := v8go.NewContext()
	val, err := ctx.RunScript(fmt.Sprintf("JSON.stringify(%s)", code), "math.js")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", val), nil
}
