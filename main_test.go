package main

import (
	"bytes"
	"strings"
	"testing"
)

func Test_run_(t *testing.T) {
	const (
		noErr  = false
		hasErr = true
	)
	cases := map[string]struct {
		input   string
		args    []string
		want    string
		wantErr bool
	}{
		// absolute path
		"ok/root":              {"/a/b/c/d/e", []string{"1", "3"}, "/a/c", noErr},
		"ok/not-root":          {"/a/b/c/d/e", []string{"2", "3"}, "b/c", noErr},
		"ok/range":             {"/a/b/c/d/e", []string{"1..3"}, "/a/b/c", noErr},
		"ok/range-not-root":    {"/a/b/c/d/e", []string{"2..4"}, "b/c/d", noErr},
		"ok/range-one":         {"/a/b/c/d/e", []string{"3..3"}, "c", noErr},
		"ok/range-two":         {"/a/b/c/d/e", []string{"3..4"}, "c/d", noErr},
		"ok/range-minus-1":     {"/a/b/c/d/e", []string{"2..-1"}, "b/c/d/e", noErr},
		"ok/range-minus-2":     {"/a/b/c/d/e", []string{"2..-2"}, "b/c/d", noErr},
		"ok/range-minus-3":     {"/a/b/c/d/e", []string{"1..-4"}, "/a/b", noErr},
		"ok/range-minus-4":     {"/a/b/c/d/e", []string{"1..-5"}, "/a", noErr},
		"err/range-overflow-1": {"/a/b/c/d/e", []string{"1..-6"}, "", hasErr},
		"err/range-overflow-2": {"/a/b/c/d/e", []string{"3..2"}, "", hasErr},
		"ok/replace":           {"/a/b/c/d/e", []string{"2", "1", "3"}, "b/a/c", noErr},
		"ok/replace-with-nega": {"/a/b/c/d/e", []string{"2", "1", "-1"}, "b/a/e", noErr},
		"ok/repeat":            {"/a/b/c/d/e", []string{"2", "2", "2"}, "b/b/b", noErr},
		"ok/no-arg":            {"/a/b/c/d/e", []string{}, "/a/b/c/d/e", noErr},
		"ok/full":              {"/a/b/c/d/e", []string{"1..-1"}, "/a/b/c/d/e", noErr},
		"ok/range-left-only-1": {"/a/b/c/d/e", []string{"1.."}, "/a/b/c/d/e", noErr},
		"ok/range-left-only-2": {"/a/b/c/d/e", []string{"3.."}, "c/d/e", noErr},
		"ok/single-positive":   {"/a/b/c/d/e", []string{"3"}, "c", noErr},
		"ok/single-negative":   {"/a/b/c/d/e", []string{"-1"}, "e", noErr},
		// relative path
		"ok/rel/simple":  {"x/y/z", []string{"1", "3"}, "x/z", noErr},
		"ok/rel/range":   {"x/y/z", []string{"2", "3"}, "y/z", noErr},
		"ok/rel/dot":     {"./x/y/z", []string{"1", "3"}, "./x/z", noErr},
		"ok/rel/not-dot": {"./x/y/z", []string{"2", "3"}, "y/z", noErr},
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
			case tt.wantErr && err == nil:
				t.Fatal("expected error did not occur")
			case !tt.wantErr && err != nil:
				t.Fatal("unexpected error:", err)
			}

			if got := strings.TrimSuffix(got.String(), "\n"); got != tt.want {
				t.Errorf("want %q, but got %q", tt.want, got)
			}
		})
	}
}
