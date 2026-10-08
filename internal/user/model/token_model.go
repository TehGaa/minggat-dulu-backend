package model

import "time"

type Token struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	AtExpires    time.Time `json:"at_expires"`
	RtExpires    time.Time `json:"rt_expires"`
}
