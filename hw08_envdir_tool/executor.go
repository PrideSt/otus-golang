package main

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	ExitNormal = iota
	ExitError
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(args []string, env Environment) (returnCode int) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "not enough arguments, pass subcommand")
		return ExitError
	}

	cmd := exec.Command(args[0], args[1:]...) //nolint:gosec
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = envToStrings(env)

	err := cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ExitError
	}

	return ExitNormal
}
