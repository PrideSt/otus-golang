// +build !bench

package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDomainStat(t *testing.T) {
	data := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

	t.Run("find 'com'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("find 'gov'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
	})

	t.Run("find 'unknown'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "unknown")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})
}

func Test_countDomains(t *testing.T) {
	type args struct {
		users  []User
		domain string
	}
	tests := []struct {
		name    string
		args    args
		want    DomainStat
		wantErr bool
	}{
		{
			name: "empty",
			args: args{
				domain: "abc",
			},
			want: DomainStat{},
		},
		{
			name: "no result",
			args: args{
				users: []User{
					{
						Name:   "vasya@gmail.com",
						Domain: "com",
					},
				},
				domain: "abc",
			},
			want: DomainStat{},
		},
		{
			name: "simple",
			args: args{
				users: []User{
					{
						Name:   "vasya@gmail.com",
						Domain: "com",
					},
				},
				domain: "com",
			},
			want: DomainStat{"vasya@gmail.com": 1},
		},

		{
			name: "several samples",
			args: args{
				users: []User{
					{
						Name:   "vasya@gmail.com",
						Domain: "com",
					},
					{
						Name:   "masha@gmail.com",
						Domain: "com",
					},
					{
						Name:   "archibald@mail.ru",
						Domain: "ru",
					},
					{
						Name:   "vasya@gmail.com",
						Domain: "com",
					},
				},
				domain: "com",
			},
			want: DomainStat{
				"vasya@gmail.com": 2,
				"masha@gmail.com": 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := countDomains(tt.args.users, tt.args.domain)
			if (err != nil) != tt.wantErr {
				t.Errorf("countDomains() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.Equal(t, mapToSliceString(got), mapToSliceString(tt.want))
		})
	}
}

func mapToSliceString(m DomainStat) []string {
	result := make([]string, 0, len(m))
	for name, cnt := range m {
		result = append(result, fmt.Sprintf("%s=%d", name, cnt))
	}

	sort.Strings(result)
	return result
}

func Test_getUsers(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []User
		wantErr bool
	}{
		{
			name: "empty",
			args: args{
				r: strings.NewReader(""),
			},
			want: []User{},
		},
		{
			name: "simple",
			args: args{
				r: strings.NewReader(`{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Brian Olson","Username":"non_quia_id","Email":"FrancesEllis@Quinu.edu","Phone":"237-75-34","Password":"cmEPhX8","Address":"Butterfield Junction 74"}
{"Id":3,"Name":"Justin Oliver Jr. Sr.","Username":"oPerez","Email":"MelissaGutierrez@Twinte.gov","Phone":"106-05-18","Password":"f00GKr9i","Address":"Oak Valley Lane 19"}`),
			},
			want: []User{
				{
					Name:   "browsedrive.gov",
					Domain: "gov",
				},
				{
					Name:   "quinu.edu",
					Domain: "edu",
				},
				{
					Name:   "twinte.gov",
					Domain: "gov",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getUsers(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("getUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, len(tt.want), len(got), "getUsers() len(got) = %d, want %d", len(got), len(tt.want))
			require.Equal(t, tt.want, got)
		})
	}
}
