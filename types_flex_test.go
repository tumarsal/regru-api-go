package regru

import (
	"encoding/json"
	"testing"
)

func TestResourceRecord_prioNumberOrString(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		want string
	}{
		{"number", `{"subname":"@","rectype":"MX","content":"mail.example.ru","prio":10}`, "10"},
		{"string", `{"subname":"@","rectype":"MX","content":"mail.example.ru","prio":"10"}`, "10"},
		{"float", `{"subname":"@","rectype":"MX","content":"mail.example.ru","prio":10.0}`, "10"},
		{"absent", `{"subname":"www","rectype":"A","content":"1.2.3.4"}`, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var r ResourceRecord
			if err := json.Unmarshal([]byte(tc.raw), &r); err != nil {
				t.Fatal(err)
			}
			if got := r.Priority.String(); got != tc.want {
				t.Fatalf("prio: got %q, want %q", got, tc.want)
			}
		})
	}
}
