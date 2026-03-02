package utils

import (
	"fmt"
	"time"
)

// FormatSalaryRange 格式化薪资范围
func FormatSalaryRange(min, max int) string {
	if max == 0 || max >= 15000 {
		return fmt.Sprintf("%d及以上", min)
	}
	return fmt.Sprintf("%d-%d", min, max)
}

// FormatWorkNature 工作性质转换(数据库int -> 字符串)
func FormatWorkNature(code int) string {
	switch code {
	case 1:
		return "实习招聘"
	case 2:
		return "校园招聘"
	default:
		return "未知"
	}
}

// ParseWorkNatureToCode 工作性质转换(字符串 -> 数据库int)
func ParseWorkNatureToCode(nature string) int {
	switch nature {
	case " 实习招聘":
		return 1
	case "校园招聘":
		return 2
	default:
		return 0
	}
}

// FormatDegreeRequirement 学历要求转换
func FormatDegreeRequirement(code int) string {
	switch code {
	case 0:
		return "本科及以上"
	case 1:
		return "硕士及以上"
	case 2:
		return "博士及以上"
	default:
		return "不限"
	}
}

// FormatDate 日期格式化为YYYY-MM-DD
func FormatDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}

// FormatDatePtr 日期指针格式化为YYYY-MM-DD
func FormatDatePtr(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}
