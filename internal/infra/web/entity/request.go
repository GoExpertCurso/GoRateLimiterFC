package entity

type Request struct {
	Ip    *Ip
	Token *Token
	Limit int64
}

func NewRequest(ip string, token string) *Request {
	ipAdress := NewIp(ip)
	accessToken := NewToken(token)
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
