package messenger

import (
	"crypto/md5"
	"encoding/hex"
	"time"
)

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func parseTime(timeStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeStr)
}
