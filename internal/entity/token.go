package entity

type Token struct {
	AccessToken string
	limit       int64
}

func NewToken(accessToken string, limit int64) *Token {
	return &Token{
		AccessToken: accessToken,
		limit:       limit,
	}
}
