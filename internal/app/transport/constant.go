package transport

import (
	"fmt"
	"time"
)

const (
	TagPassword  = "password"
	TagEmail     = "email"
	TagNickName  = "nickname"
	TagPhoneNum  = "phoneNum"
	TagTermAgree = "termAgree"
	TagIsseudAt  = "isseudAt"
)

type TimeString time.Time

func (t TimeString) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02T15:04:05+09:00"))
	return []byte(stamp), nil
}
