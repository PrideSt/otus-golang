package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	// {0: go-bin, 1: end-dir-path, 3: subcommand, 4...: subcommand-args}
	minArgs := 3

	if len(os.Args) < minArgs {
		log.Fatalf(
			"invalid argument: expected env-dir, subcommand, given: %s",
			strings.Join(os.Args[1:], ", "),
		)
	}
	env, err := ReadDir(os.Args[1])
	if err != nil {
		log.Fatalf("reading dir failed: %s", err)
	}

	osEnv := stringsToEnv(os.Environ())
	if len(env) > 0 {
		env = mergeEnv(osEnv, env)
	} else {
		env = osEnv
	}

	exitCode := RunCmd(os.Args[2:], env)
	if exitCode > 0 {
		fmt.Fprintln(os.Stderr, "subcommand failed")
		os.Exit(exitCode)
	}
}
