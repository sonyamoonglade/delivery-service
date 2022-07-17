package tgdelivery

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

func ValidatePhoneNumber(v string) bool {
	spl := strings.Split(v, "")
	if len(v) != 12 || spl[0] != "+" || spl[1] != "7" || spl[2] != "9" {
		return false
	}
	return true
}

func ValidateUsername(v string) bool {
	spl := strings.Split(v, " ")
	min := 3
	if len(spl) != 2 {
		return false
	}
	for _, w := range spl {
		wordSpl := strings.Split(w, "")
		if len(wordSpl) < min {
			return false
		}
	}

	return true
}
