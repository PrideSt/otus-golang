package main

import (
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

var (
	ntpPool = []string{
		"0.ru.pool.ntp.org",
	}
)

const (
	// ErrorCodeNTPPoolInvalidConfig happens when NTP pool not configurated
	ErrorCodeNTPPoolInvalidConfig = 1 << iota
	// ErrorCodeNTPLookup happens when NTP service unavalible
	ErrorCodeNTPLookup
)

func main() {
	// choose some NTP host
	ntpHost, err := getNTPHost(0)
	if err != nil {
		err := fmt.Errorf("invalid NTP pool configuration, %w", err)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(ErrorCodeNTPPoolInvalidConfig)
	}

	// make request to NTP server
	ntpTime, err := ntp.Time(ntpHost)
	if err != nil {
		err := fmt.Errorf("unable get NTP time, %w", err)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(ErrorCodeNTPLookup)
	}

	localTime := time.Now()

	// out response to stdout
	fmt.Printf("current time: %s\nexact time: %s\n", localTime.Round(0).UTC(), ntpTime.Round(0).UTC())
}

func getNTPHost(offset int) (string, error) {
	if len(ntpPool) == 0 {
		return "", fmt.Errorf("NTP pool is empty")
	}

	if len(ntpPool) < offset {
		return "", fmt.Errorf("NTP host offest is out of range [offset/pool size] = [%d/%d]", offset, len(ntpPool))
	}

	return ntpPool[offset], nil
}
