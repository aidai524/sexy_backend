package model

type User struct {
	ID     int64
	Email  string
	Role   string
	Status int
}

type ListenKey struct {
	ListenKey string `json:"listen_key" degate:"listenKey"`
}
