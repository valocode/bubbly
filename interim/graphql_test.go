package interim

import "testing"

func TestNewEngine(t *testing.T) {
	cases := []struct {
		desc string
	}{
		{
			desc: "",
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
		})
	}
}
