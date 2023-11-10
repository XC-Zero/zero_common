package img

import "strings"

func SuffixIsImage(fileName string) bool {
	var list = strings.Split(fileName, ".")
	if len(list) == 0 {
		return false
	}
	switch strings.ToLower(list[len(list)-1]) {
	case "tiff", "png", "gif", "jpeg", "jpg", "svg", "bmp", "webp":
		return true
	default:
		return false
	}
}
