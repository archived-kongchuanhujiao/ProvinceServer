package wenda

import (
	"bytes"
	"crypto/sha1"
	"encoding/csv"
	"encoding/hex"
	"github.com/kongchuanhujiao/server/internal/app/datahub/public/wenda"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// checkAnswerForSelect 检查为合法的 [选择题] 答案
func checkAnswerForSelect(a string) (ok bool) {
	ok, _ = regexp.MatchString("^[A-Za-z]$", a)
	return
}

// checkAnswerForFill 检查为合法的 [简答题] 答案
func checkAnswerForFill(a string) bool { return strings.HasPrefix(a, "#") }

// HashForSHA1 SHA1 散列
func HashForSHA1(d string) string {
	h := sha1.New()
	_, _ = h.Write([]byte(d))
	return hex.EncodeToString(h.Sum(nil))
}

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// AnswerToCSV 将答案转换为 CSV
func AnswerToCSV(ans []*wenda.AnswersTab) (r []byte, err error) {
	bf := bytes.NewBuffer(r)
	w := csv.NewWriter(bf)
	err = w.Write([]string{"QQ", "答题时间", "作答答案"})
	for _, an := range ans {
		err = w.Write([]string{strconv.FormatUint(an.QQ, 10), an.Time, an.Answer})
	}
	w.Flush()

	r = bf.Bytes()

	return
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
