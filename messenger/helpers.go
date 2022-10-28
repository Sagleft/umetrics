package messenger

import (
	"crypto/md5"
	"encoding/hex"
)

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
