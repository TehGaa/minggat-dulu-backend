package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Token struct {
	UserId       string    `json:"user_id"`
	UserData     bson.M    `json:"user_data"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	AccessUUID   string    `json:"access_uuid"`
	RefreshUUID  string    `json:"refresh_uuid"`
	AtExpires    time.Time `json:"at_expires"`
	RtExpires    time.Time `json:"rt_expires"`
}
