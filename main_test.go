package main

import (
	"bytes"
	"strings"
	"testing"
)

func Test_run_(t *testing.T) {
	const (
		ok     = false
		hasErr = true
	)
	cases := map[string]struct {
		input  string
		args   []string
		want   string
		hasErr bool
	}{
		// absolute path
		"abs/root":              {"/a/b/c/d/e", []string{"1", "3"}, "/a/c", ok},
		"abs/not-root":          {"/a/b/c/d/e", []string{"2", "3"}, "b/c", ok},
		"abs/range":             {"/a/b/c/d/e", []string{"1..3"}, "/a/b/c", ok},
		"abs/range-not-root":    {"/a/b/c/d/e", []string{"2..4"}, "b/c/d", ok},
		"abs/range-one":         {"/a/b/c/d/e", []string{"3..3"}, "c", ok},
		"abs/range-two":         {"/a/b/c/d/e", []string{"3..4"}, "c/d", ok},
		"abs/range-minus-1":     {"/a/b/c/d/e", []string{"2..-1"}, "b/c/d/e", ok},
		"abs/range-minus-2":     {"/a/b/c/d/e", []string{"2..-2"}, "b/c/d", ok},
		"abs/range-minus-3":     {"/a/b/c/d/e", []string{"1..-4"}, "/a/b", ok},
		"abs/range-minus-4":     {"/a/b/c/d/e", []string{"1..-5"}, "/a", ok},
		"abs/range-overflow-1":  {"/a/b/c/d/e", []string{"1..-6"}, "", hasErr},
		"abs/range-overflow-2":  {"/a/b/c/d/e", []string{"3..2"}, "", hasErr},
		"abs/replace":           {"/a/b/c/d/e", []string{"2", "1", "3"}, "b/a/c", ok},
		"abs/replace-with-nega": {"/a/b/c/d/e", []string{"2", "1", "-1"}, "b/a/e", ok},
		"abs/repeat":            {"/a/b/c/d/e", []string{"2", "2", "2"}, "b/b/b", ok},
		"abs/no-arg":            {"/a/b/c/d/e", []string{}, "", hasErr},
		"abs/full":              {"/a/b/c/d/e", []string{"1..-1"}, "/a/b/c/d/e", ok},
		"abs/range-left-only-1": {"/a/b/c/d/e", []string{"1.."}, "/a/b/c/d/e", ok},
		"abs/range-left-only-2": {"/a/b/c/d/e", []string{"3.."}, "c/d/e", ok},
		"abs/single-positive":   {"/a/b/c/d/e", []string{"3"}, "c", ok},
		"abs/single-negative":   {"/a/b/c/d/e", []string{"-1"}, "e", ok},
		// relative path
		"rel/simple":  {"x/y/z", []string{"1", "3"}, "x/z", ok},
		"rel/range":   {"x/y/z", []string{"2", "3"}, "y/z", ok},
		"rel/dot":     {"./x/y/z", []string{"1", "3"}, "./x/z", ok},
		"rel/not-dot": {"./x/y/z", []string{"2", "3"}, "y/z", ok},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var got bytes.Buffer
			cli := &CLI{
				Stdout: &got,
				Stderr: &got,
				Stdin:  strings.NewReader(tt.input),
			}
			err := cli.run(tt.args)

			switch {
			case tt.hasErr && err == nil:
				t.Fatal("expected error did not occur")
			case !tt.hasErr && err != nil:
				t.Fatal("unexpected error:", err)
			}

			if got := strings.TrimSuffix(got.String(), "\n"); got != tt.want {
				t.Errorf("want %q, but got %q", tt.want, got)
			}
		})
	}
}
