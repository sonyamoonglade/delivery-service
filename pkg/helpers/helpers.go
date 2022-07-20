package helpers

import (
	"strconv"
	"strings"
)

func SixifyOrderId(id int64) string {

	idLikeSix := strconv.Itoa(int(id))

	if len(idLikeSix) == 6 {
		return idLikeSix
	}
	for {
		l := len(idLikeSix)

		if l == 6 {
			break
		}
		idLikeSix = "0" + idLikeSix
	}
	return idLikeSix
}

func ExtractOrderId(text string) int64 {
	strOrderId := ""
	arr := strings.Split(text, "")
	it := 0
	for i := len(arr) - 1; i >= 0; i-- {
		ch := arr[i]
		if it == 6 {
			if it == 0 {
				return 0
			}
			break
		}
		strOrderId = ch + strOrderId
		it += 1
	}
	numId, _ := strconv.ParseInt(strOrderId, 10, 64)
	return numId
}

func ExtractTotalPrice(text string) int64 {
	arr := strings.Split(text, "")
	sumStr := ""
	for _, ch := range arr {
		if ch == "." {
			break
		}
		if _, err := strconv.ParseInt(ch, 10, 64); err != nil {
			continue
		}
		sumStr += ch
	}
	sum, _ := strconv.ParseInt(sumStr, 10, 64)
	return sum
}

func ExtractUsername(text string) string {
	spl := strings.Split(text, ": ")
	name := spl[1]
	return name
}
