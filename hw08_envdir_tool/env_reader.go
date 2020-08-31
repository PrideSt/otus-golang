package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

type Environment map[string]string

// getValue read the first line from file specified by path and ignore others.
func getValue(path string) string {
	file, err := os.Open(path)

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	defer func() {
		err := file.Close()
		if err != nil {
			log.Fatalf("failed close file: %s", err)
		}
	}()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	scanner.Scan()
	value := scanner.Bytes()

	value = bytes.TrimRight(value, " \t")
	value = bytes.ReplaceAll(value, []byte{'\x00'}, []byte{'\n'})

	return string(value)
}

// envToStrings converts Environment to slice of ket=val strings.
func envToStrings(e Environment) []string {
	result := make([]string, 0, len(e))
	for key, val := range e {
		result = append(result, fmt.Sprintf("%s=%s", key, val))
	}

	return result
}

// stringsToEnv converts key=value strings to Environment.
func stringsToEnv(ss []string) Environment {
	env := make(Environment, len(ss))
	var name, value string
	for _, line := range ss {
		position := strings.IndexByte(line, '=')
		if position == -1 {
			env[line] = ""
			continue
		}

		name = line[:position]
		value = line[position+1:]
		env[name] = value
	}

	return env
}

// mergeEnv lhs and rhs env variables, rhs overwrite lhs ones and unset them if rhs value os empty.
func mergeEnv(lhs Environment, rhs Environment) Environment {
	for key, val := range rhs {
		if len(val) > 0 {
			lhs[key] = val
		} else {
			if _, ok := lhs[key]; !ok {
				continue
			}
			delete(lhs, key)
		}
	}

	return lhs
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := make(Environment)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("unable to scan dir: %w", err)
	}

	for _, fInfo := range files {
		if !fInfo.Mode().IsRegular() {
			continue
		}
		if strings.ContainsAny(fInfo.Name(), " =\n\t") {
			return nil, fmt.Errorf("invalid filename %q", fInfo.Name())
		}

		value := getValue(path.Join(dir, fInfo.Name()))

		if len(value) > 0 {
			env[fInfo.Name()] = value
		}
	}
	return env, nil
}
