package models

import (
	"net/http"
	"time"
)

type Middleware func(handler http.HandlerFunc) http.HandlerFunc

type Session struct {
	UserID         int
	SessionID      string
	Username       string
	IpAddress      string
	ExpirationTime time.Time
}
