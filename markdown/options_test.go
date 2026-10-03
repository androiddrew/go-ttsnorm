package markdown_test

import (
	"testing"

	"github.com/androiddrew/go-ttsnorm/markdown"
)

func TestKeepOverrides(t *testing.T) {
	cases := []struct{ in, plain, keep string }{
		{"Use [Kubernetes](/kˌubəɹnˈɛtiz/) **today**.", "Use Kubernetes today.", "Use [Kubernetes](/kˌubəɹnˈɛtiz/) today."},
		{"I [read](-1) it and [lead](+2) on.", "I read it and lead on.", "I [read](-1) it and [lead](+2) on."},
		{"[half](0.5) and [Hello](#greeting#)", "half and Hello", "[half](0.5) and [Hello](#greeting#)"},
		{"See [the docs](https://example.com/a/) now.", "See the docs now.", "See the docs now."},
		{"- [GPU](/ʤˈipˌijˈu/) setup", "GPU setup.", "[GPU](/ʤˈipˌijˈu/) setup."},
	}
	for _, c := range cases {
		if got := markdown.ToSpeech(c.in); got != c.plain {
			t.Errorf("ToSpeech(%q) = %q, want %q", c.in, got, c.plain)
		}
		if got := markdown.ToSpeech(c.in, markdown.KeepOverrides()); got != c.keep {
			t.Errorf("ToSpeech(%q, KeepOverrides) = %q, want %q", c.in, got, c.keep)
		}
	}
}
