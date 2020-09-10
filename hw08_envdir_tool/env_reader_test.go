package main

import (
	"reflect"
	"sort"
	"testing"
)

func TestStringsToEnv(t *testing.T) {
	type args struct {
		ss []string
	}
	tests := []struct {
		name string
		args args
		want Environment
	}{
		{
			name: `simple`,
			args: args{
				ss: []string{"aaa=1", "bbb=2", "ccc=3"},
			},
			want: Environment{"aaa": "1", "bbb": "2", "ccc": "3"},
		},
		{
			name: `without value`,
			args: args{
				ss: []string{"aaa="},
			},
			want: Environment{"aaa": ""},
		},
		{
			name: `empty`,
			args: args{
				ss: []string{},
			},
			want: Environment{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stringsToEnv(tt.args.ss)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("stringsToEnv() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnvToStrings(t *testing.T) {
	type args struct {
		e Environment
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: `simple`,
			args: args{
				e: Environment{"aaa": "1", "bbb": "2", "ccc": "3"},
			},
			want: []string{"aaa=1", "bbb=2", "ccc=3"},
		},
		{
			name: `empty`,
			args: args{
				e: Environment{},
			},
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := envToStrings(tt.args.e)
			sort.Strings(got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("envToStrings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getValue(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: `simple`,
			args: args{
				path: `./testdata/simple/var1`,
			},
			want: `fooooo`,
		},
		{
			name: `value with space`,
			args: args{
				path: `./testdata/escape-value/var1`,
			},
			want: `value with space`,
		},
		{
			name: `value with trailing space`,
			args: args{
				path: `./testdata/escape-value/var2`,
			},
			want: `value_with_trailing_space`,
		},
		{
			name: `value with trailing tab`,
			args: args{
				path: `./testdata/escape-value/var3`,
			},
			want: `value_with_trailing_tab`,
		},
		{
			name: `value with trailing new line`,
			args: args{
				path: `./testdata/escape-value/var4`,
			},
			want: `value_with_trailing_new_line`,
		},
		{
			name: `multy lines in file`,
			args: args{
				path: `./testdata/escape-value/var5`,
			},
			want: `multyline_one`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getValue(tt.args.path); got != tt.want {
				t.Errorf("getValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadDir(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		want    Environment
		wantErr bool
	}{
		{
			name: `simple`,
			args: args{
				dir: `./testdata/simple`,
			},
			want: Environment{
				"var1": "fooooo",
				"var2": "barrr",
			},
		},
		{
			name: `escape value`,
			args: args{
				dir: `./testdata/escape-value`,
			},
			want: Environment{
				"var1": "value with space",
				"var2": "value_with_trailing_space",
				"var3": "value_with_trailing_tab",
				"var4": "value_with_trailing_new_line",
				"var5": "multyline_one",
			},
		},
		{
			name: `nested dirs`,
			args: args{
				dir: `./testdata/nested-dirs`,
			},
			want: Environment{
				"var1": "fooooo",
			},
		},
		{
			name: `invalid filename`,
			args: args{
				dir: `./testdata/invalid-filename`,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadDir(tt.args.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadDir() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMergeEnv(t *testing.T) {
	type args struct {
		lhs Environment
		rhs Environment
	}
	tests := []struct {
		name string
		args args
		want Environment
	}{
		{
			name: `simple`,
			args: args{
				lhs: Environment{"aaa": "1", "bbb": "2"},
				rhs: Environment{"bbb": "3", "ccc": "4"},
			},
			want: Environment{"aaa": "1", "bbb": "3", "ccc": "4"},
		},
		{
			name: `lhs empty`,
			args: args{
				lhs: Environment{},
				rhs: Environment{"bbb": "3", "ccc": "4"},
			},
			want: Environment{"bbb": "3", "ccc": "4"},
		},
		{
			name: `rhs empty`,
			args: args{
				lhs: Environment{"aaa": "1", "bbb": "2"},
				rhs: Environment{},
			},
			want: Environment{"aaa": "1", "bbb": "2"},
		},
		{
			name: `both empty`,
			args: args{
				lhs: Environment{},
				rhs: Environment{},
			},
			want: Environment{},
		},
		{
			name: `unset bbb`,
			args: args{
				lhs: Environment{"aaa": "1", "bbb": "2"},
				rhs: Environment{"bbb": "", "ccc": "4"},
			},
			want: Environment{"aaa": "1", "ccc": "4"},
		},
		{
			name: `unset bbb not exists`,
			args: args{
				lhs: Environment{"aaa": "1"},
				rhs: Environment{"bbb": "", "ccc": "4"},
			},
			want: Environment{"aaa": "1", "ccc": "4"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mergeEnv(tt.args.lhs, tt.args.rhs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mergeEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
