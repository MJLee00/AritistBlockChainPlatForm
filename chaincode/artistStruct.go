package main

type Artist struct{
	CertNo	string	`json:"CertNo"`	// 证书编号 key
	Photo	string	`json:"Photo"`	// 照片
	Params string `json:"Params"`   //参数
	Type string `json:"Type"`//类型
    Date string `json:"Date"`

	Historys	[]HistoryItem	// 当前art的历史记录
}

type HistoryItem struct {
	TxId	string
	Artist	Artist
}
