package attendee

import (
	"bytes"
	"encoding/xml"
	"flag"
	"io/ioutil"
	"path/filepath"
	"testing"
)

var update = flag.Bool("update", false, "update .golden files")

func TestAttendeeMarshal(t *testing.T) {
	tests := []struct {
		fileName string
		a        Attendee
	}{
		{
			fileName: "attendee1.golden",
			a: Attendee{
				Name:  "robpike",
				Age:   60,
				Phone: "13912345678",
			},
		},
	}

	for _, tt := range tests {
		got, err := xml.MarshalIndent(&tt.a, "", "  ")
		if err != nil {
			t.Fatalf("want nil, got %v", err)
		}

		golden := filepath.Join("testdata", tt.fileName)
		if *update {
			ioutil.WriteFile(golden, got, 0644)
		}

		want, err := ioutil.ReadFile(golden)
		if err != nil {
			t.Fatalf("open file %s failed: %v", tt.fileName, err)
		}

		if !bytes.Equal(got, want) {
			t.Errorf("want %s, got %s", string(want), string(got))
		}
	}
}
