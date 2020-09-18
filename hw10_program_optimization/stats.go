package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type User struct {
	Name   string
	Domain string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %s", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) ([]User, error) { //nolint:unparam
	var result users
	i := 0
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		name := getName(scanner.Text())
		result[i] = User{
			Name:   name,
			Domain: name[strings.LastIndex(name, ".")+1:],
		}
		i++
	}

	return result[:i], nil
}

var emailBeginPattern = `Email":"`

func getName(line string) string {
	// without regexp much faster
	line = line[strings.Index(line, emailBeginPattern)+len(emailBeginPattern)+1:]
	line = line[:strings.Index(line, `"`)]
	line = line[strings.Index(line, "@")+1:]

	return strings.ToLower(line)
}

func countDomains(users []User, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range users {
		if user.Domain == domain {
			result[user.Name]++
		}
	}

	return result, nil
}
