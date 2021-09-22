package helper

import (
	"time"

	"github.com/briandowns/spinner"
)

var Spinner = spinner.New(spinner.CharSets[35], 100*time.Millisecond)
