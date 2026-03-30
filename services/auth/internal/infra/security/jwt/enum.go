package jwt

type TokenType string

const (
	Access  TokenType = "access"
	Refresh TokenType = "refresh"
)
