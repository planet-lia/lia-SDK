package internal

import (
	"testing"
)

func TestCompareVersions(t *testing.T) {
	cases := []struct {
		tag1        string
		tag2        string
		isTag1Newer bool
	}{
		{"v1.0.0", "v1.0.0", false},
		{"v1.0.1", "v1.0.0", true},
		{"v1.0.1", "v1.0.2", false},
		{"v1.0.11", "v1.0.2", true},
		{"v2.0.11", "v1.0.2", false},
		{"v2.1.11", "v2.0.22", true},
		{"v2.0.22", "v2.3.1", false},
		{"v2.33.22", "v2.3.1", true},
		{"v1.0.21", "v1.0.20", true},
		{"v11.0.21", "v1.0.20", false},
	}

	// Run actual tests
	for i, c := range cases {
		r := isNewUpdateLargerThanCurrent(c.tag1, c.tag2)
		if r != c.isTag1Newer {
			t.Fatalf("%d: Wrong result for tag1=%s, tag2=%s, should be %v but is %v\n",
				i, c.tag1, c.tag2, c.isTag1Newer, r)
		}
	}
}
