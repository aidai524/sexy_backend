package model

type AccountData struct {
	SuperLikeTotal     int `json:"total_super_like"`     // 总super like数量
	SuperLikeAvailable int `json:"available_super_like"` // 可用的super like数量
}
