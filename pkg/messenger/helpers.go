package messenger

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/fatih/color"
)

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func parseTime(timeStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeStr)
}

func reconnect(connMethod func() error) {
	for {
		for i := 1; i <= 5; i++ {
			fmt.Println("connect to Utopia Network..")
			err := connMethod()
			if err == nil {
				color.Green("connected")
				return
			}

			color.Yellow("[WARN] connection failed")
			time.Sleep(time.Second * 12)
		}

		color.Yellow("[WARN] retry after 40s..")
		time.Sleep(time.Second * 40)
	}
}
