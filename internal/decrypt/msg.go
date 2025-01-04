package decrypt

import (
	"fmt"
	"gokdb/internal/model"
	"strings"
)

func DecryptMessage(msg *model.Message, enc uint32) error {
	var err error
	if msg.Attachment != "" && msg.Attachment != "{}" {
		msg.Attachment, err = Decrypt(
			msg.UserId, msg.Attachment, enc,
		)
		if err != nil {
			return fmt.Errorf("failed to decrypt attachment: %v", err)
		}
	}

	if strings.Trim(msg.Message, "=") != "" {
		msg.Message, err = Decrypt(
			msg.UserId, msg.Message, enc,
		)
		if err != nil {
			return fmt.Errorf("failed to decrypt message: %v", err)
		}
	}

	return nil
}
