package common
import (
	"fmt"
// "strings"
	"time"
	"unicode/utf8"
	"crypto/hmac"
	"crypto/sha256"

)

func LevelShift(tab int) string {
	res := ""
	for i := 0; i < tab; i += 1 {
		res = res + "\t"
	}
	return res
}

func GetSlice(text string, begin_pos int, end_pos int) string {
	result := ""
	for i, w := 0, 0; i < len(text); i += w {
		runeValue, width := utf8.DecodeRuneInString(text[i:])
		w = width
		s1 := string(runeValue)
		if i >= begin_pos {
			if i < end_pos {
				result = result + s1
			}
		}
	}
	return result
}

func GetMID() string {
        appKey := "GetMID"
	id := fmt.Sprintf("%v", time.Now())
	mac := hmac.New(sha256.New, []byte(id))
	mac.Write([]byte(appKey))
	return fmt.Sprintf("%x", mac.Sum(nil))
}
