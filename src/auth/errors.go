package auth

import "errors"

var (
	ErrFailedToParseClaims = errors.New("failed to parse claims")
	ErrTokenExpired        = errors.New("token expired")
)
