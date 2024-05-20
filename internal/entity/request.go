package entity

type Request struct {
	Key   string
	Limit int
}

func NewRequest(key string, limit int) *Request {
	return &Request{
		Key:   key,
		Limit: limit,
	}
}
