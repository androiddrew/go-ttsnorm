package normalize_test

import (
	"fmt"
	"testing"

	"github.com/androiddrew/go-ttsnorm/corpus"
	"github.com/androiddrew/go-ttsnorm/normalize"
)

func TestCorpus(t *testing.T) {
	cases, err := corpus.Cases()
	if err != nil {
		t.Fatal(err)
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			if got := normalize.Text(c.Input); got != c.Expected {
				source := "Python"
				if c.Override != nil {
					source = "override: " + c.Override.Reason
				}
				t.Errorf("Text(%q)\n got: %q\nwant: %q (%s)", c.Input, got, c.Expected, source)
			}
		})
	}
}

func TestExamples(t *testing.T) {
	for in, want := range map[string]string{
		"2024 budget":  "twenty twenty-four budget",
		"$3.5 million": "three point five million dollars",
		"3 GB":         "three gigabytes",
		"Visit https://example.com or email hello@example.com.": "Visit e x a m p l e dot c o m or email h e l l o at e x a m p l e dot c o m.",
		"Dr. Rivera paid $12.50 at 3:05 p.m.":                   "Doctor Rivera paid twelve dollars and fifty cents at three oh five p m.",
	} {
		if got := normalize.Text(in); got != want {
			t.Errorf("Text(%q) = %q, want %q", in, got, want)
		}
	}
}
