package constants

import "time"

// HTTP related constants
const (
	DefaultTimeout = 10 * time.Second
)

// Context key constants
type ContextKey string

const (
	UserEmailKey ContextKey = "userEmail"
	UserTokenKey ContextKey = "userToken"
)
