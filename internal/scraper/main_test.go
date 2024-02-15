package scraper

import (
	"testing"
)

func Test_searchAndWriteMetadata(t *testing.T) {
	type args struct {
		gameName string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			searchAndWriteMetadata(tt.args.gameName)
		})
	}
}
