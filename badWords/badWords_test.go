package badWords

import (
	"testing"
)

// 测试是否正常检测屏蔽字
func TestInit(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{"4r5e", "****"},
		{"5h1t", "****"},
		{"5hit", "****"},
		{"a55", "***"},
		{"anal", "****"},
		{"anus", "****"},
		{"ar5e", "****"},
		{"arrse", "*****"},
		{"a", "*"},
		{"as", "**"},
		{"ars", "***"},
	}
	for _, test := range tests {
		InitAdd(badRoot, 0, len(test.input), test.input)
		if r := HandelBad(test.input); r != test.want {
			t.Errorf("HandelBad(%s)=%s", test.input, r)
		}

	}
}
