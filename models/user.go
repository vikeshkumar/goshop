package models

import "time"

type User struct {
	Uid          string
	Email        string
	Locked       bool
	Hash         string
	Password     string
	LocalUser    bool
	SsoUser      bool
	SsoProvider  *SSOProvider
	CreationTime time.Time
	ModifiedTime time.Time
	Roles        []Role
}
