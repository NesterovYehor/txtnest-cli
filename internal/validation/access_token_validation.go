package validation

import (
	"time"
)

func ValidateAccessToken(expires time.Time) bool {
	return expires.After(time.Now())
}

func ValidateRefreshToken(expiresAt time.Time) bool {
	return expiresAt.After(time.Now())
}
