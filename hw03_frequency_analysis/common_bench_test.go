package hw03_frequency_analysis //nolint:golint,stylecheck

//go:generate go run ./cmd/cwd.go
//go:generate go run ./cmd/txtgen/txtgen.go -text-size=1000 -dict-size=100 -pkg-name=generated -out ./internal/generated/txt1KDict100.go
//go:generate go run ./cmd/txtgen/txtgen.go -text-size=10000 -dict-size=100 -pkg-name=generated -out ./internal/generated/txt10KDict100.go
//go:generate go run ./cmd/txtgen/txtgen.go -text-size=100000 -dict-size=100 -pkg-name=generated -out ./internal/generated/txt100KDict100.go
//go:generate go run ./cmd/txtgen/txtgen.go -text-size=1000000 -dict-size=100 -pkg-name=generated -out ./internal/generated/txt1MDict100.go

//go:generate go run ./cmd/txtgen/txtgen.go -text-size=100000 -dict-size=1000 -pkg-name=generated -out ./internal/generated/txt100KDict1K.go
//go:generate go run ./cmd/txtgen/txtgen.go -text-size=100000 -dict-size=10000 -pkg-name=generated -out ./internal/generated/txt100KDict10K.go

import (
	"fmt"
	"testing"

	"github.com/PrideSt/otus-golang/hw03_frequency_analysis/internal/generated"
	"github.com/PrideSt/otus-golang/hw03_frequency_analysis/internal/pool"
)

// go test -v -run=None -bench=^BenchmarkTopN$ -gcflags=-l -benchtime=30s ./...
// goos: darwin
// goarch: amd64
// pkg: github.com/PrideSt/otus-golang/hw03_frequency_analysis
// BenchmarkTopN/1_000Dict100_single-8                61569            584720 ns/op
// BenchmarkTopN/1_000Dict100_multy-8                  6885           5368810 ns/op
// BenchmarkTopN/10_000Dict100_single-8                6159           5688278 ns/op
// BenchmarkTopN/10_000Dict100_multy-8                  693          51902209 ns/op
// BenchmarkTopN/100_000Dict100_single-8                607          59495543 ns/op
// BenchmarkTopN/100_000Dict100_multy-8                  67         522850980 ns/op
// BenchmarkTopN/1_000_000Dict100_single-8               60         569253823 ns/op
// BenchmarkTopN/1_000_000Dict100_multy-8                 6        5084911734 ns/op
// BenchmarkTopN/100_000Dict1_000_single-8              601          59816168 ns/op
// BenchmarkTopN/100_000Dict1_000_multy-8                68         524167155 ns/op
// BenchmarkTopN/100_000Dict10_000_single-8             560          64313097 ns/op
// BenchmarkTopN/100_000Dict10_000_multy-8               66         524178160 ns/op
func BenchmarkTopN(b *testing.B) {
	for _, bt := range [...]struct {
		name string
		text string
	}{
		{
			name: `1_000Dict100`,
			text: generated.Txt1000Dict100,
		},
		{
			name: `10_000Dict100`,
			text: generated.Txt10000Dict100,
		},
		{
			name: `100_000Dict100`,
			text: generated.Txt100000Dict100,
		},
		{
			name: `1_000_000Dict100`,
			text: generated.Txt1000000Dict100,
		},
		{
			name: `100_000Dict1_000`,
			text: generated.Txt100000Dict1000,
		},
		{
			name: `100_000Dict10_000`,
			text: generated.Txt100000Dict10000,
		},
	} {
		b.Run(fmt.Sprintf("%s %s", bt.name, "single"), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Top10(bt.text, nil)
			}
		})

		b.Run(fmt.Sprintf("%s %s", bt.name, "multy"), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				pool.Top10(bt.text, nil)
			}
		})
	}
}
