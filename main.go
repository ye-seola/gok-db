package main

import (
	"encoding/json"
	"gokdb/internal/kdb"
	"gokdb/internal/model"
	"gokdb/internal/sendmsg"
	"gokdb/internal/utils"
	"gokdb/internal/ws"
	"log"
	"time"
)

func onSendMsg(event string, payload json.RawMessage) error {
	var sendMsg model.SendMsgPayload
	if err := json.Unmarshal(payload, &sendMsg); err != nil {
		return err
	}

	if err := sendmsg.SendMsg(sendMsg.ChatId, sendMsg.Message); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func main() {
	db, err := kdb.New()
	if err != nil {
		panic(err)
	}

	ws.SetHandler("SENDMSG", onSendMsg)
	go ws.Start(":9023")

	// init
	lastDB1ModTime, err := utils.GetDB1ModifiedTime()
	if err != nil {
		panic(err)
	}

	lastProcessedRowId, lastProcessedMsgId, err := db.GetLastId()
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(100 * time.Millisecond)

	for range ticker.C {
		currDB1ModTime, err := utils.GetDB1ModifiedTime()
		if err != nil {
			log.Println(err)
			continue
		}

		// DB1 Changed
		if currDB1ModTime != lastDB1ModTime {
			lastDB1ModTime = currDB1ModTime

			messages, err := db.GetMessagesAfterRowId(lastProcessedRowId)
			if err != nil {
				log.Println(err)
				continue
			}

			for _, msg := range messages {
				if msg.MsgId <= lastProcessedMsgId {
					continue
				}

				ws.BroadcastEvent("MSG", msg)
				lastProcessedRowId = msg.RowId
				lastProcessedMsgId = msg.MsgId
			}
		}
	}
}
