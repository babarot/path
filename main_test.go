package main

import (
	"fmt"
	"testing"
)

type MockCLI struct {
	Out string
}

func NewMockRenderer() *MockCLI {
	return &MockCLI{}
}

func (c *MockCLI) Printf(format string, a ...any) {
	c.Out += fmt.Sprintf(format, a...)
}

func Test_run_(t *testing.T) {
	cases := []struct {
		name  string
		input string
		args  []string
		want  map[string]string
	}{
		{
			name:  "kubernetes: regular case (min match #3)",
			input: "",
			want: map[string]string{
				"": "",
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			renderer := NewMockRenderer()
			err := run(renderer, tt.args)
			if err != nil {
				t.Errorf(err.Error())
			}
		})
	}
}
