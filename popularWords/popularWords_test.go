package popularWords

import (
	"testing"
)

// 测试是否正常统计高频流行词
func TestStatistic(t *testing.T) {
	strs := []string{"hello", "world", "beautiful", "world"}
	Statistic(strs)
	popular, ok := GetPopular(1)
	if popular != "world" {
		t.Errorf(`Statistic failed.popular=%s`, popular)
	}
	if !ok {
		t.Error(`GetPopular(1) = false`)
	}

}
// 测试获取流行词的时间条件
func TestNonGetPopular(t *testing.T) {
	_, ok := GetPopular(61)
	if ok {
		t.Error(`GetPopular(61) = true`)
	}
}
// 测试是否清理流行词库
func TestTimer(t *testing.T) {
	strs := []string{"hello", "world", "beautiful", "world"}
	Statistic(strs)
	timerLogic() // 清理定时器，运行test时会等待60s
	if len(wordMaps) > 0 {
		t.Error(`uncleaned`)
	}
}
