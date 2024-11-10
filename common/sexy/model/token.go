package model

type TokenInfo struct {
	Id       int    `json:"id"`
	Address  string `json:"code"`
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Icon     string `json:"icon"`
	Decimals int32  `json:"decimals"`
	OracleId string `json:"oracle_id"`
}
