package sendmsg

import (
	"encoding/json"
	"fmt"
	"log"
)

var (
	notificationReferer = mustFetchNotificationReferer()
)

func SendMsg(chatId int64, message string) error {
	cmdMutex.Lock()
	defer cmdMutex.Unlock()

	if !started {
		mustStart()
	}

	marshaled, err := json.Marshal(map[string]interface{}{
		"msg":     message,
		"chatId":  chatId,
		"notiRef": notificationReferer,
	})
	if err != nil {
		return err
	}

	_, err = stdin.Write(append(marshaled, 0xa))
	if err != nil {
		return err
	}

	if stdout.Scan() {
		v := map[string]interface{}{}
		err = json.Unmarshal(stdout.Bytes(), &v)
		if err != nil {
			return err
		}

		fmt.Printf("%+v\n", v)

		if v["success"].(bool) {
			return nil
		}

		return fmt.Errorf("%s", v["error"].(string))
	}

	if err := stdout.Err(); err != nil {
		return err
	}

	log.Fatalln("SendMsg EOF")
	return nil
}
