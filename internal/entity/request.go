package entity

import "github.com/GoExpertCurso/GoRateLimiterFC/configs"

type Request struct {
	Ip    *Ip
	Token *Token
	Limit int64
}

func NewRequest(ip string, token string, limits configs.Conf) *Request {
	ipAdress := NewIp(ip, int64(limits.IPLimit))
	accessToken := NewToken(token, int64(limits.TokenLimit))
	return &Request{
		Ip:    ipAdress,
		Token: accessToken,
	}
}

func (r *Request) LimitCheck() {
	if r.Token.limit == 0 {
		r.Limit = r.Ip.limit
	} else {
		r.Limit = r.Token.limit
	}
}
