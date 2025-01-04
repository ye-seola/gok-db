package kdb

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"gokdb/internal/constsnts"
	"gokdb/internal/decrypt"
	"gokdb/internal/model"

	"github.com/mattn/go-sqlite3"
)

type KDB struct {
	db *sql.DB
}

func New() (*KDB, error) {
	sql.Register("sqlite3_kdb",
		&sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				if _, err := conn.Exec("ATTACH DATABASE ? AS db1;", []driver.Value{
					constsnts.DB1Path,
				}); err != nil {
					return err
				}

				if _, err := conn.Exec("ATTACH DATABASE ? AS db2;", []driver.Value{
					constsnts.DB2Path,
				}); err != nil {
					return err
				}

				return nil
			},
		})

	db, err := sql.Open("sqlite3_kdb", ":memory:")
	if err != nil {
		return nil, err
	}

	return &KDB{
		db: db,
	}, nil
}

// rowId, msgId
func (k *KDB) GetLastId() (int64, int64, error) {
	rows, err := k.db.Query("select _id, id from db1.chat_logs order by _id desc limit 1")
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()

	var rowId, msgId int64
	if !rows.Next() {
		return 0, 0, fmt.Errorf("failed to get last chat_logs _id, id")
	}

	err = rows.Scan(&rowId, &msgId)
	if err != nil {
		return 0, 0, err
	}

	return rowId, msgId, nil
}

func (m *KDB) GetMessagesAfterRowId(rowId int64) ([]model.Message, error) {
	rows, err := m.db.Query(`
        SELECT _id, id, chat_id, user_id, type, message, attachment, v
        FROM db1.chat_logs
        WHERE _id > ? ORDER BY _id ASC`, rowId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vJson string
	var messages []model.Message

	for rows.Next() {
		msg := model.Message{}

		err := rows.Scan(&msg.RowId, &msg.MsgId, &msg.ChatId, &msg.UserId, &msg.MsgType, &msg.Message, &msg.Attachment, &vJson)
		if err != nil {
			continue
		}

		var v struct {
			Enc uint32 `json:"enc"`
		}

		err = json.Unmarshal([]byte(vJson), &v)
		if err != nil {
			return nil, fmt.Errorf("failed to Unmarshal: %v", err)
		}

		err = decrypt.DecryptMessage(&msg, v.Enc)
		if err != nil {
			fmt.Println("DecryptError", msg, err)
			continue
		}

		messages = append(messages, msg)
	}

	return messages, nil
}
