// +build bench

package hw10_program_optimization //nolint:golint,stylecheck

import (
	"archive/zip"
	"github.com/stretchr/testify/require"
	"testing"
)

//go test -v -run=None -bench=. -tags=bench -benchmem -count=20 -benchtime=10s ./... | tee bench-X.mem
//
// compare with base version
//name    old time/op    new time/op    delta
//Stat-8    1.47ms ± 7%    0.74ms ± 0%    -49.39%  (p=0.000 n=18+12)
//
//name    old alloc/op   new alloc/op   delta
//Stat-8    41.0kB ±14%  3208.5kB ± 0%  +7717.12%  (p=0.000 n=17+14)
//
//name    old allocs/op  new allocs/op  delta
//Stat-8       435 ±14%        15 ± 0%    -96.55%  (p=0.000 n=17+14)
func BenchmarkStat(b *testing.B) {
	r, err := zip.OpenReader("testdata/users.dat.zip")
	defer r.Close()

	data, err := r.File[0].Open()
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetDomainStat(data, "biz")
	}
}
