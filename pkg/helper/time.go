package helper

import "time"

func GetCurrentTime() time.Time {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	currentTime, _ := time.Parse("2006-01-02 15:04:05", time.Now().In(loc).Format("2006-01-02 15:04:05"))
	return currentTime
}
