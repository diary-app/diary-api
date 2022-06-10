package common

import (
	"fmt"
	"strings"
	"time"
)

type DateOnly time.Time

const dateLayout = "2006-01-02"

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse(dateLayout, s)
	if err != nil {
		return err
	}
	*d = DateOnly(t)
	return nil
}

func (d DateOnly) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(d).Format(dateLayout))
	return []byte(stamp), nil
}

func (d DateOnly) Format(s string) string {
	t := time.Time(d)
	return t.Format(s)
}
