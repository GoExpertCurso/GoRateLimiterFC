package entity

type Ip struct {
	Address string
	limit   int64
}

func NewIp(address string) *Ip {
	return &Ip{
		Address: address,
		limit:   10,
	}
}
