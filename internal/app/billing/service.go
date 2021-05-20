package billing

import (
	"time"
)

// GetMonthFirstDate 获取月第一天和最后一天
// 获取月第一天 shift为偏移量 正为下月 负为上月
func GetMonthFirstDate(t time.Time, shift int) string {
	t = t.AddDate(0, shift, -t.Day()+1)
	return t.Format(timeFormat)

}

// GetMonthLastDate 获取月最后一天 shift为偏移量 正为下月 负为上月
func GetMonthLastDate(t time.Time, shift int) string {
	t = t.AddDate(0, shift, -t.Day())
	return t.Format(timeFormat)
}

// GetMonthDate 获取月第一天和最后一天
func GetMonthDate() map[string]string {
	data := make(map[string]string)

	t := time.Now()
	data["thisMonthFirstDate"] = GetMonthFirstDate(t, 0)  // 本月第一天
	data["thisMonthLastDate"] = GetMonthLastDate(t, 1)    // 本月最后一天
	data["lastMonthFirstDate"] = GetMonthFirstDate(t, -1) //上月第一天
	data["lastMonthLastDate"] = GetMonthLastDate(t, 0)    //上月最后一天

	return data
}
