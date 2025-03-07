package models

import "time"


type Emails struct{
	Username []string `json:"username"`
	GroupName string `json:"groupname"`
	CreatedAt time.Time
}