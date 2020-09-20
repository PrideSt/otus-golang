package parser

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

func Test_scanValidators(t *testing.T) {
	type Args struct {
		s string
	}
	tests := []struct {
		name string
		Args Args
		want []ValidatorDesc
	}{
		{
			name: "empty",
			Args: Args{
				s: ``,
			},
			want: nil,
		},
		{
			name: "without Args",
			Args: Args{
				s: `empty`,
			},
			want: []ValidatorDesc{
				{
					FuncName: "empty",
					Args:     nil,
				},
			},
		},
		{
			name: "simple",
			Args: Args{
				s: `min:18`,
			},
			want: []ValidatorDesc{
				{
					FuncName: "min",
					Args:     []string{"18"},
				},
			},
		},
		{
			name: "two Args",
			Args: Args{
				s: `between:18,25`,
			},
			want: []ValidatorDesc{
				{
					FuncName: "between",
					Args:     []string{"18", "25"},
				},
			},
		},
		{
			name: "several Args",
			Args: Args{
				s: `in:18,25,36,45`,
			},
			want: []ValidatorDesc{
				{
					FuncName: "in",
					Args:     []string{"18", "25", "36", "45"},
				},
			},
		},
		{
			name: "several validators",
			Args: Args{
				s: `in:18,25|min:34`,
			},
			want: []ValidatorDesc{
				{
					FuncName: "in",
					Args:     []string{"18", "25"},
				},
				{
					FuncName: "min",
					Args:     []string{"34"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := scanValidators(tt.Args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("scanValidators() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseTags(t *testing.T) {
	type Args struct {
		s string
	}
	tests := []struct {
		name string
		Args Args
		want map[string]string
	}{
		{
			name: "empty",
			Args: Args{
				s: ``,
			},
			want: nil,
		},
		{
			name: "without package",
			Args: Args{
				s: `strange_tag_without_package`,
			},
			want: nil,
		},
		{
			name: "simple",
			Args: Args{
				s: `xml:"foo"`,
			},
			want: map[string]string{
				"xml": "foo",
			},
		},
		{
			name: "several values",
			Args: Args{
				s: `xml:"foo,bar,zxc"`,
			},
			want: map[string]string{
				"xml": "foo,bar,zxc",
			},
		},
		{
			name: "without value",
			Args: Args{
				s: `xml:""`,
			},
			want: nil,
		},
		{
			name: "several packages",
			Args: Args{
				s: `json:"foo,omitempty,string" xml:"foo"`,
			},
			want: map[string]string{
				"json": "foo,omitempty,string",
				"xml":  "foo",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseTags(tt.Args.s)
			assertMap(t, tt.want, got)
		})
	}
}

func joinMap(m map[string]string) []string {
	var ss []string
	for k, v := range m {
		ss = append(ss, fmt.Sprintf("%s=%v", k, v))
	}
	sort.Strings(ss)
	return ss
}

func assertMap(t *testing.T, want, got map[string]string) {
	w := joinMap(want)
	g := joinMap(got)

	if !reflect.DeepEqual(g, w) {
		t.Errorf("scanValidators() = %v, want %v", g, w)
	}
}
