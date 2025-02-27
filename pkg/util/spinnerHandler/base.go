package spinnerHandler

import (
	"time"

	"github.com/briandowns/spinner"
)

var CharSets = spinner.CharSets

type Option struct {
	Charts       []string      // 动画字符集
	Prefix       string        // 前缀
	Suffix       string        // 后缀
	AnmationTime time.Duration // 动画时间 单位毫秒 默认100
}

func CreateSpinner(option Option) *spinner.Spinner {
	if len(option.Charts) == 0 {
		option.Charts = spinner.CharSets[35]
	}
	if option.AnmationTime == 0 {
		option.AnmationTime = 100
	}
	s := spinner.New(option.Charts, option.AnmationTime*time.Millisecond) // Build our new spinner
	s.Prefix = option.Prefix
	s.Suffix = option.Suffix
	return s
}
