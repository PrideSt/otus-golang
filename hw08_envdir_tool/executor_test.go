package main

import (
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func unsetAllEnvVars(whitelist map[string]struct{}) {
	for key := range stringsToEnv(os.Environ()) {
		_, ok := whitelist[key]
		if !ok {
			err := os.Unsetenv(key)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func resetEnv(e Environment) {
	unsetAllEnvVars(nil)

	for key, val := range e {
		os.Setenv(key, val)
	}
}

func TestRunCmd(t *testing.T) {
	envUtilArgs := []string{"/usr/bin/env"}
	type args struct {
		args  []string
		env   Environment
		osEnv Environment
	}
	tests := []struct {
		name           string
		args           args
		want           []string
		wantReturnCode int
	}{
		{
			name: `run simple`,
			args: args{
				args: envUtilArgs,
				env:  Environment{"aaa": "1", "bbb": "2"},
			},
			want: []string{"", "aaa=1", "bbb=2"},
		},
		{
			name: `run env overwrite`,
			args: args{
				args:  envUtilArgs,
				env:   Environment{"aaa": "1", "bbb": "2"},
				osEnv: Environment{"aaa": "3"},
			},
			want: []string{"", "aaa=1", "bbb=2"},
		},
		{
			// unset obtain in main merge
			name: `run env not unset`,
			args: args{
				args:  envUtilArgs,
				env:   Environment{"aaa": ""},
				osEnv: Environment{"aaa": "1", "bbb": "2"},
			},
			want: []string{"", "aaa="},
		},
		{
			name: `run cmd not found`,
			args: args{
				args: []string{"./cmd/env/file-not-found"},
			},
			wantReturnCode: 1,
			want:           []string{""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetEnv(tt.args.osEnv)

			stdOut := os.Stdout
			buffer, err := ioutil.TempFile("/tmp", tt.name)
			if err != nil {
				log.Fatal(err)
			}

			os.Stdout = buffer
			gotReturnCode := RunCmd(tt.args.args, tt.args.env)
			os.Stdout = stdOut

			err = buffer.Close()
			if err != nil {
				log.Fatal(err)
			}

			if gotReturnCode != tt.wantReturnCode {
				t.Errorf("RunCmd() = %v, want %v", gotReturnCode, tt.wantReturnCode)
			}

			output, _ := ioutil.ReadFile(buffer.Name())
			lines := strings.Split(string(output), "\n")
			sort.Strings(lines)
			require.Equal(t, tt.want, lines)
		})
	}
}
