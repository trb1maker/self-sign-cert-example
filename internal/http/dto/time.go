package dto

import "time"

func NewTimeResponse() TimeResponse {
	return TimeResponse{TS: time.Now().UTC().Unix()}
}

type TimeResponse struct {
	TS int64 `json:"ts"`
}

func (t *TimeResponse) Time() time.Time {
	return time.Unix(t.TS, 0).Local() //nolint:gosmopolitan
}
