package cache

import "fmt"

func OTPKey(userID string, purpose string) string {
	return fmt.Sprintf("erp:otp:%s:%s", purpose, userID)
}

func LoginRateLimitKey(identifier string) string {
	return fmt.Sprintf("erp:rate-limit:login:%s", identifier)
}

func PasswordResetKey(userID string) string {
	return fmt.Sprintf("erp:password-reset:%s", userID)
}

func SessionKey(sessionID string) string {
	return fmt.Sprintf("erp:session:%s", sessionID)
}

func UserCacheKey(userID string) string {
	return fmt.Sprintf("erp:user:%s", userID)
}
