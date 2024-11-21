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
		"abs/root":              {"/a/b/c/d/e", []string{"1", "3"}, "/a/c", noErr},
		"abs/not-root":          {"/a/b/c/d/e", []string{"2", "3"}, "b/c", noErr},
		"abs/range":             {"/a/b/c/d/e", []string{"1..3"}, "/a/b/c", noErr},
		"abs/range-not-root":    {"/a/b/c/d/e", []string{"2..4"}, "b/c/d", noErr},
		"abs/range-one":         {"/a/b/c/d/e", []string{"3..3"}, "c", noErr},
		"abs/range-two":         {"/a/b/c/d/e", []string{"3..4"}, "c/d", noErr},
		"abs/range-minus-1":     {"/a/b/c/d/e", []string{"2..-1"}, "b/c/d/e", noErr},
		"abs/range-minus-2":     {"/a/b/c/d/e", []string{"2..-2"}, "b/c/d", noErr},
		"abs/range-minus-3":     {"/a/b/c/d/e", []string{"1..-4"}, "/a/b", noErr},
		"abs/range-minus-4":     {"/a/b/c/d/e", []string{"1..-5"}, "/a", noErr},
		"abs/range-overflow-1":  {"/a/b/c/d/e", []string{"1..-6"}, "", hasErr},
		"abs/range-overflow-2":  {"/a/b/c/d/e", []string{"3..2"}, "", hasErr},
		"abs/replace":           {"/a/b/c/d/e", []string{"2", "1", "3"}, "b/a/c", noErr},
		"abs/replace-with-nega": {"/a/b/c/d/e", []string{"2", "1", "-1"}, "b/a/e", noErr},
		"abs/repeat":            {"/a/b/c/d/e", []string{"2", "2", "2"}, "b/b/b", noErr},
		"abs/no-arg":            {"/a/b/c/d/e", []string{}, "", hasErr},
		"abs/full":              {"/a/b/c/d/e", []string{"1..-1"}, "/a/b/c/d/e", noErr},
		"abs/range-left-only-1": {"/a/b/c/d/e", []string{"1.."}, "/a/b/c/d/e", noErr},
		"abs/range-left-only-2": {"/a/b/c/d/e", []string{"3.."}, "c/d/e", noErr},
		"abs/single-positive":   {"/a/b/c/d/e", []string{"3"}, "c", noErr},
		"abs/single-negative":   {"/a/b/c/d/e", []string{"-1"}, "e", noErr},
		// relative path
		"rel/simple":  {"x/y/z", []string{"1", "3"}, "x/z", noErr},
		"rel/range":   {"x/y/z", []string{"2", "3"}, "y/z", noErr},
		"rel/dot":     {"./x/y/z", []string{"1", "3"}, "./x/z", noErr},
		"rel/not-dot": {"./x/y/z", []string{"2", "3"}, "y/z", noErr},
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
