package myntp

import (
	"time"

	"github.com/beevik/ntp"
)

const DefaultNTPServer = "time.google.com:123"

func GetTime(address string) (time.Time, error) {
	t, err := ntp.Time(address)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
