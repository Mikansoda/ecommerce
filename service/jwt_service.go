package service

import "time"

func ParseAccessForMiddleware(token string) (*JWTClaims, error) {
	return parseAccess(token)
}

func AccessBlacklistLookup(token string) (time.Time, bool) {
	exp, ok := accessBlacklist[token]
	return exp, ok
}
