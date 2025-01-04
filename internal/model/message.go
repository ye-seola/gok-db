package model

type Message struct {
	UserId     uint64 `json:"userId"`
	MsgType    int    `json:"msgType"`
	Message    string `json:"message"`
	Attachment string `json:"attachment"`
	RowId      int64  `json:"rowId"`
	MsgId      int64  `json:"msgId"`
	ChatId     int64  `json:"chatId"`
}
