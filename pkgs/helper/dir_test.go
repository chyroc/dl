package helper

import (
	"strings"
	"testing"

	"github.com/chyroc/go-assert"
)

func TestResolveDirOrCurrent(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    string
		wantErr bool
	}{
		{"1", "", "/dl/pkgs/helper", false},
		{"1", ".", "/dl/pkgs/helper", false},
		{"1", "./", "/dl/pkgs/helper", false},
		{"1", "..", "/dl/pkgs", false},
		{"1", "../", "/dl/pkgs", false},
		{"1", "./x", "/dl/pkgs/helper/x", false},
		{"1", "x", "/dl/pkgs/helper/x", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveDirOrCurrent(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveDirOrCurrent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.HasSuffix(got, tt.want) {
				t.Errorf("ResolveDirOrCurrent() got = %v, want %v", got, tt.want)
			}
		})
	}
	{
		res, _ := ResolveDirOrCurrent("/abc")
		assert.Equal(t, "/abc", res)
	}
}
