package metainfo

import (
	"fmt"
	"testing"
	"torora/pkg/comparison"
)

func TestParseFile(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		wantFile File
		wantErr  error
	}{
		{
			name: "file",
			input: map[string]interface{}{
				"length": int64(1),
				"path":   "a",
			},
			wantFile: File{
				Length: 1,
				Path:   "a",
			},
			wantErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotFile, gotErr := parseFile(tc.input)
			comparison.CompareError(t, tc.wantErr, gotErr)
			compareFile(t, tc.wantFile, gotFile)
		})
	}
}

func TestParseInfo(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		wantInfo Info
		wantErr  error
	}{
		{
			name: "info",
			input: map[string]interface{}{
				"piece length": int64(1),
				"pieces":       "a",
				"length":       int64(2),
			},
			wantInfo: Info{
				PieceLength: 1,
				Pieces:      "a",
				Length:      2,
				Files:       nil,
			},
			wantErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotInfo, gotErr := parseInfo(tc.input)
			comparison.CompareError(t, tc.wantErr, gotErr)
			compareInfo(t, tc.wantInfo, gotInfo)
		})
	}
}

func TestParseMetainfo(t *testing.T) {
	tests := []struct {
		name         string
		input        interface{}
		wantMetainfo Metainfo
		wantErr      error
	}{
		{
			name: "metainfo",
			input: map[string]interface{}{
				"announce": "a",
				"info": map[string]interface{}{
					"piece length": int64(1),
					"pieces":       "b",
					"length":       int64(2),
				},
			},
			wantMetainfo: Metainfo{
				Announce: "a",
				Info: Info{
					PieceLength: 1,
					Pieces:      "b",
					Length:      2,
					Files:       nil,
				},
			},
			wantErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotMetainfo, gotErr := parseMetainfo(tc.input)
			comparison.CompareError(t, tc.wantErr, gotErr)
			compareMetainfo(t, tc.wantMetainfo, gotMetainfo)
		})
	}
}

func compareFile(t *testing.T, want, got File) {
	comparison.CompareInterface(t, want.Length, got.Length, "file/length")
	comparison.CompareInterface(t, want.Path, got.Path, "file/path")
}

func compareInfo(t *testing.T, want, got Info) {
	comparison.CompareInterface(t, want.PieceLength, got.PieceLength, "info/piece length")
	comparison.CompareInterface(t, want.Pieces, got.Pieces, "info/pieces")
	comparison.CompareInterface(t, want.Length, got.Length, "info/length")

	comparison.CompareInterface(t, len(want.Files), len(got.Files), "len info/files")
	for i := range got.Files {
		comparison.CompareInterface(t, want.Files[i], got.Files[i], fmt.Sprintf("info/files[%d]", i))
	}
}

func compareMetainfo(t *testing.T, want, got Metainfo) {
	comparison.CompareInterface(t, want.Announce, got.Announce, "metainfo/announce")
	compareInfo(t, want.Info, got.Info)
}
