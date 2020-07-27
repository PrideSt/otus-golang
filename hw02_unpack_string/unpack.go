package hw02_unpack_string //nolint:golint,stylecheck

import (
	"fmt"

	"github.com/PrideSt/otus-golang/hw02_unpack_string/internal/repeatgroup"
)

// Unpack decode input string.
func Unpack(input string) (string, error) {
	gs, err := repeatgroup.ParseString(input)
	if err != nil {
		return "", fmt.Errorf("unable to parse string %s, %w", input, err)
	}

	return gs.Unpack()
}
