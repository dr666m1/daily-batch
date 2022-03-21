package date

import (
	"time"
)

const sqlDateLiteral = "2006-01-02"

func Today() string {
	location, _ := time.LoadLocation("Asia/Tokyo")
	dt := (time.Now().In(location)).Format(sqlDateLiteral)
	return dt
}

func OneDayBefore(dt string) string {
	ts, _ := time.Parse(sqlDateLiteral, dt)
	yesterday := ts.AddDate(0, 0, -1).Format(sqlDateLiteral)
	return yesterday
}
