package normalize_test

import (
	"testing"

	"github.com/androiddrew/go-ttsnorm/normalize"
)

func TestOptions(t *testing.T) {
	both := []normalize.Option{normalize.KeepPunctuation(), normalize.ProtectOverrides()}
	cases := []struct {
		in, want string
		opts     []normalize.Option
	}{
		{`"Don't," she said (at 3:05 p.m.).`, `Don t, she said at three oh five p m. .`, nil},
		{`"Don't," she said (at 3:05 p.m.).`, `"Don't," she said (at three oh five p m.).`, []normalize.Option{normalize.KeepPunctuation()}},
		{`Use [Kubernetes](/kˌubəɹnˈɛtiz/) at 10:30 AM.`, `Use [Kubernetes](/kˌubəɹnˈɛtiz/) at ten thirty a m.`, both},
		{`[1.2.3](/wˈʌn/) costs $5`, `[1.2.3](/wˈʌn/) costs five dollars`, both},
		{`Stress [read](-1) 2 [lead](+2), then 3.`, `Stress [read](-1) two [lead](+2), then three.`, both},
		{`Mail bob@example.com or visit https://go.dev at 9:05 p.m.`, `Mail B O B at E X A M P L E dot C O M or visit G O dot D E V at nine oh five PM.`, []normalize.Option{normalize.UpperLetters()}},
		{`See [the docs](https://example.com) or [GPU](/ʤˈipˌijˈu/).`, `See (the docs)(E X A M P L E dot C O M or [GPU](/ʤˈipˌijˈu/).`, []normalize.Option{normalize.KeepPunctuation(), normalize.UpperLetters(), normalize.ProtectOverrides()}},
		{`A + B = C, 3 < 5, # 1 | x * y [z] {w} ₿ & 10%`, `A + B C, three five, one x y (z) (w) & ten percent`, []normalize.Option{normalize.KeepPunctuation()}},
		{`Meet at 3:05 p.m. on May 5. Leave at 9 p.m. Then at 10:00 a.m.`, `Meet at three oh five PM on May five. Leave at nine p.m. Then at ten AM.`, []normalize.Option{normalize.UpperLetters()}},
		{`Meet at 3:05 p.m. on May 5.`, `Meet at three oh five p m. on May five.`, nil},
		{`Without protection [x](/1/) has 1`, `Without protection (x)( one ) has one`, []normalize.Option{normalize.KeepPunctuation()}},
	}
	for _, c := range cases {
		if got := normalize.Text(c.in, c.opts...); got != c.want {
			t.Errorf("Text(%q)\n got: %q\nwant: %q", c.in, got, c.want)
		}
	}
}
