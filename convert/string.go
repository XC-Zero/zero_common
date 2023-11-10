package convert

import (
	"crypto/md5"
	"encoding/hex"
	"hash"
	"strings"
	"sync"
	"time"
	"unicode"
)

var md hash.Hash
var mdLock sync.Mutex
var mdOnce sync.Once

var salt = []byte("tessan_erp")

func MD5(str string) string {
	mdOnce.Do(func() {
		md = md5.New()
	})
	mdLock.Lock()
	defer mdLock.Unlock()
	defer md.Reset()
	md.Write([]byte(str))
	return hex.EncodeToString(md.Sum(salt))
}

func splitText(str string) []string {
	nn := strings.ReplaceAll(str, "\"", "")
	return strings.Split(nn, "\t")
}

func IsFirstUpper(str string) bool {
	if len(str) == 0 {
		return false
	}
	return unicode.IsUpper([]rune(str)[0])
}

func TimeJoin(times ...time.Time) string {
	var str []string
	for _, t := range times {
		str = append(str, t.String())
	}
	return strings.Join(str, ",")
}

func EqualStrPtr(t1, t2 *string) bool {
	if t1 == t2 {
		return true
	}
	if t1 != nil && t2 != nil {
		if *t1 == *t2 {
			return true
		}
	}
	return false
}

// IsEmptyStr 判断字符串指针是否是空的 或全是空格
func IsEmptyStr(strs ...*string) bool {
	for i := 0; i < len(strs); i++ {
		if strs[i] != nil && strings.TrimSpace(*strs[i]) != "" {
			return false
		}
	}
	return true
}

func DealX(bs []byte) []rune {
	var res []rune
	for _, b := range bs {
		res = append(res, rune(b))
	}
	return res
}
