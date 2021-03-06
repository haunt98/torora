package bencoding

import (
	"bufio"
	"fmt"
	"strings"
	"testing"
	"torora/pkg/comparison"
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
		{
			name:  "dict",
			input: "d1:ai1ee",
			wantResult: map[string]interface{}{
				"a": int64(1),
			},
			wantErr: nil,
		},
		{
			name:  "dict",
			input: "d1:ali1eee",
			wantResult: map[string]interface{}{
				"a": []interface{}{
					int64(1),
				},
			},
			wantErr: nil,
		},
		{
			name:  "dict",
			input: "d1:ad1:bi3eee",
			wantResult: map[string]interface{}{
				"a": map[string]interface{}{
					"b": int64(3),
				},
			},
			wantErr: nil,
		},
		{
			name:       "dict",
			input:      "di1ei2ee",
			wantResult: map[string]interface{}{},
			wantErr:    fmt.Errorf("dict key must be string"),
		},
		{
			name:       "dict",
			input:      "d1:bi1e1:ai2ee",
			wantResult: map[string]interface{}{},
			wantErr:    fmt.Errorf("dict key must appear in lexicographical order"),
		},
		{
			name:       "dict",
			input:      "d1:ae",
			wantResult: map[string]interface{}{},
			wantErr:    fmt.Errorf("dict key missing value"),
		},
		{
			name:       "dict",
			input:      "d1:ai2e",
			wantResult: map[string]interface{}{},
			wantErr:    fmt.Errorf("dict missing e"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := bufio.NewReader(strings.NewReader(tc.input))
			decoder := NewDecoder(r)

			gotResult, gotErr := decoder.Decode()
			comparison.CompareError(t, tc.wantErr, gotErr)
			compareResult(t, tc.wantResult, gotResult)
		})
	}
}

func compareResult(t *testing.T, want, got interface{}) {
	switch want.(type) {
	case nil, int64, string:
		comparison.CompareInterface(t, want, got, "result")
	case []interface{}:
		parsedWant, ok := want.([]interface{})
		if !ok {
			t.Errorf("want result is not list")
		}

		parsedGot, ok := got.([]interface{})
		if !ok {
			t.Errorf("got result is not list")
		}

		comparison.CompareInterface(t, len(parsedWant), len(parsedGot), "len result")

		for i := 0; i < len(parsedGot); i += 1 {
			compareResult(t, parsedWant[i], parsedGot[i])
		}
	case map[string]interface{}:
		parsedWant, ok := want.(map[string]interface{})
		if !ok {
			t.Errorf("want result is not dict")
		}

		parsedGot, ok := got.(map[string]interface{})
		if !ok {
			t.Errorf("got result is not dict")
		}

		for wantKey, wantValue := range parsedWant {
			compareResult(t, wantValue, parsedGot[wantKey])
		}
	default:
		t.Errorf("not implement type result")
	}

	return
}
