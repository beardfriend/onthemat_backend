package transport

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	TagPassword   = "password"
	TagEmail      = "email"
	TagNickName   = "nickname"
	TagLogoUrl    = "logoUrl"
	TagPhoneNum   = "phoneNum"
	TagCallNumber = "callNumber"
	TagTermAgree  = "termAgree"
	TagIsseudAt   = "isseudAt"
)

type TimeString time.Time

func (t *TimeString) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(*t).Format("2006-01-02T15:04:05"))
	return []byte(stamp), nil
}

func (t *TimeString) ToString() string {
	return time.Time(*t).Format("2006-01-02T15:04:05")
}

func (t *TimeString) String() string {
	return time.Time(*t).String()
}

func (t *TimeString) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), "\"")
	parsed, _ := time.Parse("2006-01-02T15:04:05", str)
	*t = TimeString(parsed)
	return nil
}

func (t TimeString) Now() TimeString {
	return TimeString(time.Now())
}

func (t TimeString) Value() (driver.Value, error) {
	return time.Time(t), nil
}

func (t *TimeString) Scan(value interface{}) error {
	if value == nil {
		*t = TimeString(time.Now())
		return nil
	}
	if v, ok := value.(time.Time); ok {
		*t = TimeString(v)
		return nil
	}

	return errors.New("cannot scan")
}
