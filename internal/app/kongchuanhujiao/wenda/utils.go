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
func AnswerToCSV(ans []*wenda.AnswersTab, mem wenda.GroupMembers) (r []byte) {

	bf := bytes.NewBuffer(r)
	w := csv.NewWriter(bf)

	_ = w.Write([]string{"用户名", "QQ", "答题时间", "答题内容"})
	for _, an := range ans {
		_ = w.Write([]string{mem[an.QQ], strconv.FormatUint(an.QQ, 10), an.Time, an.Answer})
	}
	w.Flush()

	r = bf.Bytes()
	return
}
