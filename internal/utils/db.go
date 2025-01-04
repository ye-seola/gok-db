package utils

import (
	"gokdb/internal/constsnts"
	"os"
	"time"
)

func GetDB1ModifiedTime() (time.Time, error) {
	fio, err := os.Stat(constsnts.DB1WalPath)
	if err != nil {
		return time.Time{}, err
	}
	return fio.ModTime(), nil
}

func GetDB2ModifiedTime() (time.Time, error) {
	fio, err := os.Stat(constsnts.DB2WalPath)
	if err != nil {
		return time.Time{}, err
	}
	return fio.ModTime(), nil
}
