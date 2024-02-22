package entity

type Ip struct {
	Address string
	limit   int64
}

func NewIp(address string, limit int64) *Ip {
	return &Ip{
		Address: address,
		limit:   limit,
	}
}
