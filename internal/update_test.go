package internal

import (
	"testing"
)

func TestComapreVersions(t *testing.T) {
	cases := []struct {
		tag1 string
		tag2 string
		isTag1Newer bool
	}{
		{"v1.0.0", "v1.0.0", false},
		{"v1.0.1", "v1.0.0", true},
		{"v1.0.1", "v1.0.2", false},
		{"v1.0.11", "v1.0.2", true},
	}

	// Run actual tests
	for i, c := range cases {
		r := isNewUpdateLargerTanCurrent(c.tag1, c.tag2)
		if r != c.isTag1Newer {
			t.Fatalf("%d: Wrong result for tag1=%s, tag2=%s, should be %v but is %v\n",
				i, c.tag1, c.tag2, c.isTag1Newer, r)
		}
	}
}