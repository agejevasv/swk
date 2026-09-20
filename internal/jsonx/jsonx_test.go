package jsonx

import (
	"encoding/json"
	"testing"
)

func TestDecode_PreservesNumberLiterals(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"integer beyond float64 precision", `{"v":9007199254740993}`, "9007199254740993"},
		{"integer beyond int64", `{"v":12345678901234567890}`, "12345678901234567890"},
		{"round number stays plain", `{"v":1000000}`, "1000000"},
		{"exponent notation kept", `{"v":1e10}`, "1e10"},
		{"long decimal kept", `{"v":0.1234567890123456789}`, "0.1234567890123456789"},
		{"negative", `{"v":-42}`, "-42"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got map[string]any
			if err := Decode([]byte(tt.input), &got); err != nil {
				t.Fatalf("Decode: %v", err)
			}
			n, ok := got["v"].(json.Number)
			if !ok {
				t.Fatalf("expected json.Number, got %T (%v)", got["v"], got["v"])
			}
			if n.String() != tt.want {
				t.Errorf("Decode() = %s, want %s", n, tt.want)
			}
		})
	}
}

func TestDecode_RoundTripsThroughMarshal(t *testing.T) {
	const input = `{"id":12345678901234567890,"n":1000000}`

	var v any
	if err := Decode([]byte(input), &v); err != nil {
		t.Fatalf("Decode: %v", err)
	}
	out, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if string(out) != input {
		t.Errorf("round trip = %s, want %s", out, input)
	}
}

func TestDecode_RejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"garbage", "{invalid}"},
		{"trailing comma", `{"a":1,}`},
		{"trailing content", `{"a":1} trailing`},
		{"two top level values", `{"a":1}{"b":2}`},
		{"unterminated", `{"a":`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var v any
			if err := Decode([]byte(tt.input), &v); err == nil {
				t.Errorf("expected an error for %q", tt.input)
			}
		})
	}
}

func TestDecode_AcceptsSurroundingWhitespace(t *testing.T) {
	var v any
	if err := Decode([]byte("  {\"a\": 1}\n\n"), &v); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
