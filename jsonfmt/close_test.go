package main

import "testing"

func TestCloseJSON(t *testing.T) {
	cases := []struct {
		in  string
		out string
	}{
		{
			in:  `{"foo":42}`,
			out: `{"foo":42}`,
		},
		{
			in:  `{"foo":42`,
			out: `{"foo":42}`,
		},
		{
			in:  `[{"foo":42`,
			out: `[{"foo":42}]`,
		},
		{
			in:  `[{"foo":[42,32`,
			out: `[{"foo":[42,32]}]`,
		},
		{
			in:  `[{"foo":"bar"`,
			out: `[{"foo":"bar"}]`,
		},
		{
			in:  `[{"foo":"bar`,
			out: `[{"foo":"bar"}]`,
		},
		{
			in:  `[{"foo":"b\"ar`,
			out: `[{"foo":"b\"ar"}]`,
		},
	}

	for i, tc := range cases {
		out := tc.in + string(closeJSON([]byte(tc.in)))
		if out != tc.out {
			t.Errorf("%d: Close %s => %s", i, tc.in, out)
		}
	}
}
