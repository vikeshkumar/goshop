package models

import "time"

type Role struct {
	Code         string
	Description  string
	Expired      bool
	Locked       bool
	CreationTime time.Time
	ModifiedTime time.Time
}
