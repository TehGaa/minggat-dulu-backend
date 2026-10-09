package model

import "time"

type Token struct {
	UserId       string    `json:"user_id"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	AccessUUID   string    `json:"access_uuid"`
	RefreshUUID  string    `json:"refresh_uuid"`
	AtExpires    time.Time `json:"at_expires"`
	RtExpires    time.Time `json:"rt_expires"`
}
