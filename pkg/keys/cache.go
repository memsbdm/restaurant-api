package keys

import "time"

// Signed Purpose Token
type SPT string

// Opaque Access Token
type OAT string

const (
	AuthToken         OAT = "access_token"
	EmailVerification SPT = "email_verification"
)

var (
	AuthTokenDuration              = 7 * 24 * time.Hour
	EmailVerificationTokenDuration = 24 * time.Hour
)
