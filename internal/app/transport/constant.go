package transport

import (
	"database/sql/driver"
	"errors"
	"fmt"
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

func (t TimeString) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02T15:04:05"))
	return []byte(stamp), nil
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
