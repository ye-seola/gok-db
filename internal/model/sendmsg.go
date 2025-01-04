package model

type SendMsgPayload struct {
	ChatId  int64  `json:"chatId"`
	Message string `json:"message"`
}
