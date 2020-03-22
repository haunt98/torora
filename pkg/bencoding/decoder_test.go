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
			name:  "list",
			input: "li1eli2eee",
			wantResult: []interface{}{
				int64(1),
				[]interface{}{
					int64(2),
				},
			},
			wantErr: nil,
		},
		{
			name:  "list",
			input: "li1eli2eli3eeee",
			wantResult: []interface{}{
				int64(1),
				[]interface{}{
					int64(2),
					[]interface{}{
						int64(3),
					},
				},
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
			compareErr(t, tc.wantErr, gotErr)
			compareResult(t, tc.wantResult, gotResult)
		})
	}
}

func compareErr(t *testing.T, want, got error) {
	if want == nil && got == nil {
		return
	}

	if want != nil && got != nil {
		if got.Error() != want.Error() {
			t.Errorf("want err [%+v] got err [%+v]\n", want, got)
		}
		return
	}

	t.Errorf("want err [%+v] got err [%+v]\n", want, got)
}

func compareResult(t *testing.T, want, got interface{}) {
	switch want.(type) {
	case int64, string:
		if got != want {
			t.Errorf("want result [%+v] got result [%+v]\n", want, got)
		}
	case []interface{}:
		parsedWant, ok := want.([]interface{})
		if !ok {
			t.Errorf("want result is not list")
		}

		parsedGot, ok := got.([]interface{})
		if !ok {
			t.Errorf("got result is not list")
		}

		if len(parsedGot) != len(parsedWant) {
			t.Errorf("want result len [%d] got result len [%d]", len(parsedWant), len(parsedGot))
		}

		for i := 0; i < len(parsedWant); i += 1 {
			compareResult(t, parsedWant[i], parsedGot[i])
		}
	default:
		t.Errorf("not implement type result")
	}
}
