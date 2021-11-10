package popularWords

import "time"

type wordTimeInfo struct {
	time int64
	num  int
}

type wordInfo struct {
	popular string
	list []wordTimeInfo
}

const (
	Expire = 60
)

var wordMaps = make(map[string]wordInfo)

func Timer(){
	for{
	timeAfterTrigger := time.After(Expire * time.Second)
	<-timeAfterTrigger
		expireTime := time.Now().Unix() - Expire
		for _, r := range wordMaps {
			max := len(r.list)
			for i := max - 1; i >= 0;i --  {
				element := r.list[i]
				if element.time < expireTime {
					r.list = append(r.list[:i], r.list[i+1:]...)
				}
			}
			wordMaps[r.popular] = r
		}
	}
}
func GetPopular(limit int64) (string, bool) {
	if limit > 60 {
		return "", false
	}
	limitTime := time.Now().Unix() - limit
	max := 0
	popular := ""
	for _, r := range wordMaps {
		num := 0
		for _, e := range r.list {
			if e.time >limitTime {
				num += e.num
			}
		}
		if num > max {
			max = num
			popular = r.popular
		}
	}
	return popular, true
}
func Statistic(words []string) {
	curTime := time.Now().Unix()
	for _, r := range words {
		add(curTime, r)
	}
}
func add(time int64, string string) {
	word := wordMaps[string]
	for i, e := range word.list {
		if e.time == time {
			e.num += 1
			word.list[i] = e
			word.popular = string
			wordMaps[string] = word
			return
		}
	}
	wordTime := wordTimeInfo{time, 1}
	word.list = append(word.list, wordTime)
	word.popular = string
	wordMaps[string] = word
}
