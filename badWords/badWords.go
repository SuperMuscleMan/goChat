package badWords

import (
	"bufio"
	"os"
)

//屏蔽字库结构
type badWords struct {
	isWord bool
	words  map[string]badWords
}

// 屏蔽字库map
var badRoot = make(map[string]badWords)

// 初始化脏词配置
func Init(path string) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = f.Close()
	}()
	r := bufio.NewReader(f)
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			return
		}
		lineStr := string(line)
		InitAdd(badRoot, 0, len(lineStr), lineStr)
	}
}

// 加入map中
func InitAdd(maps map[string]badWords, index int, max int, r string) {

	if element, ok := maps[string(r[index])]; ok {
		if (index + 1) == max {
			element.isWord = true
			maps[string(r[index])] = element
			return
		}
		InitAdd(element.words, index+1, max, r)
	} else {
		subElement := badWords{false, make(map[string]badWords)}
		maps[string(r[index])] = subElement
		if (index + 1) == max {
			subElement.isWord = true
			maps[string(r[index])] = subElement
			return
		}
		InitAdd(subElement.words, index+1, max, r)
	}
}

// 检测是否有屏蔽字
func HandelBad(str string) string {
	goodStr := []byte(str)
	for i := 0; i < len(goodStr); i++ {
		if element, ok := badRoot[string(goodStr[i])]; ok {
			offset := checkBad_(element, -1, str[i+1:])
			if offset != -1 {
				for j := i; j <= i+offset; j++ {
					goodStr[j] = '*'
				}
				i += offset
			}
		}
	}
	return string(goodStr)
}
func checkBad_(badWord badWords, index int, str string) int {
	if badWord.isWord {
		index = 0
	}
	for i, r := range str {
		if element, ok := badWord.words[string(r)]; ok {
			if element.isWord {
				index = i + 1
			}
			badWord = element
		} else {
			return index
		}
	}
	return index
}
