package utils

import (
	"database/sql/driver"
	"strconv"
	"time"
)

type FormatTime time.Time

func (t FormatTime) MarshalJSON() ([]byte, error) {
	var timeStr string
	if !time.Time(t).IsZero() {
		timeStr = time.Time(t).Format("2006-01-02 15:04:05")
	}
	return []byte(strconv.Quote(timeStr)), nil
}
func (t FormatTime) Value() (driver.Value, error) {
	if time.Time(t).IsZero() {
		return nil, nil
	}
	return time.Time(t), nil
}
func (t *FormatTime) UnmarshalJSON(stringTime []byte) error {
	if t == nil {
		return nil
	}
	v, _ := strconv.Unquote(string(stringTime))
	tt, _ := time.ParseInLocation("2006-01-02 15:04:05", v, time.Local)
	*t = FormatTime(tt)
	return nil
}
