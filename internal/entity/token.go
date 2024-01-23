package entity

type Token struct {
	AccessToken string
	limit       int64
}

func NewToken(accessToken string) *Token {
	return &Token{
		AccessToken: accessToken,
		limit:       10,
	}
}
