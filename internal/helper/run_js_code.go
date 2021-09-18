package helper

import (
	"fmt"

	"rogchap.com/v8go"
)

func RunJsCode(code string) (string, error) {
	ctx, err := v8go.NewContext()
	if err != nil {
		return "", err
	}
	val, err := ctx.RunScript(fmt.Sprintf("JSON.stringify(%s)", code), "math.js")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", val), nil
}
