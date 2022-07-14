package tgdelivery

import "strconv"

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
