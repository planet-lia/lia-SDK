package analytics

import (
	"testing"
	"github.com/liagame/lia-SDK/internal/config"
)

func TestAnalytics(t *testing.T) {
	cases := []struct {
		os string
		in string
		want string
	}{
		{os: "", in: "/home/marko/bot1", want: "bot1"},
		{os: "", in: "bot1", want: "bot1"},
		{os: "", in: "mike/bot1", want: "bot1"},
		{os: "windows", in: "C:\\\\mirko\\bot1", want: "bot1"},
		{os: "windows", in: "mirko\\bot1", want: "bot1"},
		{os: "windows", in: "bot1", want: "bot1"},
	}

	for _, c := range cases {
		config.OperatingSystem = c.os
		out := TrimPath(c.in)
		if out != c.want {
			t.Fatalf("%s != %s", out, c.want)
		}
	}
}