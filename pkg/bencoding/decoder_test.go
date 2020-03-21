package bencoding

import (
	"bufio"
	"fmt"
	"strings"
	"testing"
)

func TestDecoder_Decode(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantResult interface{}
		wantErr    error
	}{
		{
			name:       "integer",
			input:      "i10e",
			wantResult: int64(10),
			wantErr:    nil,
		},
		{
			name:       "integer",
			input:      "i9",
			wantResult: int64(0),
			wantErr:    fmt.Errorf("integer missing e"),
		},
		{
			name:       "integer",
			input:      "ie",
			wantResult: int64(0),
			wantErr:    fmt.Errorf("integer missing value"),
		},
		{
			name:       "string",
			input:      "4:spam",
			wantResult: "spam",
			wantErr:    nil,
		},
		{
			name:       "string",
			input:      "4:spa",
			wantResult: "",
			wantErr:    fmt.Errorf("string missing character"),
		},
		{
			name:  "list",
			input: "li1ei2ee",
			wantResult: []interface{}{
				int64(1),
				int64(2),
			},
			wantErr: nil,
		},
		{
			name:  "list",
			input: "l4:spami100ee",
			wantResult: []interface{}{
				"spam",
				int64(100),
			},
			wantErr: nil,
		},
		{
			name:       "list",
			input:      "li1e",
			wantResult: []interface{}{},
			wantErr:    fmt.Errorf("list missing e"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := bufio.NewReader(strings.NewReader(tc.input))
			decoder := NewDecoder(r)

			gotResult, gotErr := decoder.Decode()

			if tc.wantErr == nil {
				if gotErr != nil {
					t.Errorf("want err [%+v] got err [%+v]\n", tc.wantErr, gotErr)
				}
			} else {
				if gotErr == nil {
					t.Errorf("want err [%+v] got err [%+v]\n", tc.wantErr, gotErr)
				}
				if gotErr.Error() != tc.wantErr.Error() {
					t.Errorf("want err [%+v] got err [%+v]\n", tc.wantErr, gotErr)
				}
			}

			switch tc.wantResult.(type) {
			case int64, string:
				if gotResult != tc.wantResult {
					t.Errorf("want result [%+v] got result [%+v]\n", tc.wantResult, gotResult)
				}
			case []interface{}:
				parsedWant, ok := tc.wantResult.([]interface{})
				if !ok {
					t.Errorf("want result is not list")
				}
				parsedGot, ok := gotResult.([]interface{})
				if !ok {
					t.Errorf("got result is not list")
				}

				if len(parsedGot) != len(parsedWant) {
					t.Errorf("want result len [%d] got result len [%d]", len(parsedWant), len(parsedGot))
				}

				for i := 0; i < len(parsedWant); i += 1 {
					if parsedGot[i] != parsedWant[i] {
						t.Errorf("want result [%+v] got result [%+v]\n", parsedWant[i], parsedGot[i])
					}
				}
			default:
				t.Errorf("not implement testcase")
			}
		})
	}
}
