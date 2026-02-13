package utils

import (
	"errors"
	"time"
)

// ParseYYYYMM 将 "YYYY-MM" 格式字符串转换为 time.Time
// 例如: "2024-06" -> time.Time{2024, 6, 1, 0, 0, 0, 0}
func ParseYYYYMM(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, errors.New("date string is empty")
	}

	// 解析为 "2006-01" 格式 (Go的时间格式化模板)
	// 统一设置为每月1日
	t, err := time.Parse("2006-01", dateStr)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

// FormatYYYYMM 将 time.Time 格式化为 "YYYY-MM" 字符串
func FormatYYYYMM(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01")
}

// FormatYYYYMMPtr 处理 *time.Time 类型,支持 NULL 值
// 如果指针为 nil 或时间为零值,返回 nil
func FormatYYYYMMPtr(t *time.Time) *string {
	if t == nil || t.IsZero() {
		return nil
	}
	formatted := t.Format("2006-01")
	return &formatted
}

// ParseYYYYMMOptional 解析可选的日期字符串
// 如果字符串为空,返回零值 time.Time 而不是错误
func ParseYYYYMMOptional(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, nil
	}
	return ParseYYYYMM(dateStr)
}
