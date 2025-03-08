package models

import "time"


type Emails struct{
	Username []string `json:"username" bson:"username"`
	GroupName string `json:"groupname" bson:"groupname"`
	CreatedAt time.Time
}