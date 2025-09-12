package service

import (
	"time"
	"ecommerce/helper"
)

func ParseAccessForMiddleware(token string) (*helper.JWTClaims, error) {
	return helper.ParseAccess(token)
}

func AccessBlacklistLookup(token string) (time.Time, bool) {
	exp, ok := accessBlacklist[token]
	return exp, ok
}
