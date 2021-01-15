package wenda

import (
	"crypto/sha1"
	"encoding/hex"
	"regexp"
	"strings"
)

// checkAnswerForSelect 检查为合法的 [选择题] 答案
func checkAnswerForSelect(a string) (ok bool) {
	ok, _ = regexp.MatchString("^[A-Za-z]$", a)
	return
}

// checkAnswerForFill 检查为合法的 [简答题] 答案
func checkAnswerForFill(a string) bool { return strings.HasPrefix(a, "#") }

// HashForSHA1 SHA1 散列 FIXME 如果没有引用则应该废弃该函数
func HashForSHA1(d string) string {
	h := sha1.New()
	_, _ = h.Write([]byte(d))
	return hex.EncodeToString(h.Sum(nil))
}

/* TODO 分词功能由其他包实现
func DoWordSplit(s string) (words []string, err error) {
	// FIXME
	words, err = AC.Cli.C.GetWordSegmentation(s)

	if err != nil {
		return
	}

	for k, v := range words {
		words[k] = strings.ReplaceAll(v, "\u0000", "")
	}

	return words, nil
}
*/
